package gitutil

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/common/colorutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/logger"
	gg "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"strings"
)

func GitClone(l *logger.Logger, url, dir string) error {
	l.Infof("Cloning: %s\n", colorutil.ColorGreen(url))
	r, err := gg.PlainClone(dir, false, &gg.CloneOptions{
		URL:               url,
		RecurseSubmodules: gg.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		if err == gg.ErrRepositoryAlreadyExists {
			return GitPull(l, dir, "master")
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

func GitCheckout(l *logger.Logger, dir, tag string) error {
	l.Infof("Checkout: %s\n", colorutil.ColorGreen(tag))
	r, err := gg.PlainOpen(dir)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	return w.Checkout(&gg.CheckoutOptions{
		Branch: plumbing.ReferenceName(tag),
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
	return w.Pull(&gg.PullOptions{
		RemoteName: tag,
	})
}

func GitRemoteInfo(dir string) (string, error) {
	r, err := gg.PlainOpen(dir)
	if err != nil {
		return "", err
	}
	remote, err := r.Remote("origin")
	if err != nil {
		return "", err
	}
	var s string
	sses := strings.Split(remote.String(), " ")
	sse := strings.Split(sses[0], "\t")
	s = sse[1]
	return s, nil
}
