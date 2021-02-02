package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func delete(name string, filter []string) {
	v := load(name, false, filter)
	if v == nil {
		return
	}

	for _, name := range v.unfiltered {
		fmt.Printf("%s is not vendored - ignoring\n", name)
	}

	var newList = make([]*Package, 0, len(v.Packages))
	for _, p := range v.Packages {
		if p.filtered {
			if !p.delete() {
				newList = append(newList, p)
			}
		} else {
			newList = append(newList, p)
		}
	}
	if len(v.Packages) != len(newList) {
		if args.DryRun {
			fmt.Printf("%d packages deleted. (not really - dry-run)\n", len(v.Packages)-len(newList))
		} else {
			fmt.Printf("%d packages deleted.\n", len(v.Packages)-len(newList))
			v.Packages = newList
			v.save()
		}
	}
}

func (p *Package) delete() bool {
	fmt.Printf("deleting %s from vendoring...\n", p.Path)
	tDir := filepath.Join(filepath.Dir(args.Filename), p.Path)

	_, err := os.Stat(tDir)
	if err != nil {
		fmt.Printf("unable to delete '%s' - skipping\n", tDir)
		return false
	}
	if args.DryRun {
		fmt.Printf("dry-run: would remove '%s'\n", tDir)
	} else {
		err = os.RemoveAll(tDir)
		if err != nil {
			fmt.Printf("error deleting '%s': %s\n", tDir, err)
			return false
		}
		removeEmptyFolders(filepath.Dir(tDir))
	}

	return true
}

func removeEmptyFolders(name string) {
	f, err := os.Open(name)
	if err != nil {
		return
	}

	names, err := f.Readdirnames(1)
	f.Close()
	if err != nil && err != io.EOF {
		fmt.Printf("error purging %s: %s\n", name, err)
	}
	if len(names) > 0 {
		return
	}

	err = os.Remove(name)
	if err != nil {
		fmt.Printf("error removing %s: %s\n", name, err)
	}

	removeEmptyFolders(filepath.Dir(name))
}
