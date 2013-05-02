package main

import(
	"fmt"
	"strings"
)

type TextOut struct{

}

func (to TextOut) output(s []string, rpmInfo map[string] *RpmInfo) {

	for r := range s {	
		fmt.Println("")
		fmt.Println(rpmInfo[s[r]].Name)
		for k,v := range rpmInfo[s[r]].Tags{
			v = strings.Replace(v, "\n", "\n               ", -1)
			fmt.Println("  " + k + ": " + v)
		}
	}

}