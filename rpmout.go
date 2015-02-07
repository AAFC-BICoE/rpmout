package main

// Glen Newton
// glen.newton@gmail.com

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"
)

func init() {
	flag.StringVar(&outputFormat, "outputFormat", "html", "Values: html|html2|json|txt|latex|exhibit")
	flag.StringVar(&header, "header", "Installed Software", "gggg")
	flag.BoolVar(&doR, "R", false, "Find R packages")
}

var outputFormat string
var header string
var doR = false

type PackageInfo struct {
	Name      string
	IndexName string
	IsR       bool
	Tags      map[string]string
}

type RpmWriter interface {
	output(string, []string, []string, map[string]*PackageInfo, map[string]bool, map[string]*Node) error
}

func handleParameters() bool {
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
	switch outputFormat {
	case "html":
		rpmWriter = new(HtmlOut)
	case "html2":
		rpmWriter = new(HtmlOut2)
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

var allRpmsExec = []string{"rpm", "-qa"}

var numFindOfInterestRpmsWorkers = 6
var numFindRpmTagsWorkers = 2

var rpmWriter RpmWriter

func main() {

	runtime.GOMAXPROCS(8)

	_, err := exec.LookPath(allRpmsExec[0])

	if err != nil {
		log.Fatal("'", allRpmsExec[0], "' is not in your path: needed to run; it is usually in '/usr/bin/rpm'")
	}

	if !handleParameters() {
		log.Fatal("Bad parameters")
	}

	rand.Seed(time.Now().Unix())

	var dirsOfInterest []string
	if len(flag.Args()) == 0 {
		dirsOfInterest = []string{"/"}
	} else {
		dirsOfInterest = flag.Args()
	}

	//fmt.Println(dirsOfInterest)

	rpmListChannel, err := makeRpmList()
	if err != nil {
		log.Fatal("jjjj ", err)
	}

	ofInterestChannel, ofInterestDoneChannel, err := findOfInterestRpms(dirsOfInterest, rpmListChannel)

	tagInfoChannel, tagInfoDoneChannel, err := findRpmTags(ofInterestChannel)

	packageInfoMap := make(map[string]*PackageInfo, 200)
	//packageInfoMap2 := make(map[string] PackageInfo, 200)
	groupSet := make(map[string]bool)

	var nodes map[string]*Node
	nodes = make(map[string]*Node)

	addResultsDoneChannel := addResultsToMap(packageInfoMap, groupSet, tagInfoChannel, nodes)

	for i := 0; i < numFindOfInterestRpmsWorkers; i++ {
		_ = <-ofInterestDoneChannel
	}
	close(ofInterestChannel)

	for i := 0; i < numFindRpmTagsWorkers; i++ {
		_ = <-tagInfoDoneChannel
	}
	close(tagInfoChannel)

	_ = <-addResultsDoneChannel

	var rpackageInfoMap map[string]*PackageInfo = nil
	if doR {
		rpackageInfoMap = findRPackages()
	}

	add(packageInfoMap, rpackageInfoMap)
	rpmWriter.output(header, dirsOfInterest, sortStringKeyMap(packageInfoMap), packageInfoMap, groupSet, nodes)
}

func xlog(m string) {
	fmt.Println(m)
}

func makeRpmList() (chan *PackageInfo, error) {
	rpmListChannel := make(chan *PackageInfo, 200)

	stringListChannel := runExec(allRpmsExec)

	go func() {
		for stringInfo := range stringListChannel {
			//fmt.Println("makeRpmList: [",  stringInfo.val, "] ", &stringInfo.val," ", stringInfo.done)
			packageInfo := new(PackageInfo)
			packageInfo.Name = stringInfo.val
			packageInfo.Tags = make(map[string]string)
			rpmListChannel <- packageInfo
		}
		close(rpmListChannel)
	}()
	return rpmListChannel, nil
}

func findOfInterestRpms(dirsOfInterest []string, rpmListChannel chan *PackageInfo) (chan *PackageInfo, chan bool, error) {
	ofInterestChannel := make(chan *PackageInfo, 100)
	doneChannel := make(chan bool, numFindOfInterestRpmsWorkers)

	for i := 0; i < numFindOfInterestRpmsWorkers; i++ {
		go func() {
			for packageInfo := range rpmListChannel {
				//fmt.Println("findOfInterestRpms ", packageInfo.Name)
				cmd := exec.Command(allRpmsExec[0], "-ql", packageInfo.Name)
				out, err := cmd.Output()

				if err != nil {
					log.Fatal("jjjj ", err)
				}

				lines := strings.Split(string(out), "\n")
				done := false
				for _, line := range lines {
					if done {
						break
					}
					for _, dir := range dirsOfInterest {
						if strings.Index(line, dir) == 0 {
							ofInterestChannel <- packageInfo
							done = true
						}
					}
				}
			}
			doneChannel <- true
		}()
	}
	return ofInterestChannel, doneChannel, nil
}

func findRpmTags(ofInterestChannel chan *PackageInfo) (chan *PackageInfo, chan bool, error) {
	tagInfoChannel := make(chan *PackageInfo, 200)
	doneChannel := make(chan bool, numFindRpmTagsWorkers)

	for i := 0; i < numFindRpmTagsWorkers; i++ {
		go func() {
			numTags := random(100, 400)
			tagBuffer := make([]*PackageInfo, numTags)
			count := 0
			for packageInfo := range ofInterestChannel {
				//fmt.Println("\t findRpmTags ", packageInfo.name)
				//fmt.Println("--",count," ", packageInfo.name)
				if count >= numTags {
					out := runCommand(tagBuffer, count)
					parseAndSend(count, tagBuffer, out, tagInfoChannel)
					//fmt.Println("+++ ", packageInfo.name, " --- ", out)
					count = 0
					tagInfoChannel <- packageInfo

				}
				tagBuffer[count] = packageInfo
				count = count + 1
			}
			if count > 0 {
				out := runCommand(tagBuffer, count)
				parseAndSend(count, tagBuffer, out, tagInfoChannel)
				//_ = runCommand(tagBuffer, count)
				//fmt.Println("+++ ", " --- ", out)
			}
			doneChannel <- true
		}()
	}
	return tagInfoChannel, doneChannel, nil
}

const tagSeparator = "--==--"

const recordSeparator = "|||"

var rpmFormat = "name:%{NAME}" + tagSeparator + "os:%{OS}" + tagSeparator + "version:%{VERSION}" + tagSeparator + "release:%{RELEASE}" + tagSeparator + "arch:%{ARCH}" + tagSeparator + "installtime:%{INSTALLTIME:date}" + tagSeparator + "group:%{GROUP}" + tagSeparator + "size:%{SIZE}" + tagSeparator + "license:%{LICENSE}" + tagSeparator + "sourcerpm:%{SOURCERPM}" + tagSeparator + "buildtime:%{BUILDTIME}" + tagSeparator + "buildhost:%{BUILDHOST}" + tagSeparator + "packager:%{PACKAGER}" + tagSeparator + "vendor:%{VENDOR}" + tagSeparator + "url:%{URL}" + tagSeparator + "summary:%{SUMMARY}" + tagSeparator + "description:%{DESCRIPTION}" + tagSeparator + "" + tagSeparator + "distribution:%{DISTRIBUTION}" + tagSeparator + "packager:%{PACKAGER}" + tagSeparator + "patch:%{PATCH}" + recordSeparator

func parseAndSend(count int, tagInfoBuffer []*PackageInfo, out string, tagInfoChannel chan *PackageInfo) {
	records := strings.Split(out, recordSeparator)
	if len(records) != count+1 {
		log.Fatal("jiijjj ")
	}
	for i := 0; i < count; i++ {
		tags := strings.Split(records[i], tagSeparator)
		for j := 0; j < len(tags); j++ {
			//fmt.Println("\t", j, " ", tags[j])
			parts := strings.SplitN(tags[j], ":", 2)
			if len(parts) == 2 {
				parts[0] = strings.Trim(parts[0], " ")
				parts[1] = strings.Trim(parts[1], " ")
				if len(parts[1]) > 0 && parts[1] != "(none)" {
					tagInfoBuffer[i].Tags[parts[0]] = parts[1]
				}
			}
		}
		tagInfoChannel <- tagInfoBuffer[i]
	}
}

func runCommand(tagBuffer []*PackageInfo, count int) string {
	args := []string{}
	args = append(args, "--qf")
	args = append(args, rpmFormat)
	args = append(args, "-q")
	for j := 0; j < count; j++ {
		args = append(args, tagBuffer[j].Name)
	}
	cmd := exec.Command("rpm", args...)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal("help", err)
	}
	return string(out)
}

func makeArgs(cmd *exec.Cmd, tagBuffer []*PackageInfo, count int) {
	cmd.Args = make([]string, 3+count)
	cmd.Args[0] = "--qf"
	cmd.Args[1] = rpmFormat
	cmd.Args[2] = "-q"
	for i := 0; i < count; i++ {
		cmd.Args[3+i] = tagBuffer[i].Name
	}
}

func addResultsToMap(rpmMap map[string]*PackageInfo, groupSet map[string]bool, tagInfoChannel chan *PackageInfo, nodes map[string]*Node) chan bool {
	doneChannel := make(chan bool)
	go func() {
		for packageInfo := range tagInfoChannel {
			rpmMap[strings.ToLower(packageInfo.Name)] = packageInfo
			groupSet[packageInfo.Tags["group"]] = true
			//fmt.Println("Adding to map ", packageInfo.Name)
			extractHierarchy(packageInfo, nodes)
		}
		doneChannel <- true
	}()
	return doneChannel
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func sortStringKeyMap2(m map[string]*struct{}) []string {
	sm := make([]string, len(m))
	return sm
}

//func sortStringKeyMap(m map[string]struct{}) []string{
func sortStringKeyMap(m map[string]*PackageInfo) []string {
	sm := make([]string, len(m))
	i := 0
	for k, _ := range m {
		sm[i] = k
		i++
	}
	sort.Strings(sm)
	return sm
}

func sortLicenseCountMap(licenseCount map[string]int) []string {
	sm := make([]string, len(licenseCount))
	i := 0
	for k, _ := range licenseCount {
		sm[i] = k
		i++
	}
	sort.Strings(sm)
	return sm
}

type Pair struct {
	Key   string
	Value int
}

// From: https://groups.google.com/forum/#!topic/golang-nuts/FT7cjmcL7gw
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]int) PairList {
	pp := make(PairList, len(m))
	i := 0
	for k, v := range m {
		pp[i] = Pair{k, v}
	}
	sort.Sort(pp)
	return pp
}

func sortStringKeyNodeMap(m map[string]*Node) []string {
	sm := make([]string, len(m))
	i := 0
	for k, _ := range m {
		sm[i] = k
		i++
	}
	sort.Strings(sm)
	return sm
}

// Add the R packages to the RPM packages
func add(rpms map[string]*PackageInfo, rpacks map[string]*PackageInfo) {
	for k, v := range rpacks {
		rpms[k] = v
	}
}
