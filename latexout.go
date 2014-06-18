package main

import (
	"fmt"
	"strings"
)

type LaTeXOut struct {
}

func (lo LaTeXOut) output(s []string, rpmInfo map[string]*RpmInfo, groupSet map[string]bool, nodes map[string]*Node) error {
	fmt.Println("\\documentclass[11pt,landscape]{article}")
	fmt.Println("")
	fmt.Println("\\usepackage[landscape,paperwidth=10in,paperheight=8.5in]{geometry}")
	fmt.Println("\\usepackage{longtable,microtype,savetrees}")
	fmt.Println("\\usepackage[hyphens]{url}")
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
	fmt.Println("\\thispagestyle{empty}")
	fmt.Println("\\pagestyle{empty}")
	//fmt.Println("\\tableofcontents")
	//fmt.Println("\\newpage")
	//fmt.Println("\\begin{landscape}")
	fmt.Println("\\renewcommand*{\arraystretch}{1.4}")
	fmt.Println("\\begin{longtable}{|p{2cm}|p{1.4cm}|p{4cm}|p{5cm}|p{4cm}|p{3cm}|}")
	fmt.Println("\\hline")
	fmt.Println("\\textbf{Name}& \\textbf{Version}& \\textbf{Summary}& \\textbf{Description}& \\textbf{URL}& \\textbf{Install Time}\\\\")
	fmt.Println("\\hline")
	fmt.Println("\\endfirsthead")
	fmt.Println("\\hline")
	fmt.Println("\\textbf{Name}& \\textbf{Version}& \\textbf{Summary}& \\textbf{Description}& \\textbf{URL}& \\textbf{Install Time}\\\\")
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
		fmt.Println(escapeLatex(rpmInfo[s[r]].Tags["name"]) + "&")
		fmt.Println("\\foo{1.4cm}{" + escapeLatex(rpmInfo[s[r]].Tags["version"]) + "}&")
		fmt.Println(escapeLatex(rpmInfo[s[r]].Tags["summary"]) + "&")
		fmt.Println(escapeLatex(rpmInfo[s[r]].Tags["description"]) + "&")
		fmt.Println("\\url{" + escapeLatex(rpmInfo[s[r]].Tags["url"]) + "}&")
		fmt.Println(escapeLatex(rpmInfo[s[r]].Tags["installtime"]))
		fmt.Println("\\\\ \\hline")
	}
	//fmt.Println("\\end{itemize}")
	//fmt.Println("\\end{section}")
	fmt.Println("\\end{longtable}")
	//fmt.Println("\\end{landscape}")
	fmt.Println("\\end{document}")

	return nil
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
