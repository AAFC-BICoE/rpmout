package main

import (
	"encoding/json"
	"fmt"
)

//type RpmWriter interface {
//	output(string, []string, []string, map[string]*PackageInfo, map[string]bool, map[string]*Node) error
//}

type ExhibitOut struct {
}

func (lo ExhibitOut) output(outputLocation string, header string, dirsOfInterest []string, sortedKeys []string, packageInfo map[string]*PackageInfo, groupSet map[string]bool, nodes map[string]*Node) error {
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

		//fmt.Println(info)
		//fmt.Println(p)
		//fmt.Println(item)
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
	fmt.Println(string(b))
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

// see http://logd.tw.rpi.edu/tutorial/using_mit_simile_exhibit

const htmlTemplate = `
 <html>
 <head>
    <title>Cluster Software</title>
  
    <script src=\"http://trunk.simile-widgets.org/exhibit/api/exhibit-api.js\"></script>
    <link href=\"{{.jsonDataFileName}}.js\"  rel=\"exhibit/data\" />
 </head> 

 <body>
   <h1>Cluster Software</h1>
       <table width=\"100%\">
        <tr valign=\"top\">
            <td ex:role=\"viewPanel\">
	      <div ex:role=\"view\" ex:label=\"List\"></div>
	                  </td>
            <td width=\"25%\">
	      <div ex:role=\"facet\" ex:facetClass=\"TextSearch\"></div>
	      <div ex:role=\"facet\" ex:expression=\".type\" ex:facetLabel=\"Type\"></div>
	      <div ex:role=\"facet\" ex:expression=\".license\" ex:facetLabel=\"License\"></div>
	      <div ex:role=\"facet\" ex:expression=\".group\" ex:facetLabel=\"Group\"></div>
            </td>
        </tr>
    </table>
 </body>
`
