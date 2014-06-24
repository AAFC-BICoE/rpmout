package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type RPackage struct {
	Built                 string
	Depends               string
	Description           string
	Enhances              string
	Imports               string
	LibPath               string
	License               string
	License_is_FOSS       string
	License_restricts_use string
	LinkingTo             string
	MD5sum                string
	NeedsCompilation      string
	OS_type               string
	Package               string
	Priority              string
	Suggests              string
	Title                 string
	URL                   string
	Version               string
}

const Description = "description"
const License = "license"
const Name = "name"
const Summary = "summary"
const URL = "url"
const Version = "version"

func findRPackages() map[string]*PackageInfo {
	cmd := exec.Command("rpmout.R")
	jsonBytes, err := cmd.Output()
	if err != nil {
		log.Fatal("jjjj ", err)
	}

	var rpackages []RPackage
	err = json.Unmarshal(jsonBytes, &rpackages)
	if err != nil {
		fmt.Println("error:", err)
	}
	//fmt.Printf("%+v", rpackages)

	packageInfoMap := convertRPackages(rpackages)
	return packageInfoMap
}

func convertRPackages(rpackages []RPackage) map[string]*PackageInfo {
	packageInfoMap := make(map[string]*PackageInfo, 200)

	for _, p := range rpackages {
		//for k, _ := range m { ... }
		// fmt.Println("8888888888888888888888888")
		// fmt.Println(i)
		// fmt.Printf("%+v", p)
		// fmt.Println("")
		packageInfo := convertRPackage(p)
		//fmt.Printf("*************** ", Name, packageInfo.Tags[Name])
		//packageInfoMap[strings.ToLower(packageInfo.IndexName)] = packageInfo
		packageInfoMap[strings.ToLower(packageInfo.Name)] = packageInfo
		//packageInfoMap[packageInfo.Name] = packageInfo

	}

	return packageInfoMap
}

func convertRPackage(rpackage RPackage) *PackageInfo {
	packageInfo := new(PackageInfo)
	packageInfo.IsR = true

	packageInfo.Tags = make(map[string]string)
	packageInfo.Tags[Name] = rpackage.Package
	packageInfo.Name = rpackage.Package
	packageInfo.IndexName = "R-packages!" + packageInfo.Tags[Name]

	packageInfo.Tags[Summary] = rpackage.Title
	packageInfo.Tags[Description] = rpackage.Description
	packageInfo.Tags[Version] = rpackage.Version
	packageInfo.Tags[License] = rpackage.License
	packageInfo.Tags[URL] = rpackage.URL
	return packageInfo
}
