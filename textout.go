package main

import (
	"fmt"
	"strings"
)

type TextOut struct {
}

func (to TextOut) output(outputLocation string, header string, dirsOfInterest []string, s []string, packageInfo map[string]*PackageInfo, groupSet map[string]bool, nodes map[string]*Node) error {

	for r := range s {
		fmt.Println("")
		fmt.Println(packageInfo[s[r]].Name)
		for k, v := range packageInfo[s[r]].Tags {
			v = strings.Replace(v, "\n", "\n               ", -1)
			fmt.Println("  " + k + ": " + v)
		}
	}

	return nil

}
