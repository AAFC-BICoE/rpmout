package main
// Glen Newton
// glen.newton@gmail.com

import(
	"log"
"fmt"
"os"
"os/exec"
"sort"
"strings"
"math/rand"
"time"
"flag"
"runtime"

)

func handleParameters() (bool){
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s: Output rpm info for rpms that have install components in certain directories\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\t %s <args> <rootDir0>...<rootDirN>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\t default <rootDir>: /\n")
		fmt.Fprintf(os.Stderr, "Args:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:  %s -outputFormat=html /opt /usr/local\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nNote that the 'rpm' program (http://www.rpm.org/max-rpm/rpm.8.html) needs to be in your path\n\n")
	}

	flag.Parse()
	switch outputFormat{
	case "html":
		rpmWriter = new(HtmlOut)
	case "json":
		rpmWriter = new(JsonOut)

	case "txt":
		rpmWriter = new(TextOut)

	case "latex":
		rpmWriter = new(LaTeXOut)
	default:
		log.Print("Unknown/unsupported output format: " + outputFormat)
		return false
	}
	return true
}

var outputFormat string
func init() {
	flag.StringVar(&outputFormat, "outputFormat", "html", "Values: html|json|txt|latex")
}

type RpmInfo struct{
	Name string
	Tags map[string]string
}

type RpmWriter interface{
	output([]string, map[string] *RpmInfo)
}


var allRpmsExec = []string{"rpm", "-qa"}
var numFindOfInterestRpmsWorkers = 6
var numFindRpmTagsWorkers = 2

var rpmWriter RpmWriter


func main(){
	runtime.GOMAXPROCS(8)
	_, err := exec.LookPath(allRpmsExec[0])
	if err != nil {
		log.Fatal("'", allRpmsExec[0],"' is not in your path: needed to run; it is usually in '/usr/bin/rpm'")
	}

	if !handleParameters(){
		log.Fatal("Bad parameters")
	}


	rand.Seed(time.Now().Unix())

	var dirsOfInterest []string
	if len(flag.Args()) == 0{
		dirsOfInterest = [] string{"/"}
	}else{
		dirsOfInterest = flag.Args()
	}

	//fmt.Println(dirsOfInterest)

	rpmListChannel, err := makeRpmList()
	if err != nil {
		log.Fatal("jjjj ", err)
	}

	ofInterestChannel, ofInterestDoneChannel, err := findOfInterestRpms(dirsOfInterest, rpmListChannel)

	tagInfoChannel,tagInfoDoneChannel, err := findRpmTags(ofInterestChannel)

	rpmInfoMap := make(map[string] *RpmInfo, 200)
	addResultsDoneChannel := addResultsToMap(rpmInfoMap, tagInfoChannel)

	for i:=0; i<numFindOfInterestRpmsWorkers; i++{
		_ = <- ofInterestDoneChannel
	}
	close(ofInterestChannel)

	for i:=0; i<numFindRpmTagsWorkers; i++{
		_ = <- tagInfoDoneChannel
	}
	close(tagInfoChannel)

	_ = <- addResultsDoneChannel

	mk := make([]string, len(rpmInfoMap))
	i := 0
	for k, _ := range rpmInfoMap {
		mk[i] = k
		i++
	}
	sort.Strings(mk)
	rpmWriter.output(mk, rpmInfoMap)
}


func xlog(m string){
	fmt.Println(m)
}


func makeRpmList()(chan *RpmInfo, error){
	rpmListChannel := make(chan *RpmInfo, 200)
	
	stringListChannel := runExec(allRpmsExec)

	go func(){
		for stringInfo := range stringListChannel{
			//fmt.Println("makeRpmList: [",  stringInfo.val, "] ", &stringInfo.val," ", stringInfo.done)
			rpmInfo := new(RpmInfo)
			rpmInfo.Name = stringInfo.val
			rpmInfo.Tags = make(map[string] string)
			rpmListChannel<- rpmInfo
		}
		close(rpmListChannel)
	}()
	return 	rpmListChannel, nil
}

func findOfInterestRpms(dirsOfInterest []string, rpmListChannel chan *RpmInfo, )(chan *RpmInfo, chan bool, error){
	ofInterestChannel := make(chan *RpmInfo, 100)
	doneChannel := make(chan bool, numFindOfInterestRpmsWorkers)

	for i:=0; i<numFindOfInterestRpmsWorkers; i++{
		go func(){
			for rpmInfo := range rpmListChannel{
				//fmt.Println("findOfInterestRpms ", rpmInfo.Name)
				cmd := exec.Command(allRpmsExec[0], "-ql", rpmInfo.Name)
				out,err := cmd.Output()
				
				if err != nil {
					log.Fatal("jjjj ", err)
				}
				
				lines := strings.Split(string(out), "\n")
				done := false
				for _, line := range lines{
					if done{
						break;
					}
					for _, dir := range dirsOfInterest{
						if strings.Index(line, dir) == 0{
							ofInterestChannel <- rpmInfo
							done = true
							//fmt.Println("findOfInterestRpms ********* ", rpmInfo.Name)
						}
					}
				}
			}
			doneChannel <- true
		}()
	}
	return 	ofInterestChannel, doneChannel, nil
}



