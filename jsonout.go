package main

import(
	"fmt"
"encoding/json"
)

type JsonOut struct{

}

func (jo JsonOut) output(s []string, rpmInfo map[string] *RpmInfo, groupSet map[string]bool, nodes map[string]*Node) error{

	b, _ := json.Marshal(rpmInfo)
	fmt.Println(string(b))

	return nil
}