package gitutil

import (
	"fmt"
	"strings"

	"github.com/getcouragenow/core-bs/sdk/pkg/common/colorutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/logger"
	gg "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GitInfo contains sha commit,
// tag, branch, and remote of type origin
// for current repo.
type GitInfo struct {
	ShaCommit string
	LatestTag string
	Branch    string
	Origin    string
}

func IsGitDir(dir string) bool {
	_, err := gg.PlainOpen(dir)
	if err != nil {
		return false
	}
	return true
}

func GitRemoteInfo(dir string) (*string, error) {
	r, err := gg.PlainOpen(dir)
	if err != nil {
		return nil, err
	}
	rem, err := r.Remote("origin")
	if err != nil {
		return nil, err
	}
	curOrigin := rem.Config().URLs[0]
	if &curOrigin == nil {
		return nil, err
	}
	return &curOrigin, nil
}

// GitGetInfo function is a utility to get current repo's
// tag, sha commit, and branch.
func GitGetInfo(dir string) (*GitInfo, error) {
	r, err := gg.PlainOpen(dir)
	if err != nil {
		return nil, err
	}
	rev, err := r.Head()
	if err != nil {
		return nil, err
	}
	commitSha := rev.Hash().String()
	curBranch := rev.Strings()[0]
	tags, err := r.Tags()
	if err != nil {
		return nil, err
	}
	lastTag, err := tags.Next()
	if err != nil {
		return nil, err
	}
	curTag := lastTag.Name().Short()
	rem, err := r.Remote("origin")
	if err != nil {
		return nil, err
	}
	curOrigin := rem.Config().URLs
	return &GitInfo{
		ShaCommit: commitSha,
		Branch:    curBranch,
		Origin:    fmt.Sprintf("%s", curOrigin),
		LatestTag: curTag,
	}, nil
}

func GitClone(l *logger.Logger, url, dir string) error {
	l.Infof("Cloning: %s\n", colorutil.ColorGreen(url))
	r, err := gg.PlainClone(dir, false, &gg.CloneOptions{
		URL:               url,
		RecurseSubmodules: gg.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		if err == gg.ErrRepositoryAlreadyExists {
			return GitPull(l, dir, "refs/heads/master")
		} else {
			return err
		}
	}
	rev, err := r.Head()
	if err != nil {
		return err
	}
	l.Infof("Successfully cloned: %s, sha: %s\n",
		url, colorutil.ColorGreen(rev.Hash()))
	return nil
}

func GitCheckout(l *logger.Logger, dir, branch string) error {
	l.Infof("Checkout: %s\n", colorutil.ColorGreen(branch))
	r, err := gg.PlainOpen(dir)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	cBranch := strings.Join([]string{"refs/heads/", branch}, "")
	return w.Checkout(&gg.CheckoutOptions{
		Branch: plumbing.ReferenceName(cBranch),
	})
}

func GitPull(l *logger.Logger, dir, tag string) error {
	l.Infof("Pulling from: %s\n", dir)
	r, err := gg.PlainOpen(dir)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	if err = w.Pull(&gg.PullOptions{}); err != nil {
		if err == gg.NoErrAlreadyUpToDate {
			return nil
		}
		return err
	}
	return nil
}
