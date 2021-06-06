package cli

import (
	"bytes"
	"fmt"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

type Command interface {
	Call(args ...string) ([]byte, error)
	Run(args ...string) error
}

func GetVersion(command Command) (Version, error) {
	out, e := command.Call("--version")
	if e != nil {
		return Version{}, e
	}

	buffer := bytes.NewBuffer(out)
	var major, minor, patch int
	n, e := fmt.Fscanf(buffer, "git version %d.%d.%d\n", &major, &minor, &patch)
	if e != nil || n != 3 {
		return Version{}, fmt.Errorf("cannot parse: %v", string(out))
	}

	return Version{major, minor, patch}, nil
}
