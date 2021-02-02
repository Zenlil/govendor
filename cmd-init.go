package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func initVendor(name string) {
	v := load(name, false, nil)
	if v != nil {
		fmt.Println("vendor-file already exists")
		return
	}

	wd, _ := os.Getwd()
	v = &Vendor{
		Ignore:   "test",
		RootPath: strings.TrimPrefix(wd, fmt.Sprintf("%s%c", filepath.Join(args.GoPath, "src"), filepath.Separator)),
		filename: name,
	}
	if args.DryRun {
		fmt.Printf("dry-run: would create '%s' with package = '%s'\n", name, v.RootPath)
	} else {
		os.MkdirAll(filepath.Dir(name), 0)
		v.save()
	}
}
