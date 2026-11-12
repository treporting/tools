package review

import (
	"os/exec"
	"strings"
)

// ExecGit asks the git of the machine about the repo in Root.
type ExecGit struct {
	Root string
}

func (g ExecGit) run(args ...string) (string, error) {
	out, err := exec.Command("git", append([]string{"-C", g.Root}, args...)...).Output()
	return string(out), err
}

// CommitExists says whether rev names a commit.
func (g ExecGit) CommitExists(rev string) bool {
	if rev == "" {
		return false
	}
	_, err := g.run("cat-file", "-e", rev+"^{commit}")
	return err == nil
}

// Show answers the content of path in the commit rev.
func (g ExecGit) Show(rev, path string) ([]byte, error) {
	out, err := g.run("show", rev+":"+path)
	return []byte(out), err
}

// FileCommits answers the commits up to rev that changed path, the newest first.
func (g ExecGit) FileCommits(rev, path string) ([]string, error) {
	out, err := g.run("log", "--format=%H", rev, "--", path)
	if err != nil {
		return nil, err
	}
	return strings.Fields(out), nil
}

// Parent answers the first parent of a commit, and false for the first commit.
func (g ExecGit) Parent(commit string) (string, bool) {
	out, err := g.run("rev-parse", "--verify", "--quiet", commit+"^1")
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(out), true
}

// IsMerge says whether a commit has more than one parent.
func (g ExecGit) IsMerge(commit string) bool {
	_, err := g.run("rev-parse", "--verify", "--quiet", commit+"^2")
	return err == nil
}

// Short answers the short hash of rev.
func (g ExecGit) Short(rev string) (string, error) {
	out, err := g.run("rev-parse", "--short", rev)
	return strings.TrimSpace(out), err
}

// Authors answers the author of a commit and the names of its co-authors, without their address.
func (g ExecGit) Authors(commit string) (string, []string, error) {
	out, err := g.run("show", "-s", "--format=%an%x00%(trailers:key=Co-Authored-By,valueonly,separator=%x1f)", commit)
	if err != nil {
		return "", nil, err
	}
	author, trailers, _ := strings.Cut(strings.TrimRight(out, "\n"), "\x00")
	var coAuthors []string
	for _, t := range strings.Split(trailers, "\x1f") {
		name, _, _ := strings.Cut(strings.TrimSpace(t), " <")
		if name != "" {
			coAuthors = append(coAuthors, name)
		}
	}
	return author, coAuthors, nil
}

// Changed answers the files changed between base and HEAD.
func (g ExecGit) Changed(base string) ([]string, error) {
	out, err := g.run("diff", "--name-only", base, "HEAD")
	if err != nil {
		return nil, err
	}
	return strings.Fields(out), nil
}
