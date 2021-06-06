package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/oakcask/git-stale/git"
)

type option struct {
	delete bool
	force  bool
}

func getOptions() option {
	deleteLong := flag.Bool("-delete", false, "delete stale branches")
	deleteShort := flag.Bool("d", false, "short alias for --delete")
	forceLong := flag.Bool("-force", false, "force")
	forceShort := flag.Bool("f", false, "short alias for --force")
	flag.Parse()

	return option{
		delete: *deleteLong || *deleteShort,
		force:  *forceLong || *forceShort,
	}
}

type actualGitCommand struct{}

func (c *actualGitCommand) Call(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	return cmd.Output()
}

func (c *actualGitCommand) Run(args ...string) error {
	cmd := exec.Command("git", args...)
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
		if ref.Gone {
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
