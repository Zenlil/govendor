package main

import (
	"fmt"
	"strings"
	"time"
)

func list(name string, filter []string) {
	v := load(name, true, filter)
	if v == nil {
		return
	}

	fmt.Printf("Vendor-list for '%s'\n", v.RootPath)
	var l = 0
	for _, p := range v.Packages {
		if p.filtered {
			if len(p.Path) > l {
				l = len(p.Path)
			}
		}
	}
	spaces := strings.Repeat(" ", l)

	for _, p := range v.Packages {
		if p.filtered {
			path := (p.Path + spaces)[:l]
			fmt.Printf("%s  %s\n", path, p.RevisionTime.Format(time.RFC3339))
		}
	}
}
