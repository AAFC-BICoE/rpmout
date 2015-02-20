package main

import (
	"fmt"
	//"container/list"
	//	"html"
	//"sort"
	//"html/template"
	"io"
	"os"
	"strconv"
	"strings"
)

type HtmlOut2 struct {
}

var alpha = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

func (ho HtmlOut2) output(outputLocation string, outputBaseFileName string, header string, dirsOfInterest []string, s []string, packageInfo map[string]*PackageInfo, groupSet map[string]bool, nodes map[string]*Node) error {
	rpmFileCreated := make(map[string]bool)
	err := makeL0(nodes, rpmFileCreated)
	return err
}

func makeL0(nodes map[string]*Node, rpmFileCreated map[string]bool) error {
	fmt.Println("makeL0")
	file, err := newFile("h0.html")
	if err != nil {
		return err
	}
	_, err = io.WriteString(file, "<ol>\n")
	snodes := sortStringKeyNodeMap(nodes)

	for nodeKey := range snodes {

		//for _, node := range nodes {
		node := nodes[snodes[nodeKey]]
		makeL1(node, rpmFileCreated)
		fmt.Println(node.Name)
		n, err := io.WriteString(file, "<li><a href=\""+makeFileName(node.Name, 1)+"\">"+node.Name+"</a></li>\n")
		if err != nil {
			fmt.Println(n, err)
			return err
		}
	}
	_, err = io.WriteString(file, "</ol>\n")
	file.Close()
	/*
		for _, node = range nodes {

		}
	*/
	return nil
}

func makeL1(node *Node, rpmFileCreated map[string]bool) (int, error) {
	fmt.Println("makeL1")

	filename := makeFileName(node.Name, 1)
	file, err := newFile(filename)
	if err != nil {
		return 0, err
	}
	_, err = io.WriteString(file, "<ol>\n")
	for _, leaf := range node.Children {
		pkgList := leaf.Packages
		n, err := io.WriteString(file, "\n<li><a href=\""+makeFileName(node.Name+"__"+leaf.Name, 2)+"\">"+leaf.Name+"</a> "+strconv.Itoa(len(pkgList))+"</li>")
		if err != nil {
			fmt.Println(n, err)
			return 0, err
		}
		//fmt.Println("\t",leaf.Name, " ", pkgList.Len())
		makeL2(node.Name+"__"+leaf.Name, pkgList, rpmFileCreated, len(pkgList))
	}
	_, err = io.WriteString(file, "</ol>\n")
	file.Close()
	return 5, nil
}

func makeL2(filename string, pkgs map[string]*PackageInfo, rpmFileCreated map[string]bool, numPackages int) error {
	filename = makeFileName(filename, 2)

	file, err := newFile(filename)
	if err != nil {
		return err
	}
	_, err = io.WriteString(file, "<ol>\n")
	fmt.Println("makeL2")

	sortedPackages := sortStringKeyMap(pkgs)
	//alphaUsed, alphaFirst := alphaList(sortedPackages, pkgs)
	alphaUsed, alphaFirst, alphaName := alphaList(sortedPackages, pkgs)

	for i := range alpha {
		if alphaUsed[alpha[i]] {
			_, _ = io.WriteString(file, "<a href=\"#"+alpha[i]+"\">"+alpha[i]+"</a> ")
		} else {
			_, _ = io.WriteString(file, alpha[i]+" ")
		}
	}

	for p := range sortedPackages {
		packageInfo := pkgs[sortedPackages[p]]
		if alphaFirst[packageInfo.Name] {
			_, _ = io.WriteString(file, "\n<br>"+"<a name=\""+alphaName[packageInfo.Name]+"\">"+strings.ToUpper(alphaName[packageInfo.Name]))
		}
		//_,_ = io.WriteString(file, "\n<li><a href=\"" + makeFileName(packageInfo.Name,3) + "\">" + packageInfo.Name + "</a></li>")
		_, _ = io.WriteString(file, "\n<li><a href=\""+makeFileName(packageInfo.Name, 3)+"\">"+packageInfo.Tags["name"]+"</a></li>")
		makeL3(packageInfo, rpmFileCreated)
	}
	/*
		for pkgs != nil {
			packageInfo := pkgs.Value.(*PackageInfo)
			//fmt.Println("\t\t", packageInfo.Name)
			n, err := io.WriteString(file, "\n<li><a href=\"" + makeFileName(packageInfo.Name,3) + "\">" + packageInfo.Name + "</a></li>")
			makeL3(packageInfo, rpmFileCreated)
			if err != nil {
				fmt.Println(n, err)
				return err
			}
			pkgs = pkgs.Next()

		}
	*/
	_, err = io.WriteString(file, "</ol>\n")
	file.Close()
	return nil
}

func makeL3(packageInfo *PackageInfo, rpmFileCreated map[string]bool) error {
	filename := makeFileName(packageInfo.Name, 3)
	if _, ok := rpmFileCreated[filename]; ok {
		return nil
	}
	file, err := newFile(filename)
	if err != nil {
		return err
	}
	_, err = io.WriteString(file, "\n"+packageInfo.Name)
	_, err = io.WriteString(file, "\n<br>"+packageInfo.Tags["group"])
	file.Close()
	return nil
}

func newFile(filename string) (*os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	return f, err
}

func makeFileName(base string, i int) string {
	base = strings.Replace(strings.ToLower(strings.Trim(base, " \t\n")), " ", "_", -1)

	return "h" + strconv.Itoa(i) + "_" + base + ".html"
}

func alphaList(pkgList []string, pkgs map[string]*PackageInfo) (map[string]bool, map[string]bool, map[string]string) {
	alphaUsed := make(map[string]bool)
	alphaFirst := make(map[string]bool)
	alphaName := make(map[string]string)

	for i := range alpha {
		alphaUsed[alpha[i]] = false
	}

	//alphaCount := 0
	for i := range pkgList {
		pkg := pkgs[pkgList[i]]
		for j := range alpha {
			if !alphaUsed[alpha[j]] {
				if strings.HasPrefix(strings.ToLower(pkg.Name), alpha[j]) {
					alphaUsed[alpha[j]] = true
					alphaFirst[pkg.Name] = true
					alphaName[pkg.Name] = alpha[j]
				}
			}
		}
		/*
			if strings.HasPrefix(strings.ToLower(pkg.Name), alpha[alphaCount]){
				alphaUsed[alpha[alphaCount]] = true
				alphaFirst[alpha[alphaCount]] = pkg.Name
				alphaCount++
			}
		*/
	}
	return alphaUsed, alphaFirst, alphaName
}
