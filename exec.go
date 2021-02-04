package main

import (
	"bufio"
	"os/exec"
	"sync"
)

func execProc(folder, cmdname string, args ...string) ([]string, error) {
	cmd := exec.Command(cmdname, args...)
	cmd.Dir = folder

	pipe, _ := cmd.StdoutPipe()

	scan := bufio.NewScanner(pipe)
	var wg sync.WaitGroup
	var list []string

	wg.Add(1)
	go func() {
		for scan.Scan() {
			line := scan.Text()
			list = append(list, line)
		}
		wg.Done()
	}()

	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	wg.Wait()

	return list, nil
}
