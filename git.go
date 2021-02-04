package main

import (
	"fmt"
	"strings"
	"time"
)

func (p *Package) updateGit() bool {
	var head string

	lines, err := execProc(p.sourcefolder, "git", "rev-parse", "HEAD")
	if err != nil {
		fmt.Printf("error updating git-info for %s: %s\n", p.Path, err)
		return false
	}
	if len(lines) > 0 {
		head = lines[0]
	}

	lines, err = execProc(p.sourcefolder, "git", "rev-list", "--format=format:'%ci'", "--max-count=1", head)
	if err != nil {
		fmt.Printf("error updating git-info for %s: %s\n", p.Path, err)
		return false
	}

	for _, line := range lines {
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
