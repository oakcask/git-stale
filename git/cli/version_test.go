package cli

import (
	"errors"
	"fmt"
	"testing"
)

type fakeCommand struct {
	out string
	err error
}

func (c *fakeCommand) Call(args ...string) ([]byte, error) {
	return []byte(c.out), c.err
}

func (c *fakeCommand) Run(args ...string) error {
	fmt.Print(c.out)
	return c.err
}

func TestGetVersion(t *testing.T) {
	testCases := []struct {
		command  *fakeCommand
		expected Version
		err      error
	}{
		{
			&fakeCommand{
				out: "git version 1.23.456\n",
				err: nil,
			},
			Version{Major: 1, Minor: 23, Patch: 456},
			nil,
		},
		{
			&fakeCommand{
				out: "git version 1.23.456a\n",
				err: nil,
			},
			Version{},
			errors.New("cannot parse: git version 1.23.456a\n"),
		},
		{
			&fakeCommand{
				out: "git version 1.23.\n",
				err: nil,
			},
			Version{},
			errors.New("cannot parse: git version 1.23.\n"),
		},
		{
			&fakeCommand{
				out: "",
				err: errors.New("exec failed"),
			},
			Version{},
			errors.New("exec failed"),
		},
	}

	for _, tc := range testCases {
		actual, err := GetVersion(tc.command)
		if tc.err != nil {
			if err == nil {
				t.Errorf("want error but got nil")
			} else if tc.err.Error() != err.Error() {
				t.Errorf("want error `%v` but got `%v`", tc.err, err)
			}
		}
		if tc.expected != actual {
			t.Errorf("want %+v but got %+v", tc.expected, actual)
		}
	}
}
