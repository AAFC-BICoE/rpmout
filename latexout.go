package main

import (
	"bytes"
	"fmt"
	"strings"
)

type LaTeXOut struct {
}

func (lo LaTeXOut) output(header string, dirsOfInterest []string, s []string, rpmInfo map[string]*RpmInfo, groupSet map[string]bool, nodes map[string]*Node) error {
	fmt.Println("\\documentclass[11pt,landscape]{article}")
	fmt.Println("")
	fmt.Println("\\usepackage[landscape,paperwidth=10in,paperheight=8.5in]{geometry}")
	fmt.Println("\\usepackage{longtable,microtype,savetrees}")
	fmt.Println("\\usepackage{fancyhdr}")
	//fmt.Println("\\usepackage[hyphens]{url}")
	fmt.Println("\\usepackage[yyyymmdd,hhmmss]{datetime}")
	fmt.Println("\\usepackage{hyperref}")
	fmt.Println("\\usepackage{seqsplit}")

	fmt.Println("")
	fmt.Println("\\oddsidemargin -.5cm")
	fmt.Println("\\evensidemargin -.5cm")
	fmt.Println("")
	fmt.Println("\\newcommand\\foo[2]{%")
	fmt.Println("\\begin{minipage}{#1}")
	fmt.Println("\\seqsplit{#2}")
	fmt.Println("\\end{minipage}")
	fmt.Println("}")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("\\begin{document}")
	fmt.Println("\\pagestyle{fancy}")

	fmt.Println("\\cfoot{Updated: \\today\\ at \\currenttime}")
	fmt.Println("\\rfoot{\\thepage}")
	fmt.Println("\\lhead{\\bf", header, "}")
	fmt.Println("\\rhead{RPMs in directories: [/] }")

	fmt.Println("%\\thispagestyle{empty}")
	fmt.Println("%\\pagestyle{empty}")
	//fmt.Println("\\tableofcontents")
	//fmt.Println("\\newpage")
	//fmt.Println("\\begin{landscape}")
	fmt.Println("\\renewcommand*{\\arraystretch}{1.4}")
	fmt.Println("\\begin{longtable}{|p{2cm}|p{1.4cm}|p{4cm}|p{5cm}|p{4cm}|p{3cm}|}")
	fmt.Println("\\hline")
	fmt.Println("\\textbf{Name}& \\textbf{Version}& \\textbf{Summary}& \\textbf{Description}& \\textbf{URL}& \\textbf{License}\\\\")
	fmt.Println("\\hline")
	fmt.Println("\\endfirsthead")
	fmt.Println("\\hline")
	fmt.Println("\\textbf{Name}& \\textbf{Version}& \\textbf{Summary}& \\textbf{Description}& \\textbf{URL}& \\textbf{License}\\\\")
	fmt.Println("\\hline")
	fmt.Println("\\endhead")

	//fmt.Println("\\begin{enumerate}")
	for r := range s {
		//fmt.Println("\\section{" + escapeLatex(rpmInfo[s[r]].Name) + "}")
		//fmt.Println("\\item{" + escapeLatex(rpmInfo[s[r]].Name) + "}")
		//fmt.Println("\\begin{itemize}")
		//for k,v := range rpmInfo[s[r]].Tags{

		//	v = strings.Replace(v, "\n", " ", -1)
		//fmt.Println("\\item {\\bf" + escapeLatex("  " + k + ": ") + "}" + escapeLatex(v))
		//fmt.Println("\\newline")
		//	fmt.Println("\\hline")
		name := escapeLatex(insertSpaces(rpmInfo[s[r]].Tags["name"]))

		fmt.Println(name + "&")
		fmt.Println("\\foo{1.4cm}{" + escapeLatex(rpmInfo[s[r]].Tags["version"]) + "}&")
		fmt.Println(escapeLatex(rpmInfo[s[r]].Tags["summary"]) + "&")
		fmt.Println(escapeLatex(rpmInfo[s[r]].Tags["description"]) + "&")
		fmt.Println("\\small \\url{" + escapeLatex(rpmInfo[s[r]].Tags["url"]) + "}&")
		fmt.Println(escapeLatex(rpmInfo[s[r]].Tags["license"]))
		fmt.Println("\\\\ \\hline")
	}
	//fmt.Println("\\end{itemize}")
	//fmt.Println("\\end{section}")
	fmt.Println("\\end{longtable}")
	//fmt.Println("\\end{landscape}")
	fmt.Println("\\end{document}")

	return nil
}

func insertSpacesAtCertainCharacters(v string) string {
	var buffer bytes.Buffer
	minRange := 4
	for i, c := range v {
		buffer.WriteString(string(c))
		if i > 0 && minRange >= 4 && (string(c) == "-" || string(c) == "_") {
			buffer.WriteString(" ")
			minRange = 0
		}

		minRange++
	}
	return buffer.String()
}

// This could be better....break at Capitals or hyphens / underscores
func insertSpaces(v string) string {
	if len(v) < 11 {
		return v
	}

	v = insertSpacesAtCertainCharacters(v)

	var buffer bytes.Buffer
	count := 0
	for i, c := range v {
		if string(c) == " " {
			count = 0
		} else {
			if string(c) != " " && string(c) != "\\" && i > 0 && count >= 8 {
				buffer.WriteString(" ")
				count = 0
			}
		}
		buffer.WriteString(string(c))
		count++
	}
	return buffer.String()
}

func escapeLatex(v string) string {
	v = strings.Replace(v, "\\", "\\textbackslash{}", -1)
	v = strings.Replace(v, "_", "\\_", -1)
	v = strings.Replace(v, "$", "\\$", -1)

	v = strings.Replace(v, "#", "\\#", -1)
	v = strings.Replace(v, "%", "\\%", -1)
	v = strings.Replace(v, "^", "\\^{}", -1)
	v = strings.Replace(v, "&", "\\&", -1)

	v = strings.Replace(v, "{", "\\{", -1)

	v = strings.Replace(v, "}", "\\}", -1)
	v = strings.Replace(v, "~", "\\~{}", -1)

	return v
}
