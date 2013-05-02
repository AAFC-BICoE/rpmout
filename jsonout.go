package main

import(
	"fmt"
"encoding/json"
)

type JsonOut struct{

}

func (jo JsonOut) output(s []string, rpmInfo map[string] *RpmInfo) {

	b, _ := json.Marshal(rpmInfo)
	fmt.Println(string(b))
}