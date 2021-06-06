package git

import (
	"reflect"
	"testing"

	"github.com/oakcask/git-stale/git/cli"
)

type fakeCommand struct {
	out  string
	err  error
	runs [][]string
}

func (c *fakeCommand) Call(args ...string) ([]byte, error) {
	return []byte(c.out), c.err
}

func (c *fakeCommand) Run(args ...string) error {
	c.runs = append(c.runs, args)
	return c.err
}

func TestGit_GetBranches(t *testing.T) {
	testCases := []struct {
		command     fakeCommand
		outBranches []Branch
	}{
		{
			fakeCommand{
				out: `a, gone
b, behind 1
c, behind 1, ahead 2
d
`,
				err: nil,
			},
			[]Branch{
				{
					Name: RefName("a"),
					Gone: true,
				},
				{
					Name: RefName("b"),
					Gone: false,
				},
				{
					Name: RefName("c"),
					Gone: false,
				},
				{
					Name: RefName("d"),
					Gone: false,
				},
			},
		},
	}

	for _, tc := range testCases {
		g := git{&tc.command, cli.Version{}}

		branches, e := g.GetBranches()
		if e != nil {
			t.Errorf("wants %+v but got unexpected error %v", tc.outBranches, e)
		} else {
			if !reflect.DeepEqual(branches, tc.outBranches) {
				t.Errorf("wants %+v but got %+v", tc.outBranches, branches)
			}
		}
	}
}

func TestGit_RemoveBranches(t *testing.T) {

	testCases := []struct {
		force        bool
		branches     []RefName
		expectedRuns [][]string
	}{
		{
			force:        false,
			branches:     []RefName{},
			expectedRuns: nil,
		},
		{
			force:    false,
			branches: []RefName{RefName("a"), RefName("b")},
			expectedRuns: [][]string{
				{"branch", "-d", "a", "b"},
			},
		},
	}

	for _, tc := range testCases {
		fakeCmd := fakeCommand{}
		g := git{&fakeCmd, cli.Version{}}

		err := g.RemoveBranches(tc.force, tc.branches...)

		if err != nil {
			t.Errorf("unexpectedly got error %v", err)
		} else {
			if len(fakeCmd.runs) != len(tc.expectedRuns) {
				t.Errorf("expected %v time(s) but got %v time(s)", len(tc.expectedRuns), len(fakeCmd.runs))
			} else if !reflect.DeepEqual(fakeCmd.runs, tc.expectedRuns) {
				t.Errorf("expected runs are %v but got %v", tc.expectedRuns, fakeCmd.runs)
			}
		}
	}
}
