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
)


type RpmInfo struct{
	name string
	tags map[string]string
}


var allRpmsExec = []string{"rpm", "-qa"}
var numFindOfInterestRpmsWorkers = 8
var numFindRpmTagsWorkers = 3

var	tagSeparator = "--==--"
var	recordSeparator = "|||"

var rpmFormat = "\"name:%{NAME}" + tagSeparator + "os:%{OS}" + tagSeparator + "version:%{VERSION}" + tagSeparator + "release:%{RELEASE}" + tagSeparator + "arch:%{ARCH}" + tagSeparator + "installtime:%{INSTALLTIME:date}" + tagSeparator + "group:%{GROUP}" + tagSeparator + "size:%{SIZE}" + tagSeparator + "license:%{LICENSE}" + tagSeparator + "sourcerpm:%{SOURCERPM}" + tagSeparator + "buildtime:%{BUILDTIME}" + tagSeparator + "buildhost:%{BUILDHOST}" + tagSeparator + "packager:%{PACKAGER}" + tagSeparator + "vendor:%{VENDOR}" + tagSeparator + "url:%{URL}" + tagSeparator + "summary:%{SUMMARY}" + tagSeparator + "description:%{DESCRIPTION}" + tagSeparator + "" + tagSeparator + "distribution:%{DISTRIBUTION}" + tagSeparator + "packager:%{PACKAGER}" + tagSeparator + "patch:%{PATCH}" + recordSeparator + "\""

//var rpmFormat = "\"name:%{NAME}\""

func main(){
	rand.Seed(time.Now().Unix())
	args := os.Args
	numArgs := len(args)
	var dirsOfInterest []string

	if(numArgs > 1){
		dirsOfInterest = args[1:numArgs]
	}else{
		dirsOfInterest = [] string{"/"}
	}

	xlog("Start")

	rpmListChannel, rpmListDoneChannel, err := makeRpmList()
	if err != nil {
		log.Fatal("jjjj ", err)
	}

	ofInterestChannel, ofInterestDoneChannel, err := findOfInterestRpms(dirsOfInterest, rpmListChannel, rpmListDoneChannel)

	tagInfoChannel, tagInfoDoneChannel, err :=	findRpmTags(ofInterestChannel, ofInterestDoneChannel)

	//addToMapDoneChannel, rpmInfoMap := addResultsToMap(tagInfoChannel, tagInfoDoneChannel)
	addToMapDoneChannel, rpmInfoMap := addResultsToMap(tagInfoChannel, tagInfoDoneChannel)
	_ = <- addToMapDoneChannel

	mk := make([]string, len(rpmInfoMap))
    i := 0
    for k, _ := range rpmInfoMap {
        mk[i] = k
        i++
    }
    sort.Strings(mk)
    fmt.Println(mk)

	xlog("End")
	fmt.Println("__________________________________")
}

func xlog(m string){
	fmt.Println(m)
}

func makeRpmList()(chan *RpmInfo, chan bool, error){
	xlog("Start makeRpmList")
	rpmListChannel := make(chan *RpmInfo, 1000)
	rpmListDoneChannel := make(chan bool) 
	
	stringListChannel, stringDoneChannel := runExec(allRpmsExec)

	go func(){
		xlog("Start makeRpmList GO 1")
		finished := false
		for !finished{
			select{
			case rpmName := <- stringListChannel:
				rpmInfo := new(RpmInfo)
				rpmInfo.name = rpmName
				rpmInfo.tags = make(map[string] string)
				rpmListChannel<- rpmInfo
			case <- stringDoneChannel:
				finished = true;
				for i:=0; i<numFindOfInterestRpmsWorkers; i++{
					rpmListDoneChannel<- true
				}
			}
		}
	}()
	return 	rpmListChannel, rpmListDoneChannel, nil
}

