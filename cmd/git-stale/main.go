package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/oakcask/git-stale/git"
	"github.com/oakcask/go-since"
)

type option struct {
	delete   bool
	force    bool
	prefixes []string
	since    *time.Time
	push     bool
}

func parseSince(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}

	dateString := strings.TrimSpace(s)

	t, e := time.Parse("2006-01-02", dateString)
	if e == nil {
		return &t, nil
	}

	t, e = since.Since(dateString, time.Now())
	if e == nil {
		return &t, nil
	}

	return nil, fmt.Errorf("unrecognized date: %v", dateString)
}

func getOptions() (option, error) {
	deleteLong := flag.Bool("delete", false, "delete stale branches")
	deleteShort := flag.Bool("d", false, "short alias for --delete")
	forceLong := flag.Bool("force", false, "force")
	forceShort := flag.Bool("f", false, "short alias for --force")
	since := flag.String("since", "", "select branches by last commit date")
	push := flag.Bool("push", false, "with -d option, execute `git push --delete` to remove remote branch")
	flag.Parse()

	sinceValue, err := parseSince(*since)
	if err != nil {
		return option{}, err
	}

	return option{
		delete:   *deleteLong || *deleteShort,
		force:    *forceLong || *forceShort,
		prefixes: flag.Args(),
		since:    sinceValue,
		push:     *push,
	}, nil
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
	opts, e := getOptions()
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}

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

	goneBranches := []git.Branch{}
	if opts.since == nil {
		for _, ref := range branches {
			if ref.Gone && matchPrefix(string(ref.Name), opts.prefixes) {
				goneBranches = append(goneBranches, ref)
			}
		}
	} else {
		for _, ref := range branches {
			if ref.CommitDate.Sub(*opts.since) < 0 && matchPrefix(string(ref.Name), opts.prefixes) {
				goneBranches = append(goneBranches, ref)
			}
		}
	}

	if opts.delete && opts.push {
		e = client.RemoveRemoteBranches(opts.force, goneBranches...)
		if e != nil {
			if _, ok := e.(*exec.ExitError); !ok {
				fmt.Println(e.Error())
			}

			os.Exit(1)
		}
	} else if opts.delete {
		e = client.RemoveBranches(opts.force, goneBranches...)
		if e != nil {
			if _, ok := e.(*exec.ExitError); !ok {
				fmt.Println(e.Error())
			}

			os.Exit(1)
		}
	} else {
		for _, ref := range goneBranches {
			fmt.Println(ref.Name.String())
		}
	}
}
