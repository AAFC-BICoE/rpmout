package main

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type LaTeXOut struct {
}

const IndexSoftware = "Index, Installed Software"
const IndexLicenses = "Index, Licenses"
const IndexGroup = "Index, Group"

func (lo LaTeXOut) output(outputLocation string, outputBaseFileName string, header string, dirsOfInterest []string, sortedKeys []string, packageInfo map[string]*PackageInfo, groupSet map[string]bool, nodes map[string]*Node) error {
	fmt.Println("\\documentclass[11pt]{article}")
	fmt.Println("")
	fmt.Println("\\usepackage{longtable,microtype,savetrees}")
	fmt.Println("\\usepackage{fancyhdr}")
	fmt.Println("\\usepackage[yyyymmdd,hhmmss]{datetime}")
	fmt.Println("\\usepackage{hyperref}")
	fmt.Println("\\usepackage{seqsplit}")
	fmt.Println("\\usepackage[usenames,dvipsnames]{color}")
	fmt.Println("\\usepackage[makeindex]{splitidx}")

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
	//	fmt.Println("\\lhead{\\ }")
	//	fmt.Println("\\chead{\\bf", header, "}")
	fmt.Println("\\rhead{\\hyperlink{toc}{TOC}}")
	fmt.Println("")

	fmt.Println("\\cfoot{Updated: \\today\\ at \\currenttime}")
	fmt.Println("\\rfoot{\\thepage}")
	fmt.Println("")

	fmt.Println("\\newindex[" + IndexSoftware + "]{s}")
	fmt.Println("\\newindex[" + IndexLicenses + "]{l}")
	fmt.Println("\\newindex[" + IndexGroup + "]{g}")

	fmt.Println("")

	fmt.Println("\\title{Installed Software Report}")

	fmt.Println("\\begin{document}")
	fmt.Println("\\thispagestyle{fancy}")
	fmt.Println("")
	fmt.Println("\\pagestyle{fancy}")

	fmt.Println("\\maketitle")

	fmt.Println("\\tableofcontents")
	fmt.Println("\\label{toc}")

	//fmt.Println("{\\newpage\\label{rpmIndex}\\printindex[s]}")

	count, rCount := makeSoftwareTable(sortedKeys, packageInfo)
	makeStatistics(count, rCount, countLicenses(packageInfo))

	makeSoftwareIndex()
	makeGroupIndex()
	makeLicenseIndex()

	fmt.Println("\\end{document}")

	return nil
}

func countLicenses(packageInfo map[string]*PackageInfo) map[string]int {
	counts := make(map[string]int)

	for _, p := range packageInfo {
		license := p.Tags["license"]
		if _, ok := counts[license]; ok {
			//do something here
			counts[license] = 1 + counts[license]
		} else {
			counts[license] = 1
		}
	}
	return counts
}

func makeStatistics(count, rCount int, licenseCount map[string]int) {
	fmt.Println("\\newpage")
	fmt.Println("\\section{Statistics}")

	fmt.Println("\\noindent Total \\# packages: " + strconv.Itoa(count))
	fmt.Println("\\newline")
	fmt.Println("\\noindent Total  \\# {\\bf R} packages: " + strconv.Itoa(rCount))
	fmt.Println("\\vfill")
	fmt.Println("\\noindent {Made with: \\tt \\href{https://github.com/AAFC-MBB/rpmout}{rpmout}}")
	//fmt.Println("\\end{landscape}")

	// License count
	fmt.Println("\\newpage")
	fmt.Println("\\subsection{Licenses}")
	fmt.Println("\\begin{longtable}{|p{10cm}|c|}")
	fmt.Println("\\hline")
	fmt.Println("\\hline")
	fmt.Println("\\textbf{License}& \\textbf{Count}\\\\")
	fmt.Println("\\hline")
	fmt.Println("\\hline")
	fmt.Println("\\endfirsthead")
	fmt.Println("\\hline")
	fmt.Println("\\hline")
	fmt.Println("\\textbf{License}& \\textbf{Count}\\\\")
	fmt.Println("\\hline")
	fmt.Println("\\hline")
	fmt.Println("\\endhead")

	//sortedLicenses := sortMapByValue(licenseCount)

	sortedLicenses := sortedKeys(licenseCount)

	//for k, v := range licenseCount {
	for _, v := range sortedLicenses {
		//fmt.Println("\\newline")
		//v := licenseCount[k]
		//fmt.Println(escapeLatex(pair.Key) + "& " + strconv.Itoa(pair.Value) + "\\\\")
		fullLicense := escapeLatex(strings.TrimSpace(v))
		fmt.Println(fullLicense + "& " + strconv.Itoa(licenseCount[v]) + "\\\\")
		fmt.Println("\\hline")
	}
	fmt.Println("\\end{longtable}")
	fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
}

func makeLicenseIndex() {
	fmt.Println("\\phantomsection")
	fmt.Println("\\cleardoublepage")
	fmt.Println("\\addcontentsline{toc}{section}{" + IndexLicenses + "}")
	fmt.Println("\\printindex[l]")
	//fmt.Println("{\\newpage\\label{rpmIndex}\\printindex[l]}")
	fmt.Println("")
}