func findOfInterestRpms(dirsOfInterest []string, rpmListChannel chan *RpmInfo, rpmListDoneChannel chan bool)(chan *RpmInfo, chan bool, error){
	ofInterestChannel := make(chan *RpmInfo, 1000)
	ofInterestDoneChannel := make(chan bool) 

	for i:=0; i<numFindOfInterestRpmsWorkers; i++{
		go func(){
			finished := false
			for !finished{
				select{
				case rpmInfo := <- rpmListChannel:
					cmd := exec.Command("rpm", "-ql", rpmInfo.name)
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
							}
						}
					}
				case <- rpmListDoneChannel:
					finished = true
					for j:=0; j<numFindRpmTagsWorkers; j++{
						ofInterestDoneChannel <- true
					}
					
				}
			}
		}()
	}
	return 	ofInterestChannel, ofInterestDoneChannel, nil
}



func findRpmTags(ofInterestChannel chan *RpmInfo, ofInterestDoneChannel chan bool)(chan *RpmInfo, chan bool, error){
	tagInfoChannel := make(chan *RpmInfo, 1000)
	tagInfoDoneChannel := make(chan bool) 

	for i:=0; i<numFindRpmTagsWorkers; i++{
		go func(){
			numTags := random(100,400)
			tagBuffer:= make([] *RpmInfo, numTags)
			count :=0
			finished := false
			for !finished{
				select{
				case rpmInfo := <- ofInterestChannel:
					fmt.Println("--",count," ", rpmInfo.name)
					if count >= numTags{
						out := runCommand(tagBuffer, count)
						parseAndSend(count, tagBuffer, out, tagInfoChannel)

						//fmt.Println("+++ ", rpmInfo.name, " --- ", out)
						count = 0
						tagInfoChannel<- rpmInfo

					}
					tagBuffer[count] = rpmInfo
					count = count + 1

				case <- ofInterestDoneChannel:
					fmt.Println("+++++++++++++++++++++")
					if count > 0{
						out := runCommand(tagBuffer, count)
						parseAndSend(count, tagBuffer, out, tagInfoChannel)
						//_ = runCommand(tagBuffer, count)
						//fmt.Println("+++ ", " --- ", out)
					}
					finished = true
					tagInfoDoneChannel <- true
				}
			}
		}()
	}
	return 	tagInfoChannel, tagInfoDoneChannel, nil
}


func parseAndSend(count int, tagInfoBuffer [] *RpmInfo, out string, tagInfoChannel chan *RpmInfo){
	//records := strings.Split(out, recordSeparator)
	records := strings.Split(out, recordSeparator)
	if len(records) != count+1{
		fmt.Println(len(records), "-", count)
		log.Fatal("jiijjj ")
	}
	for i:=0; i<count; i++{
		fmt.Println(tagInfoBuffer[i].name, "%%%%%%%%%%%%%%%%%%%%%%%%%%::" + records[i])
		tags := strings.Split(records[i], tagSeparator)
		for j:=0; j<len(tags); j++{
			fmt.Println("\t", j, " ", tags[j])
			parts := strings.SplitN(tags[j], ":", 2)
			if(len(parts) == 2){
				tagInfoBuffer[i].tags[parts[0]] = parts[1]
			}
		}
	}
	
}

func runCommand(tagBuffer []*RpmInfo, count int)(string){
	args := []string{}
	args = append(args, "--qf")
	args = append(args, rpmFormat)
						args = append(args,  "-q",)
	for j:=0; j<count; j++{
		args = append(args,  tagBuffer[j].name)
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
		cmd.Args[3+i] = tagBuffer[i].name
	}
}




func addResultsToMap(tagInfoChannel chan *RpmInfo, tagInfoDoneChannel chan bool) ( chan bool, map[string]*RpmInfo) {
	addToMapDoneChannel := make(chan bool) 
	rpmMap := make(map[string] *RpmInfo)

	go func(){
		finished := false
		for !finished{
			select{
			case rpmInfo := <- tagInfoChannel:
				//case _ = <- tagInfoChannel:
				rpmMap[rpmInfo.name] = rpmInfo
			case <- tagInfoDoneChannel:
				finished = true
				addToMapDoneChannel <- true
			}
		}
	}()
	
	return 	addToMapDoneChannel, rpmMap

}


func random(min, max int) int {
    return rand.Intn(max - min) + min
}