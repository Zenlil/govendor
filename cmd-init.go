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
		RootPath: strings.TrimPrefix(wd, extras.goSrc),
		filename: name,
	}
	if args.DryRun {
		fmt.Printf("dry-run: would create '%s' with package = '%s'\n", name, v.RootPath)
	} else {
		os.MkdirAll(filepath.Dir(name), defaultAccess)
		v.save()
	}
}
