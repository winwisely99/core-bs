package gitutil

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/common/colorutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
	gg "github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
	"strings"
)

func GitClone(url string, dir string) error {
	log.Infof("Cloning: %s\n", colorutil.ColorGreen(url))
	r, err := gg.PlainClone(dir, false, &gg.CloneOptions{
		URL:               url,
		RecurseSubmodules: gg.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return err
	}
	rev, err := r.Head()
	if err != nil {
		return err
	}
	log.Infof("Successfully cloned: %s, sha: %s\n",
		url, colorutil.ColorGreen(rev.Hash()))
	return nil
}

func GitPull(dir string, tag string) error {
	log.Infof("Pulling from: %s\n", dir)
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
	remotes, err := r.Remotes()
	if err != nil {
		return "", err
	}
	var s string
	for _, remote := range remotes {
		if remote.Config().Name == "origin" {
			sses := strings.Split(remote.String(), " ")
			sse := strings.Split(sses[0], "\t")
			s = sse[1]
		}
	}
	return s, nil
}

func GitRemove(dir string) error {
	log.Infof("Removing directory....")
	_, err := osutil.RunUnixCmd(`rm`, `-rf`, dir)
	return err
}
