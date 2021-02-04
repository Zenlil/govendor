package main

import (
	"bufio"
	"os/exec"
	"strings"
	"sync"
)

type execStdErr []string

func (e execStdErr) Error() string {
	return strings.Join([]string(e), ", ")
}

func execProc(folder, cmdname string, args ...string) ([]string, error) {
	cmd := exec.Command(cmdname, args...)
	cmd.Dir = folder

	pipe, _ := cmd.StdoutPipe()
	errpipe, _ := cmd.StderrPipe()

	scan := bufio.NewScanner(pipe)
	errscan := bufio.NewScanner(errpipe)
	var wg sync.WaitGroup
	var list []string
	var errlist execStdErr

	wg.Add(2)
	go func() {
		for scan.Scan() {
			line := scan.Text()
			list = append(list, line)
		}
		wg.Done()
	}()

	go func() {
		for errscan.Scan() {
			line := errscan.Text()
			if len(line) > 0 {
				errlist = append(errlist, line)
			}
		}
		wg.Done()
	}()

	err := cmd.Run()
	wg.Wait()

	if len(errlist) > 0 {
		return list, errlist
	}

	if err != nil {
		return nil, err
	}

	return list, nil
}
