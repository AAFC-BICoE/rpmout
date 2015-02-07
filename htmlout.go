package main

import (
	"fmt"
	"html"
	"sort"
	//"html/template"
	//"os"
)

type HtmlOut struct {
}

func (ho HtmlOut) output(header string, dirsOfInterest []string, s []string, packageInfo map[string]*PackageInfo, groupSet map[string]bool, nodes map[string]*Node) error {
	//t, _ := template.New("foo").Parse(`{{define "T"}}Hello, {{.}} how are you!{{end}}`)

	index := new(Index)
	index.Init()

	groupSetList := make([]string, len(groupSet))
	i := 0
	for g, _ := range groupSet {
		groupSetList[i] = g
		i++
	}
	sort.Strings(groupSetList)
	for g := range groupSetList {
		fmt.Println("<br>[", groupSetList[g], "]")
	}

	fmt.Println("<ol>")
	for r := range s {
		//fmt.Println("html " + packageInfo[s[r]].name)
		//_ = t.ExecuteTemplate(os.Stdout, "T", packageInfo[s[r]].name)
		rpm := packageInfo[s[r]]
		fmt.Println("<li>")
		fmt.Println("<b>" + html.EscapeString(rpm.Tags["name"]) + "</b>")
		fmt.Println(" - " + html.EscapeString(rpm.Tags["summary"]))
		fmt.Println("<br> V." + html.EscapeString(rpm.Tags["version"]))
		fmt.Println("<br><b><tt>" + html.EscapeString(rpm.Name) + "</tt></b>")
		_, ok := rpm.Tags["description"]
		if ok {
			fmt.Println("<br><br><em>" + html.EscapeString(rpm.Tags["description"]) + "</em>")
			fmt.Println("<br>")
		}
		_, ok = rpm.Tags["url"]
		if ok {
			fmt.Println("<br><b>URL</b>: <a href=\"" + html.EscapeString(rpm.Tags["url"]) + "\">" + html.EscapeString(rpm.Tags["url"]) + "</a>")
		}
		fmt.Println("<br><b>Installed/updated:</b> " + html.EscapeString(rpm.Tags["installtime"]))
		_, ok = rpm.Tags["packager"]
		if ok {
			fmt.Println("<br><b>Packager:</b> " + html.EscapeString(rpm.Tags["packager"]))
		}
		fmt.Println("<br><b>License:</b> " + html.EscapeString(rpm.Tags["license"]))
		if rpm.IsR {
			fmt.Println("<br><b>R package</b> ")
		} else {
			fmt.Println("<br><b>Group:</b> " + html.EscapeString(rpm.Tags["group"]))
		}

		fmt.Println("<br><br>")
		fmt.Println("</li>")
	}
	fmt.Println("</ol>")
	fmt.Println("")

	return nil
}
