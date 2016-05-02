package main

import (
       "os"
       "strings"
	"encoding/csv"
)

type CsvOut struct {
}

func (to CsvOut) output(outputLocation string, outputBaseFileName string, header string, dirsOfInterest []string, s []string, packageInfo map[string]*PackageInfo, groupSet map[string]bool, nodes map[string]*Node) error {

     //     	var b bytes.Buffer
	w := csv.NewWriter(os.Stdout)

	w.Write([]string{"name","group","license","version","buildtime","description"})

	for r := range s {
	        s := []string{packageInfo[s[r]].Name, packageInfo[s[r]].Tags["group"], packageInfo[s[r]].Tags["license"], packageInfo[s[r]].Tags["version"], packageInfo[s[r]].Tags["buildtime"], strings.Replace(packageInfo[s[r]].Tags["description"], "\n", "   ", -1)}
		w.Write(s)
//		w.WriteTo(os.Stdout)
	}

	return nil

}
