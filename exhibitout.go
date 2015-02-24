package main

import (
	"bufio"
	//	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"text/template"
)

type ExhibitOut struct {
	OutputDir string
}

type TemplateData struct {
	Header string
}

func (lo ExhibitOut) output(outputDir string, outputFileBaseName string, header string, dirsOfInterest []string, sortedKeys []string, packageInfo map[string]*PackageInfo, groupSet map[string]bool, nodes map[string]*Node) error {

	content := TemplateData{header}

	exists, err := exists(outputDir)
	if err != nil {
		return err
	}
	if !exists {
		err = os.MkdirAll(outputDir, 0700)
		if err != nil {
			return err
		}
	}

	einfo := new(EInfo)
	info := make([]*EItem, 0)

	for _, p := range packageInfo {
		item := new(EItem)
		item.Label = p.Name
		if p.IsR {
			item.Type = "R package"
		} else {
			item.Type = "RPM Software"
		}

		item.Group = p.Tags["group"]
		item.License = p.Tags["license"]
		item.Description = p.Tags["description"]
		item.Label = p.Tags["name"]
		item.ShortDescription = p.Tags["summary"]
		item.Url = p.Tags["url"]

		info = append(info, item)
	}
	einfo.Items = info
	//fmt.Printf("%+v\n", info)
	b, _ := json.Marshal(einfo)
	//fmt.Println(string(b))
	err = ioutil.WriteFile("allSoftware.js", b, 0644)
	if err != nil {
		return err
	}

	tmpl, err := template.New("test").Parse(exhibitTemplate)

	// var buffer bytes.Buffer
	// writer := bufio.NewWriter(&buffer)

	// err = tmpl.Execute(writer, content)
	// if err != nil {
	// 	return err
	// }

	// //err = ioutil.WriteFile("rpmout.html", []byte(exhibitTemplate), 0644)
	// err = ioutil.WriteFile("rpmout.html", buffer.Bytes(), 0644)
	// if err != nil {
	// 	return err
	// }

	f, err := os.Create("rpmout.html")
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	err = tmpl.Execute(w, content)
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}

type EItem struct {
	Label            string `json:"label"`
	Group            string
	ShortDescription string `json:"ShortDescription"`
	License          string
	Description      string
	Url              string
	Type             string
}

type EInfo struct {
	Items []*EItem `json:"items"`
}

// From: http://stackoverflow.com/a/10510783
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
