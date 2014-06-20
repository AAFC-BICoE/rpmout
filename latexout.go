package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type LaTeXOut struct {
}

func (lo LaTeXOut) output(header string, dirsOfInterest []string, s []string, rpmInfo map[string]*RpmInfo, groupSet map[string]bool, nodes map[string]*Node) error {
	fmt.Println("\\documentclass[11pt]{article}")
	fmt.Println("")
	fmt.Println("\\usepackage{longtable,microtype,savetrees}")
	fmt.Println("\\usepackage{fancyhdr}")
	fmt.Println("\\usepackage[yyyymmdd,hhmmss]{datetime}")
	fmt.Println("\\usepackage{hyperref}")
	fmt.Println("\\usepackage{seqsplit}")
	fmt.Println("\\usepackage[usenames,dvipsnames]{color}")
	fmt.Println("\\usepackage{makeidx}")

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

	fmt.Println("\\renewcommand{\\thispagestyle}[1]{}")
	fmt.Println("")
	fmt.Println("\\setlength{\\footskip}{20pt}")
	fmt.Println("\\pagestyle{fancy}")
	fmt.Println("")
	fmt.Println("\\lhead{\\ }")
	fmt.Println("\\chead{\\bf", header, "}")
	fmt.Println("\\rhead{\\hyperlink{rpmIndex}{Index}}")
	fmt.Println("")

	fmt.Println("\\cfoot{Updated: \\today\\ at \\currenttime}")
	fmt.Println("\\rfoot{\\thepage}")
	fmt.Println("")
	fmt.Println("\\makeindex")
	fmt.Println("")
	fmt.Println("\\begin{document}")
	fmt.Println("\\thispagestyle{fancy}")
	fmt.Println("")
	fmt.Println("\\pagestyle{fancy}")
	fmt.Println("{\\newpage\\label{rpmIndex}\\printindex}")
	fmt.Println("")
	fmt.Println("\\renewcommand*{\\arraystretch}{1.4}")
	fmt.Println("\\begin{longtable}{|p{3.5cm}|p{4cm}|p{9.67cm}|}")
	fmt.Println("\\hline")
	fmt.Println("\\textbf{Name}& \\textbf{Summary}& \\textbf{Description}\\\\")
	fmt.Println("\\hline")
	fmt.Println("\\endfirsthead")
	fmt.Println("\\hline")
	fmt.Println("\\textbf{Name}& \\textbf{Summary}& \\textbf{Description}\\\\")
	fmt.Println("\\hline")
	fmt.Println("\\endhead")

	count := 0

	for r := range s {
		count += 1
		name := escapeLatex(rpmInfo[s[r]].Tags["name"])

		fmt.Println("{\\bf \\color{blue}" + name + " \\index{" + name + "}}")
		fmt.Println("")
		fmt.Println("\\vspace{3mm}")

		fmt.Println("Version: ")
		fmt.Println(escapeLatex(rpmInfo[s[r]].Tags["version"]))
		fmt.Println("&")

		fmt.Println(escapeLatex(rpmInfo[s[r]].Tags["summary"]))
		fmt.Println("&")

		description := escapeLatex(rpmInfo[s[r]].Tags["description"])
		if len(description) > 3000 {
			description = description[:3000] + "\\ldots"
		}
		//http: //www.latex-community.org/forum/viewtopic.php?f=51&t=8096
		fmt.Println("\\em ")
		fmt.Println(description)
		fmt.Println("")

		fmt.Println("\\vspace{3mm}")
		url := rpmInfo[s[r]].Tags["url"]
		fmt.Println("\\noindent URL:")
		if len(url) > 0 {
			fmt.Print("{\\bf\\url{" + url + "}}")
		} else {
			fmt.Println("NA")
		}
		fmt.Println("")
		fmt.Println("")
		fmt.Println("\\vspace{3mm}")
		fmt.Println("\\noindent License: {\\bf\\color{Sepia} " + escapeLatex(rpmInfo[s[r]].Tags["license"]) + "}")
		fmt.Println("\\\\ \\hline")
	}
	//fmt.Println("\\end{itemize}")
	//fmt.Println("\\end{section}")
	fmt.Println("\\end{longtable}")
	fmt.Println("\\noindent Total \\# packages: " + strconv.Itoa(count))
	fmt.Println("\\vfill")
	fmt.Println("\\noindent Made with: \\tt \\href{https://github.com/AAFC-MBB/rpmout}{rpmout}")
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
	v = strings.Replace(v, "{", "\\{", -1)
	v = strings.Replace(v, "}", "\\}", -1)
	v = strings.Replace(v, "\\", "\\textbackslash{}", -1)
	v = strings.Replace(v, "_", "\\_", -1)
	v = strings.Replace(v, "$", "\\$", -1)

	v = strings.Replace(v, "#", "\\#", -1)
	v = strings.Replace(v, "%", "\\%", -1)
	v = strings.Replace(v, "^", "\\^{}", -1)
	v = strings.Replace(v, "&", "\\&", -1)

	v = strings.Replace(v, "~", "\\~{}", -1)

	return v
}
