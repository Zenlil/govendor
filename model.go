package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Vendor struct {
	Comment  string     `json:"comment"`
	Ignore   string     `json:"ignore"`
	Packages []*Package `json:"package,omitempty"`
	RootPath string     `json:"rootPath"`

	filename   string
	unfiltered []string
}

type Package struct {
	//	ChecksumSHA1 []byte    `json:"checksumSHA1,omitempty"`
	Path         string    `json:"path"`
	Revision     string    `json:"revision"`
	RevisionTime time.Time `json:"revisionTime"`

	filtered bool

	sourcefolder string
	target       string
	toCopy       map[string]int64
	toDelete     map[string]bool
	sumSize      int64
}

func load(name string, logErr bool, filter []string) *Vendor {
	var v = new(Vendor)

	buf, err := ioutil.ReadFile(name)
	if err != nil {
		if logErr {
			fmt.Printf("error reading '%s': %s\n", name, err)
		}
		return nil
	}

	err = json.Unmarshal(buf, v)
	if err != nil {
		if logErr {
			fmt.Printf("error parsing '%s': %s\n", name, err)
		}
		return nil
	}

	//var any = false
	for _, name := range filter {
		var used bool = false
		//any = true
		for _, p := range v.Packages {
			if strings.EqualFold(name, p.Path) {
				used = true
				p.filtered = true
			}
		}
		if !used {
			v.unfiltered = append(v.unfiltered, name)
		}
	}
	if len(filter) == 0 {
		for _, p := range v.Packages {
			p.filtered = true
		}
	}

	v.filename = name

	return v
}

func (v *Vendor) save() {
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("error marshalling '%s': %s\n", v.filename, err)
		return
	}

	f, err := os.Create(v.filename)
	if err != nil {
		fmt.Printf("error writing '%s': %s\n", v.filename, err)
		return
	}
	defer f.Close()

	_, _ = f.Write(buf)
}
