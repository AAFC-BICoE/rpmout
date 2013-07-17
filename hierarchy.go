
package main

import(
	"strings"
//"container/list"
//"fmt"
)


type Leaf struct{
	Name string
	//Packages *list.List
	Packages map[string]*RpmInfo
}

type Node struct{
	Name string
	Children map[string]*Leaf
}


func extractHierarchy(r *RpmInfo, nodes map[string]*Node){
	if group,ok := r.Tags["group"]; ok {
		group := strings.TrimSpace(group)
		if len(group) > 0{
			parts := strings.Split(group, "/")
			node := getNode(nodes, parts[0])
			//fmt.Println("++++++++++++++++++++++++++ ", group, " ", node)
			if len(parts) > 1 {
				//fmt.Println("++++++++++++++++++++++++++ ", group, " ", node)
				if node.Children == nil {
					//node.Children = make(map[string]Leaf)
					node.Children = make(map[string]*Leaf)
				}
				var leaf *Leaf

				if leaf,ok = node.Children[parts[1]]; ok {
					//
				}else{
					leaf = new(Leaf)
					leaf.Name = parts[1]
					leaf.Packages = make(map[string]*RpmInfo)
					node.Children[parts[1]] = leaf
				}
				leaf.Packages[strings.ToLower(r.Name)] = r
			}
		}
	}
}

func getNode(nodes map[string]*Node, part string) *Node{
	var node *Node
	var ok bool
	if node, ok = nodes[part]; ok{
		
	}else{
		node = new(Node)
		nodes[part] = node
	}
	node.Name = part
	return node
}


