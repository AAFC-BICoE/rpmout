package main

import (
	"encoding/json"
	"fmt"
)

type JsonOut struct {
}

func (jo JsonOut) output(header string, dirsOfInterest []string, s []string, rpmInfo map[string]*RpmInfo, groupSet map[string]bool, nodes map[string]*Node) error {

	b, _ := json.Marshal(rpmInfo)
	fmt.Println(string(b))

	return nil
}
