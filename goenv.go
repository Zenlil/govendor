package main

import "strings"

var (
	goEnvMap map[string]string
)

func goEnv(key string) string {
	if goEnvMap == nil {
		lines, err := execProc(".", "go", "env")
		if err != nil {
			return ""
		}
		goEnvMap = make(map[string]string, len(lines))
		for _, line := range lines {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				value := parts[1]
				if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
					value = value[1 : len(value)-1]
				}
				goEnvMap[parts[0]] = value
			}
		}
	}

	return goEnvMap[key]
}
