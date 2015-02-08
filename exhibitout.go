package main

//type RpmWriter interface {
//	output(string, []string, []string, map[string]*PackageInfo, map[string]bool, map[string]*Node) error
//}

type ExhibitOut struct {
}

func (lo ExhibitOut) output(outputLocation string, header string, dirsOfInterest []string, sortedKeys []string, packageInfo map[string]*PackageInfo, groupSet map[string]bool, nodes map[string]*Node) error {
	return nil
}
