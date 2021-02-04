package main

import (
	"fmt"
	"os"
)

func main() {

	if fi, err := os.Stat(args.GoPath); err != nil || !fi.IsDir() {
		fmt.Printf("error: the GOPATH '%s' is not an existing directory\n", args.GoPath)
		os.Exit(1)
	}

	if args.Init == nil {
		if fi, err := os.Stat(args.Filename); err != nil || fi.IsDir() {
			fmt.Printf("error: the vendorfile '%s' is not an existing file\n", args.Filename)
			os.Exit(1)
		}
	}

	switch true {
	case args.Init != nil:
		initVendor(args.Filename)

	case args.List != nil:
		list(args.Filename, args.List.Names)

	case args.Add != nil:
		add(args.Filename, args.Add.Names)

	case args.Delete != nil:
		delete(args.Filename, args.Delete.Names)

	case args.Update != nil:
		update(args.Filename, args.Update.Names)

	case args.Get != nil:
		get(args.Filename, args.Get.Names)

	case args.Deps != nil:
		deps(args.Filename, false, args.Deps.Dump)

	case args.Tidy != nil:
		deps(args.Filename, true, args.Tidy.Dump)
	}
}
