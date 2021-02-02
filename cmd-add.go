package main

import "fmt"

func add(name string, names []string) {
	v := load(name, true, names)
	if v == nil {
		return
	}

	v.AddNames()
}

func (v *Vendor) AddNames() {
	for _, p := range v.Packages {
		if p.filtered {
			fmt.Printf("%s is already vendored\n", p.Path)
		}
	}

	var added int = 0
	for _, arg := range v.unfiltered {
		if arg != "" {
			if v.Add(arg) {
				added++
			}
		}
	}
	if added > 0 {
		if args.DryRun {
			fmt.Printf("%d packages added (not saved -- dry-run)\n", added)
		} else {
			fmt.Printf("%d packages added\n", added)
			v.save()
		}
	}
}

func (v *Vendor) Add(arg string) bool {
	path := locate(arg)
	if path == "" {
		fmt.Printf("Unable to locate '%s'.\n", arg)
		return false
	}

	var p = &Package{
		Path:         arg,
		sourcefolder: path,
		filtered:     true,
	}

	if !p.updateGit() {
		fmt.Printf("Unable to find git-info about '%s'.\n", arg)
		return false
	}

	v.Packages = append(v.Packages, p)
	fmt.Printf("adding '%s'...\n", arg)
	return true
}
