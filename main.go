package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/magefile/mage/mage"
)

const (
	maxSize = 10
)

var (
	mageVersion = "<unknown>"
	timestamp   = "<unknown>"
	hash        = "<unknown>"
	tag         = "<unknown>"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "-version" {
			fmt.Printf("mage-select CLI frontend: %s\nBuild Date: %s\nCommit: %s\nMage: %s\n", tag, timestamp, hash, mageVersion)
			os.Exit(0)
		}
	}

	cmd := exec.Command("mage", "-l")
	cmd.Stderr = os.Stderr

	cmdEnviron := []string{}
	cmdEnviron = append(cmdEnviron, cmd.Environ()...)
	cmdEnviron = append(cmdEnviron, "MAGEFILE_ENABLE_COLOR=0")
	cmd.Env = cmdEnviron

	out, err := cmd.Output()
	if err != nil {
		return
	}

	scan := bufio.NewScanner(bytes.NewBuffer(out))

	targets := []string{}
	for scan.Scan() {
		line := scan.Text()
		if strings.HasPrefix(line, "Targets:") {
			continue
		} else if strings.TrimSpace(line) == "" {
			break
		}
		line = strings.TrimSpace(strings.ReplaceAll(line, "*", ""))
		targets = append(targets, line)
	}

	var options []huh.Option[string]
	for _, target := range targets {
		options = append(options, huh.NewOption(strings.ToUpper(target), target))
	}

	var result string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a mage target").
				Options(options...).
				Value(&result),
		),
	)

	err = form.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	result = strings.Split(result, " ")[0]

	fmt.Printf("mage %s\n", result)

	os.Exit(mage.ParseAndRun(os.Stdout, os.Stderr, os.Stdin, []string{result}))
}
