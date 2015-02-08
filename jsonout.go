package main

import (
	"encoding/json"
	"fmt"
)

type JsonOut struct {
}

func (jo JsonOut) output(outputLocation string, header string, dirsOfInterest []string, s []string, packageInfo map[string]*PackageInfo, groupSet map[string]bool, nodes map[string]*Node) error {

	b, _ := json.Marshal(packageInfo)
	fmt.Println(string(b))

	return nil
}
