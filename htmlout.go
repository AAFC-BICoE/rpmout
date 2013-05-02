package main

import(
	"fmt"
	"html"
//"html/template"
//"os"
)

type HtmlOut struct{

}


func (ho HtmlOut) output(s []string, rpmInfo map[string] *RpmInfo) {
	//t, _ := template.New("foo").Parse(`{{define "T"}}Hello, {{.}} how are you!{{end}}`)

	fmt.Println("<ol>")
	for r := range s {	
		//fmt.Println("html " + rpmInfo[s[r]].name)
		//_ = t.ExecuteTemplate(os.Stdout, "T", rpmInfo[s[r]].name)
		rpm := rpmInfo[s[r]]
		fmt.Println("<li>")
		fmt.Println("<b>" + html.EscapeString(rpm.Tags["summary"])+"</b>")
		fmt.Println("<br> V." +  html.EscapeString(rpm.Tags["version"]))
		fmt.Println("<br><b><tt>" + html.EscapeString(rpm.Name) +  "</tt></b>")
		_, ok := rpm.Tags["description"]
		if ok{
			fmt.Println("<br><br><em>" + html.EscapeString(rpm.Tags["description"]) +  "</em>")
			fmt.Println("<br>")
		}
		_, ok = rpm.Tags["url"]
		if ok{
			fmt.Println("<br><b>URL</b>: <a href=\"" + html.EscapeString(rpm.Tags["url"]) + "\">" + html.EscapeString(rpm.Tags["url"])  + "</a>")
		}
		fmt.Println("<br><b>Installed/updated:</b> " +  html.EscapeString(rpm.Tags["installtime"]))
		_, ok = rpm.Tags["packager"]
		if ok{
			fmt.Println("<br><b>Packager:</b> " +  html.EscapeString(rpm.Tags["packager"]))
		}
		fmt.Println("<br><b>License:</b> " +  html.EscapeString(rpm.Tags["license"]))
		fmt.Println("<br><b>Group:</b> " +  html.EscapeString(rpm.Tags["group"]))

		fmt.Println("<br><br>")
		fmt.Println("</li>")
	}
	fmt.Println("</ol>")
	fmt.Println("")
}