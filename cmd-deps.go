package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	depsMatch = regexp.MustCompile(`^[a-zA-Z0-9]+\.`)
)

type depStatus int

const (
	dsNotUsed depStatus = iota
	dsVendored
	dsAdd
	dsUpdate
	dsDelete
	dsSubvendored
)

func deps(name string, tidy bool, noRemove bool, dump bool) {
	v := load(name, true, nil)
	if v == nil {
		return
	}

	changed, all, err := v.getDeps(args.Source, dump)
	if err != nil {
		fmt.Printf("error getting dependencies: %s\n", err)
		return
	}

	var current = make(map[string]depStatus)
	for _, p := range v.Packages {
		current[p.Path] = dsNotUsed
	}

	var toRemove []string

	for _, name := range changed {
		_, found := current[name]
		if !found {
			current[name] = dsAdd
		} else {
			current[name] = dsUpdate
		}
	}

	for _, name := range all {
		used, found := current[name]
		if !found {
			if v.hasDepParent(name) {
				current[name] = dsSubvendored
			} else {
				fmt.Printf("ERROR: Pkg %s !found but should be in 'changed'\n", name)
			}
		} else {
			if used == dsNotUsed {
				current[name] = dsVendored
			}
		}
	}

	var allNames = make([]string, 0, len(current))
	for name, used := range current {
		allNames = append(allNames, name)
		if used == dsNotUsed {
			toRemove = append(toRemove, name)
		}
	}

	sort.SliceStable(allNames, func(a, b int) bool { return allNames[a] < allNames[b] })

	for _, name := range allNames {
		switch current[name] {
		case dsNotUsed:
			fmt.Printf("??????  %s\n", name)
		case dsAdd:
			fmt.Printf("ADD     %s\n", name)
		case dsUpdate:
			fmt.Printf("UPDATE  %s\n", name)
		case dsDelete:
			if noRemove {
				fmt.Printf("unused? %s\n", name)

			} else {
				fmt.Printf("DELETE  %s\n", name)
			}
		case dsVendored:
			fmt.Printf("ok      %s\n", name)
		case dsSubvendored:
			fmt.Printf("(impl)  %s\n", name)
		}
	}
	var empty = false

	if !args.DryRun && tidy {
		if len(changed) > 0 {
			if !empty {
				fmt.Println()
				empty = true
			}
			v = load(name, true, changed)
			v.AddNames(false)
			v.Update()
		}

		if len(toRemove) > 0 && !noRemove {
			if !empty {
				fmt.Println()
			}
			delete(name, toRemove)
		}
	}
}

type goDepModel struct {
	Dir            string            `json:"Dir"`
	ImportPath     string            `json:"ImportPath"`
	RealImportPath string            `json:"RealImportPath"`
	Name           string            `json:"Name"`
	Target         string            `json:"Target"`
	Root           string            `json:"Root"`
	Imports        []string          `json:"Imports"`
	ImportMap      map[string]string `json:"ImportMap"`
	Deps           []string          `json:"Deps"`
}

func (v *Vendor) getDeps(folder string, doDump bool) ([]string, []string, error) {
	lines, err := execProc(folder, "go", "list", "-deps", "-json")
	if err != nil {
		return nil, nil, err
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteRune('[')
	var first bool = true
	for _, line := range lines {
		if line == "{" {
			if !first {
				buf.WriteRune(',')
			}
			first = false
		}
		buf.WriteString(line)
	}
	buf.WriteRune(']')

	var depModel []goDepModel
	err = json.Unmarshal(buf.Bytes(), &depModel)
	if err != nil {
		return nil, nil, err
	}

	var localVendor = fmt.Sprintf("%s%c", filepath.Join(v.RootPath, "vendor"), filepath.Separator)

	for i := range depModel {
		if strings.HasPrefix(depModel[i].ImportPath, localVendor) {
			depModel[i].RealImportPath = strings.TrimPrefix(depModel[i].ImportPath, localVendor)
		} else {
			depModel[i].RealImportPath = depModel[i].ImportPath
		}
	}

	sort.SliceStable(depModel, func(a, b int) bool { return depModel[a].RealImportPath < depModel[b].RealImportPath })

	var goRoot = goEnv("GOROOT")

	var dump []goDepModel
	var toCheck []string
	var all []string
	for _, pkg := range depModel {
		if strings.HasPrefix(pkg.Dir, goRoot) {
			// ignoring system packages
		} else {
			switch true {
			case strings.HasPrefix(pkg.ImportPath, localVendor):
				// package is locally vendored
				pkgName := strings.TrimPrefix(pkg.ImportPath, localVendor)
				all = append(all, pkgName)
				dump = append(dump, pkg)

			case pkg.ImportPath == v.RootPath:
				// ignore the self-reference

			default:
				toCheck = append(toCheck, pkg.ImportPath)
				all = append(all, pkg.ImportPath)
				dump = append(dump, pkg)
			}
		}
	}

	if doDump {
		buf2, _ := json.MarshalIndent(&dump, "", "  ")
		_ = ioutil.WriteFile("deps.json", buf2, 0777)
	}

	return toCheck, all, nil
}

func (v *Vendor) hasDepParent(name string) bool {
	parent := filepath.Dir(name)
	if len(parent) < 2 || !strings.ContainsAny(parent, string(filepath.Separator)) {
		return false
	}

	for _, p := range v.Packages {
		if parent == p.Path {
			return true
		}
	}
	return v.hasDepParent(parent)
}
