package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func (p *Package) updateGit() bool {
	var head string

	git1 := exec.Command("git", "rev-parse", "HEAD")
	git1.Dir = p.sourcefolder

	out1, _ := git1.StdoutPipe()
	_ = git1.Start()
	_ = git1.Wait

	scan := bufio.NewScanner(out1)
	if scan.Scan() {
		head = scan.Text()
	}

	git := exec.Command("git", "rev-list", "--format=format:'%ci'", "--max-count=1", head)
	git.Dir = p.sourcefolder

	out, err := git.StdoutPipe()
	if err != nil {
		fmt.Printf("error updating git-info for %s: %s\n", p.Path, err)
		return false
	}
	defer out.Close()

	outScanner := bufio.NewScanner(out)

	go func() {
		err = git.Run()
		if err != nil {
			fmt.Printf("error updating git-info for %s: %s\n", p.Path, err)
		}
	}()

	for outScanner.Scan() {
		line := outScanner.Text()
		if strings.HasPrefix(line, "commit ") {
			p.Revision = strings.TrimPrefix(line, "commit ")
		} else {
			line = strings.Trim(line, "'")
			t, err := time.Parse("2006-01-02 15:04:05 -0700", line)
			if err != nil {
				fmt.Printf("error parsing RevisionTime '%s': %s\n\n", line, err)
			} else {
				p.RevisionTime = t.UTC()
			}
		}
	}

	return true
}
