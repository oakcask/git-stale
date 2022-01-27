package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/oakcask/git-stale/git"
)

type option struct {
	delete   bool
	force    bool
	prefixes []string
}

func getOptions() option {
	deleteLong := flag.Bool("delete", false, "delete stale branches")
	deleteShort := flag.Bool("d", false, "short alias for --delete")
	forceLong := flag.Bool("force", false, "force")
	forceShort := flag.Bool("f", false, "short alias for --force")
	flag.Parse()

	return option{
		delete:   *deleteLong || *deleteShort,
		force:    *forceLong || *forceShort,
		prefixes: flag.Args(),
	}
}

func matchPrefix(s string, prefixes []string) bool {
	if len(prefixes) < 1 {
		return true
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}

	return false
}

type actualGitCommand struct{}

func (c *actualGitCommand) Call(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...) // #nosec G204
	return cmd.Output()
}

func (c *actualGitCommand) Run(args ...string) error {
	cmd := exec.Command("git", args...) // #nosec G204
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	opts := getOptions()

	client, e := git.NewGit(&actualGitCommand{})
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}

	branches, e := client.GetBranches()
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}

	goneBranches := []git.RefName{}
	for _, ref := range branches {
		if ref.Gone && matchPrefix(string(ref.Name), opts.prefixes) {
			goneBranches = append(goneBranches, ref.Name)
		}
	}

	if opts.delete {
		e = client.RemoveBranches(opts.force, goneBranches...)
		if e != nil {
			if _, ok := e.(*exec.ExitError); !ok {
				fmt.Println(e.Error())
			}

			os.Exit(1)
		}
	} else {
		for _, ref := range goneBranches {
			fmt.Println(ref.String())
		}
	}
}
