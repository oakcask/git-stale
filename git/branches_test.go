package git

import (
	"reflect"
	"testing"
	"time"

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

func stdISO8601(s string) time.Time {
	t, _ := time.Parse("20060102T150405-0700", s)
	return t
}

func TestGit_GetBranches(t *testing.T) {
	testCases := []struct {
		command     fakeCommand
		outBranches []Branch
	}{
		{
			fakeCommand{
				out: `a, Thu Feb 4 20:38:23 2021 +0900, origin, gone
b, Tue Mar 30 22:22:02 2021 +0900, origin, behind 1
c, Fri Apr 23 17:36:01 2021 +0900, origin, behind 1, ahead 2
d, Thu Jun 10 08:12:17 2021 +0900, 
`,
				err: nil,
			},
			[]Branch{
				{
					Name:       RefName("a"),
					Gone:       true,
					RemoteName: "origin",
					CommitDate: stdISO8601("20210204T203823+0900"),
				},
				{
					Name:       RefName("b"),
					Gone:       false,
					RemoteName: "origin",
					CommitDate: stdISO8601("20210330T222202+0900"),
				},
				{
					Name:       RefName("c"),
					Gone:       false,
					RemoteName: "origin",
					CommitDate: stdISO8601("20210423T173601+0900"),
				},
				{
					Name:       RefName("d"),
					Gone:       false,
					RemoteName: "",
					CommitDate: stdISO8601("20210610T081217+0900"),
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
		branches     []Branch
		expectedRuns [][]string
	}{
		{
			force:        false,
			branches:     []Branch{},
			expectedRuns: nil,
		},
		{
			force:    false,
			branches: []Branch{{Name: RefName("a")}, {Name: RefName("b")}},
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

func TestGit_RemoveRemoteBranches(t *testing.T) {

	testCases := []struct {
		force        bool
		branches     []Branch
		expectedRuns [][]string
	}{
		{
			force:        false,
			branches:     []Branch{},
			expectedRuns: nil,
		},
		{
			force: false,
			branches: []Branch{
				{RemoteName: "origin", Name: RefName("a")},
				{RemoteName: "origin", Name: RefName("b")},
				{RemoteName: "upstream", Name: RefName("c")},
				{RemoteName: "", Name: RefName("d")},
			},
			expectedRuns: [][]string{
				{"push", "--delete", "origin", "a", "b"},
				{"push", "--delete", "upstream", "c"},
			},
		},
		{
			force: true,
			branches: []Branch{
				{RemoteName: "origin", Name: RefName("a")},
				{RemoteName: "upstream", Name: RefName("b")},
				{RemoteName: "", Name: RefName("c")},
			},
			expectedRuns: [][]string{
				{"push", "--delete", "-f", "origin", "a"},
				{"push", "--delete", "-f", "upstream", "b"},
			},
		},
	}

	for _, tc := range testCases {
		fakeCmd := fakeCommand{}
		g := git{&fakeCmd, cli.Version{}}

		err := g.RemoveRemoteBranches(tc.force, tc.branches...)

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
