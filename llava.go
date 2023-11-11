package main

import (
	"os"
	"os/exec"
	"strings"
)

func llava(prompt string, path string) (string, error) {
	wd, _ := os.Getwd()
	cmd := exec.Command(wd+"/llava/llava-cli",
		"-m", wd+"/llava/ggml-model-q4_k.gguf",
		"--mmproj", wd+"/llava/mmproj-model-f16.gguf",
		"--image", path,
		"-p", prompt)
	cmd.Dir = wd

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	var collect bool = false
	var output []string
	for _, line := range strings.Split(strings.TrimSuffix(string(out), "\n"), "\n") {
		if collect && strings.HasPrefix(line, "main: ") {
			break
		}
		if !collect && strings.HasPrefix(line, " ") {
			collect = true
		}
		if collect {
			output = append(output, line)
		}
	}
	answer := strings.TrimSpace(strings.Join(output, " "))
	return answer, nil
}
