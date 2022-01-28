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
	out, e := git.command.Call("branch", "--format", "%(refname:short), %(committerdate), %(upstream:remotename), %(upstream:track,nobracket)")
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

		remoteName := fields[2]

		gone := len(fields) == 4 && fields[3] == "gone"

		branches = append(branches,
			Branch{
				Name:       RefName(fields[0]),
				Gone:       gone,
				CommitDate: commitDate,
				RemoteName: remoteName,
			})
	}

	return branches, nil
}

func (git *git) RemoveBranches(force bool, branches ...Branch) error {
	if len(branches) == 0 {
		return nil
	}

	args := []string{"branch", "-d"}
	if force {
		args = append(args, "-f")
	}
	for _, branch := range branches {
		args = append(args, branch.Name.String())
	}

	return git.command.Run(args...)
}

func (git *git) RemoveRemoteBranches(force bool, branches ...Branch) error {
	if len(branches) == 0 {
		return nil
	}

	branchesWillBePush := make(map[string][]RefName)
	remotes := []string{}

	for _, branch := range branches {
		if branch.RemoteName != "" {
			if len(branchesWillBePush[branch.RemoteName]) == 0 {
				remotes = append(remotes, branch.RemoteName)
			}
			branchesWillBePush[branch.RemoteName] = append(branchesWillBePush[branch.RemoteName], branch.Name)
		}
	}

	for _, remote := range remotes {
		refs := branchesWillBePush[remote]
		args := []string{"push", "--delete"}
		if force {
			args = append(args, "-f")
		}

		args = append(args, remote)

		for _, ref := range refs {
			args = append(args, ref.String())
		}

		err := git.command.Run(args...)
		if err != nil {
			return err
		}
	}

	return nil
}
