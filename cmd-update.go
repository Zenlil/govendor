package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
)

func update(name string, filter []string) {
	v := load(name, true, filter)
	if v == nil {
		return
	}

	v.Update()
}

func (v *Vendor) Update() {

	target := path.Dir(v.filename)
	if target == "" || target == "/" {
		target = "."
	}

	var totSize int64 = 0
	for _, p := range v.Packages {
		if p.filtered {
			fmt.Printf("...scanning...                   %s\r", p.Path)

			if p.scanToCopy(target) {
				p.scanToDelete()
				totSize += p.sumSize
				var toDelete string
				if len(p.toDelete) > 0 {
					toDelete = fmt.Sprintf(" (remove %d files)", len(p.toDelete))
				}
				fmt.Printf("%5d files  %12s bytes  %s%s%s\n", len(p.toCopy), humanize.Comma(p.sumSize), p.Path, toDelete, clearEOL)
			}
		}
	}

	var copiedSize int64 = 0
	var updated bool = false
	for _, p := range v.Packages {
		if p.filtered {
			if args.DryRun {
				fmt.Printf("dry-run: would copy/update %d and remove %d files for %s\n", len(p.toCopy), len(p.toDelete), p.Path)
			} else {
				fmt.Printf("\n%s%s%s\r", p.Path, clearEOL, moveUp)
				for name, size := range p.toCopy {
					fSize, err := copyFile(name, p.sourcefolder, p.target)
					if err != nil {
						fmt.Printf("error copying %s: %s%s\n", filepath.Join(p.Path, name), err, clearEOL)
					} else if fSize != size {
						fmt.Printf("error copying %s: copied %d, expected %d%s\n", filepath.Join(p.Path, name), fSize, size, clearEOL)
					}

					copiedSize += size
					progress(console.w, copiedSize, totSize)
				}

				for name := range p.toDelete {
					err := os.Remove(filepath.Join(p.target, name))
					if err != nil {
						fmt.Printf("error deleting %s: %s%s\n", filepath.Join(p.target, name), err, clearEOL)
					}
				}

				if p.updateGit() {
					updated = true
				}
			}
		}
	}
	if updated {
		v.save()
	}
	fmt.Printf("\nall done%s\n", clearEOL)
}

func progress(w int, n, tot int64) {
	a := int(n * int64(w-margin) / tot)
	b := w - margin - a
	fmt.Printf("[%s%s]\r", strings.Repeat("#", a), strings.Repeat(".", b))
}

func (p *Package) scanToCopy(target string) bool {
	p.target = path.Join(target, p.Path)

	p.sourcefolder = locate(p.Path)
	if p.sourcefolder == "" {
		fmt.Printf("--> source is not available, plese do `go get %s'\n", p.Path)
		return false
	}

	err := os.MkdirAll(p.target, defaultAccess)
	if err != nil {
		fmt.Printf("--> target is not available: %s\n", err)
		return false
	}

	p.toCopy = make(map[string]int64)
	findFilesToCopy(fmt.Sprintf("%s%c", p.sourcefolder, filepath.Separator), p.sourcefolder, p.toCopy)
	for _, size := range p.toCopy {
		p.sumSize += size
	}
	return true
}

func locate(source string) string {
	exppath := path.Join(args.GoPath, "src", source)
	if fi, err := os.Stat(exppath); err != nil || !fi.IsDir() {
		return ""
	}
	return exppath
}

func findFilesToCopy(realSource, source string, copied map[string]int64) {
	var walkFn filepath.WalkFunc = func(path string, fi os.FileInfo, inErr error) error {
		return walkCopyFile(path, fi, inErr, realSource, source, copied)
	}

	err := filepath.Walk(source, walkFn)
	if err != nil {
		panic(err)
	}
}

func walkCopyFile(path string, fi os.FileInfo, _ error, realSource, source string, copied map[string]int64) error {
	name := fi.Name()
	if name == "." || name == ".." || path == source {
		return nil
	}

	relName := strings.TrimPrefix(path, realSource)
	if relName == "" {
		return nil
	}

	if fi.IsDir() {
		if strings.HasPrefix(relName, ".git") {
			return filepath.SkipDir
		}
	} else {
		copied[relName] = fi.Size()
	}
	return nil
}

func copyFile(name, from, to string) (int64, error) {
	source := filepath.Join(from, name)
	target := filepath.Join(to, name)

	targetDir := filepath.Dir(target)
	_ = os.MkdirAll(targetDir, defaultAccess)

	targetF, err := os.Create(target)
	if err != nil {
		return 0, err
	}
	defer targetF.Close()

	sourceF, err := os.Open(source)
	if err != nil {
		return 0, err
	}
	defer sourceF.Close()

	return io.Copy(targetF, sourceF)
}

func (p *Package) scanToDelete() {
	p.toDelete = make(map[string]bool)

	realTarget := fmt.Sprintf("%s%c", p.target, filepath.Separator)

	var walkFn filepath.WalkFunc = func(path string, fi os.FileInfo, inErr error) error {
		return p.walkDeleteFile(path, fi, inErr, realTarget)
	}

	err := filepath.Walk(p.target, walkFn)
	if err != nil {
		panic(err)
	}
}

func (p *Package) walkDeleteFile(path string, fi os.FileInfo, _ error, realTarget string) error {
	name := fi.Name()
	if name == "." || name == ".." || path == p.target {
		return nil
	}

	relName := strings.TrimPrefix(path, realTarget)
	if relName == "" {
		return nil
	}

	if fi.IsDir() {
		if strings.HasPrefix(relName, ".git") {
			return filepath.SkipDir
		}
	} else {
		if _, ok := p.toCopy[relName]; ok {
			return nil
		}
		p.toDelete[relName] = true
	}
	return nil
}
