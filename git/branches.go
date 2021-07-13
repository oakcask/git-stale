package git

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"time"
)

func (git *git) GetBranches() ([]Branch, error) {
	// HACK: :track,nobracket will return string like "ahead 1, behind 2"
	// so using ", " as delimiter between refname and upstream enables us to
	// split these parts by ", ".
	out, e := git.command.Call("branch", "--format", "%(refname:short), %(committerdate), %(upstream:track,nobracket)")
	if e != nil {
		return nil, e
	}

	var branches []Branch

	buffer := bytes.NewBuffer(out)
	scanner := bufio.NewScanner(buffer)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fields := strings.SplitN(scanner.Text(), ", ", 4)
		if len(fields) < 1 {
			return nil, fmt.Errorf("unexpected text returned from git")
		}

		commitDate, e := time.Parse("Mon Jan 2 15:04:05 2006 -0700", fields[1])
		if e != nil {
			return nil, fmt.Errorf("unexpected date time format returned from git: %v", fields[1])
		}

		gone := len(fields) == 3 && fields[2] == "gone"

		branches = append(branches,
			Branch{
				Name:       RefName(fields[0]),
				Gone:       gone,
				CommitDate: commitDate,
			})
	}

	return branches, nil
}

func (git *git) RemoveBranches(force bool, branches ...RefName) error {
	if len(branches) == 0 {
		return nil
	}

	args := []string{"branch", "-d"}
	if force {
		args = append(args, "-f")
	}
	for _, ref := range branches {
		args = append(args, ref.String())
	}

	return git.command.Run(args...)
}