func makeSoftwareIndex() {
	fmt.Println("\\phantomsection")
	fmt.Println("\\cleardoublepage")
	fmt.Println("\\addcontentsline{toc}{section}{" + IndexSoftware + "}")
	//%{\newpage\label{rpmIndex}\printindex[s]}
	fmt.Println("\\label{rpmIndex}")
	fmt.Println("\\printindex[s]")
}

func makeGroupIndex() {
	fmt.Println("\\phantomsection")
	fmt.Println("\\cleardoublepage")
	fmt.Println("\\addcontentsline{toc}{section}{" + IndexGroup + "}")
	//%{\newpage\label{rpmIndex}\printindex[s]}
	fmt.Println("\\printindex[g]")
}

func makeSoftwareTable(sortedKeys []string, packageInfo map[string]*PackageInfo) (int, int) {
	fmt.Println("\\newpage")
	fmt.Println("\\section{Installed Software}")
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
	rCount := 0

	for r := range sortedKeys {
		count += 1
		tpackage := packageInfo[sortedKeys[r]]
		if tpackage.IsR {
			rCount += 1
		}
		name := escapeLatex(tpackage.Tags["name"])
		var indexName string
		if tpackage.IndexName == "" {
			indexName = name
		} else {
			indexName = tpackage.IndexName
		}

		fmt.Println("{\\bf \\color{blue}" + name + " \\sindex[s]{" + indexName + "}")
		if tpackage.IsR {
			fmt.Println("\\sindex[s]{" + name + "\\ {\\em (R-package)}}")
		}

		fmt.Println("}")

		fmt.Println("")
		fmt.Println("\\vspace{3mm}")
		fmt.Println("Version: ")
		fmt.Println(escapeLatex(tpackage.Tags["version"]))

		if tpackage.IsR {
			fmt.Println("")
			fmt.Println("\\vspace{3mm}")
			fmt.Println("{ \\bf\\color{Sepia} R package}")
		}

		group := strings.TrimSpace(tpackage.Tags["group"])
		if len(group) > 0 {
			fmt.Println("\\sindex[g]{" + strings.Replace(group, "/", "!", -1) + "!" + name + "}")
			group = escapeLatex(group)
			fmt.Println("")
			fmt.Println("\\vspace{3mm}")
			fmt.Println("{\\small \\bf RPM Group: \\newline \\color{Sepia}" + strings.Replace(group, "/", "\\newline $\\Rightarrow$", -1) + "}")
		}

		fmt.Println("&")

		fmt.Println(escapeLatex(tpackage.Tags["summary"]))
		fmt.Println("&")

		description := escapeLatex(tpackage.Tags["description"])
		if len(description) > 3000 {
			description = description[:3000] + "\\ldots"
		}
		//http: //www.latex-community.org/forum/viewtopic.php?f=51&t=8096
		fmt.Println("\\em ")
		fmt.Println(description)
		fmt.Println("")

		fmt.Println("\\vspace{3mm}")
		url := tpackage.Tags["url"]
		fmt.Println("\\noindent URL:")
		if len(url) > 0 {
			fmt.Print("{\\bf\\url{" + url + "}}")
		} else {
			fmt.Println("NA")
		}
		fmt.Println("")
		fmt.Println("")
		fmt.Println("\\vspace{3mm}")
		fullLicense := strings.TrimSpace(tpackage.Tags["license"])
		license := fullLicense

		fmt.Println("%% license raw:[" + license + "]")
		licenseLength := len(license)
		if licenseLength > 45 {
			license = license[0:44]
		}

		license = escapeLatex(license)

		fmt.Println("\\noindent License: {\\bf\\color{Sepia} " + license)

		fmt.Println("\\sindex[l]{" + escapeLatex(fullLicense))
		if licenseLength > 45 {
			fmt.Print("\\ldots!")
		}
		fmt.Println("!" + name + "}")

		fmt.Println("")
		fmt.Println("}")
		fmt.Println("\\\\ \\hline")
	}
	//fmt.Println("\\end{itemize}")
	//fmt.Println("\\end{section}")

	fmt.Println("\\end{longtable}")
	fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
	return count, rCount
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

	v = strings.Replace(v, ">", "{\\textgreater}", -1)
	v = strings.Replace(v, "<", "{\\textless}", -1)

	v = strings.Replace(v, "#", "\\#", -1)
	v = strings.Replace(v, "%", "\\%", -1)
	v = strings.Replace(v, "^", "{\\textasciicircum}", -1)
	v = strings.Replace(v, "&", "\\&", -1)
	v = strings.Replace(v, "~", "{\\textasciitilde}", -1)

	return v
}

func makePairList(licenseCount map[string]int) PairList {
	pl := make(PairList, len(licenseCount))

	i := 0
	for k, v := range licenseCount {
		p := new(Pair)
		p.Key = k
		p.Value = v
		pl[i] = *p
	}

	return pl
}

type sortedMap struct {
	m map[string]int
	s []string
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func sortedKeys(m map[string]int) []string {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}