func findRpmTags(ofInterestChannel chan *RpmInfo)(chan *RpmInfo, chan bool, error){
	tagInfoChannel := make(chan *RpmInfo, 200)
	doneChannel := make(chan bool, numFindRpmTagsWorkers)

	for i:=0; i<numFindRpmTagsWorkers; i++{
		go func(){
			numTags := random(100,400)
			tagBuffer:= make([] *RpmInfo, numTags)
			count :=0
			for rpmInfo := range ofInterestChannel{
				//fmt.Println("\t findRpmTags ", rpmInfo.name)
				//fmt.Println("--",count," ", rpmInfo.name)
				if count >= numTags{
					out := runCommand(tagBuffer, count)
					parseAndSend(count, tagBuffer, out, tagInfoChannel)
					//fmt.Println("+++ ", rpmInfo.name, " --- ", out)
					count = 0
					tagInfoChannel<- rpmInfo
					
				}
				tagBuffer[count] = rpmInfo
				count = count + 1
			}
			if count > 0{
				out := runCommand(tagBuffer, count)
				parseAndSend(count, tagBuffer, out, tagInfoChannel)
				//_ = runCommand(tagBuffer, count)
				//fmt.Println("+++ ", " --- ", out)
			}
			doneChannel <- true
		}()
	}
	return 	tagInfoChannel, doneChannel, nil
}


var	tagSeparator = "--==--"
var	recordSeparator = "|||"

var rpmFormat = "name:%{NAME}" + tagSeparator + "os:%{OS}" + tagSeparator + "version:%{VERSION}" + tagSeparator + "release:%{RELEASE}" + tagSeparator + "arch:%{ARCH}" + tagSeparator + "installtime:%{INSTALLTIME:date}" + tagSeparator + "group:%{GROUP}" + tagSeparator + "size:%{SIZE}" + tagSeparator + "license:%{LICENSE}" + tagSeparator + "sourcerpm:%{SOURCERPM}" + tagSeparator + "buildtime:%{BUILDTIME}" + tagSeparator + "buildhost:%{BUILDHOST}" + tagSeparator + "packager:%{PACKAGER}" + tagSeparator + "vendor:%{VENDOR}" + tagSeparator + "url:%{URL}" + tagSeparator + "summary:%{SUMMARY}" + tagSeparator + "description:%{DESCRIPTION}" + tagSeparator + "" + tagSeparator + "distribution:%{DISTRIBUTION}" + tagSeparator + "packager:%{PACKAGER}" + tagSeparator + "patch:%{PATCH}" + recordSeparator

func parseAndSend(count int, tagInfoBuffer [] *RpmInfo, out string, tagInfoChannel chan *RpmInfo){
	records := strings.Split(out, recordSeparator)
	if len(records) != count+1{
		log.Fatal("jiijjj ")
	}
	for i:=0; i<count; i++{
		tags := strings.Split(records[i], tagSeparator)
		for j:=0; j<len(tags); j++{
			//fmt.Println("\t", j, " ", tags[j])
			parts := strings.SplitN(tags[j], ":", 2)
			if(len(parts) == 2){
				parts[0] = strings.Trim(parts[0], " ")
				parts[1] = strings.Trim(parts[1], " ")
				if len(parts[1]) > 0 && parts[1] != "(none)"{
					tagInfoBuffer[i].Tags[parts[0]] = parts[1]
				}
			}
		}
		tagInfoChannel<-tagInfoBuffer[i]
	}
}

func runCommand(tagBuffer []*RpmInfo, count int) string {
	args := []string{}
	args = append(args, "--qf")
	args = append(args, rpmFormat)
		args = append(args,  "-q",)
	for j:=0; j<count; j++{
		args = append(args,  tagBuffer[j].Name)
	}
	cmd := exec.Command("rpm",  args...)
	out,err := cmd.Output()
	if err != nil {
		log.Fatal("help" , err)
	}
	return string(out)
}

func makeArgs(cmd *exec.Cmd, tagBuffer []*RpmInfo, count int){
	cmd.Args = make([] string ,3 + count)
	cmd.Args[0] = "--qf"
	cmd.Args[1] = rpmFormat
	cmd.Args[2] = "-q"
	for i:=0; i<count; i++{
		cmd.Args[3+i] = tagBuffer[i].Name
	}
}




func addResultsToMap(rpmMap map[string] *RpmInfo,
	tagInfoChannel chan *RpmInfo) chan bool{
	
	doneChannel := make(chan bool)

	go func(){
		for rpmInfo := range tagInfoChannel{
			rpmMap[strings.ToLower(rpmInfo.Name)] = rpmInfo
			//fmt.Println("Adding to map ", rpmInfo.Name)
		}
		doneChannel <- true
	}()
	return doneChannel
}


func random(min, max int) int {
    return rand.Intn(max - min) + min
}



