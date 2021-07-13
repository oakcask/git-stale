package git

import (
	"fmt"
	"time"

	"github.com/oakcask/git-stale/git/cli"
)

type Git interface {
	GetBranches() ([]Branch, error)
	RemoveBranches(force bool, refnames ...RefName) error
}

type Branch struct {
	Name       RefName
	Gone       bool
	CommitDate time.Time
}

type RefName string

func (ref *RefName) String() string {
	return string(*ref)
}

type git struct {
	command cli.Command
	version cli.Version
}

func NewGit(command cli.Command) (Git, error) {
	version, e := cli.GetVersion(command)
	if e != nil {
		return nil, fmt.Errorf("cannot recognize git version")
	}

	return &git{command, version}, nil
}
