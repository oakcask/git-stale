package git

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func (git *git) GetBranches() ([]Branch, error) {
	// HACK: :track,nobracket will return string like "ahead 1, behind 2"
	// so using ", " as delimiter between refname and upstream enables us to
	// split these parts by ", ".
	out, e := git.command.Call("branch", "--format", "%(refname:short), %(upstream:track,nobracket)")
	if e != nil {
		return nil, e
	}

	var branches []Branch

	buffer := bytes.NewBuffer(out)
	scanner := bufio.NewScanner(buffer)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fields := strings.SplitN(scanner.Text(), ", ", 3)
		if len(fields) < 1 {
			return nil, fmt.Errorf("unexpected text returned from git")
		}

		gone := len(fields) == 2 && fields[1] == "gone"

		branches = append(branches,
			Branch{
				Name: RefName(fields[0]),
				Gone: gone,
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
