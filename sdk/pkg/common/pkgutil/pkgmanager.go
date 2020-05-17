package pkgutil

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
)

type PkgManager interface {
	String() string
	Setup() error
	CleanCache() error
	Search(pkg string) (*string, error)
	Install(pkgs ...string) error
	Uninstall(pkgs ...string) error
	Upgrade() error
	Update(pkg string) error
}

// MacOS
type brew struct{}

// Windows
type scoop struct{}

// dnf (Fedora / RHEL / CentOS)
type dnf struct{}

// apt (Ubuntu / Debian)
type apt struct{}

func (b *brew) String() string {
	return "brew"
}

func (b *brew) Setup() error {
	_, err := osutil.RunCmd(
		true, `xcode-select`, `--install`,
	)
	if err != nil{ return err}
	_, err = osutil.RunCmd(
		true,
		"/bin/bash",
		"-c",
		"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)",
	)
	return err
}

func (b *brew) CleanCache() error {
	_, err := osutil.RunCmd(true, `brew`, `cleanup`)
	return err
}

func (b *brew) Search(pkg string) (*string, error) {
	return osutil.RunCmd(true, `brew`, `search`, pkg)
}

func (b *brew) Install(pkgs ...string) error {
	flags := append([]string{"install"}, pkgs...)
	_, err := osutil.RunCmd(true, `brew`, flags...)
	return err
}

func (b *brew) Uninstall(pkgs ...string) error {
	flags := append([]string{"uninstall"}, pkgs...)
	_, err := osutil.RunCmd(true, `brew`, flags...)
	return err
}

func (b *brew) Upgrade() error {
	_, err := osutil.RunCmd(true, `brew`, `upgrade`)
	return err
}

func (b *brew) Update(pkg string) error {
	_, err := osutil.RunCmd(true, `brew`, `update`, pkg)
	return err
}

func (s *scoop) String() string {
	return "scoop"
}

func (s *scoop) Setup() error {
	if _, err := osutil.RunCmd(
		true,
		"Set-ExecutionPolicy RemoteSigned -scope CurrentUser",
	); err != nil {
		return err
	}
	_, err := osutil.RunCmd(
		true,
		"iwr", `-useb, get.scoop.sh | iex`,
	)
	if err != nil { return err }

	_, err = osutil.RunCmd(
		true,
		"scoop", `bucket add extras`,
	)
	_, err = osutil.RunCmd(
		true,
		"scoop", `bucket add java`,
	)
	if err = s.Install("aria2"); err != nil { return err }
	_, err = osutil.RunCmd(
		true,
		"scoop", `config aria2-enabled true`,
	)
	return err
}

func (s *scoop) CleanCache() error {
	_, err := osutil.RunCmd(true, `scoop`, `cache`, `rm *`)
	return err
}

func (s *scoop) Search(pkg string) (*string, error) {
	return osutil.RunCmd(true, `scoop`, `search`, pkg)
}

func (s *scoop) Install(pkgs ...string) error {
	flags := append([]string{"install"}, pkgs...)
	_, err := osutil.RunCmd(true, `scoop`, flags...)
	return err
}

func (s *scoop) Uninstall(pkgs ...string) error {
	flags := append([]string{"uninstall"}, pkgs...)
	_, err := osutil.RunCmd(true, `scoop`, flags...)
	return err
}

func (s *scoop) Upgrade() error {
	_, err := osutil.RunCmd(true, `scoop`, `update`, `*`)
	return err
}

func (s *scoop) Update(pkg string) error {
	_, err := osutil.RunCmd(true, `brew`, `update`, pkg)
	return err
}

func (a *apt) Setup() error {
	_, err := osutil.SudoRunCmd(
		true,
		`apt`, `update`, `-y`,
	)
	if err != nil { return err }
	_, err = osutil.SudoRunCmd(true, `apt`, `dist-upgrade`, `-y`)
	return err
}

func (a *apt) String() string {
	return "apt"
}

func (a *apt) CleanCache() error {
	_, err := osutil.SudoRunCmd(
		true, `apt-get`, `clean` ,`-y`)
	return err
}

func (a *apt) Search(pkg string) (*string, error) {
	return osutil.RunCmd(true, `apt`, `search`, pkg)
}

func (a *apt) Install(pkgs ...string) error {
	flags := append([]string{"install", "-y"}, pkgs...)
	_, err := osutil.SudoRunCmd(true, `apt`, flags...)
	return err
}

func (a *apt) Uninstall(pkgs ...string) error {
	flags := append([]string{"purge", "-y"}, pkgs...)
	_, err := osutil.SudoRunCmd(true, `apt`, flags...)
	return err
}

func (a *apt) Upgrade() error {
	_, err := osutil.SudoRunCmd(true, `apt`, `upgrade`, `-y`)
	return err
}

func (a *apt) Update(pkg string) error {
	_, err := osutil.SudoRunCmd(true, `apt`,
		`upgrade`, `-y`)
	return err
}

func (d *dnf) Setup() error {
	return d.Upgrade()
}

func (d *dnf) CleanCache() error {
	_, err := osutil.SudoRunCmd(
		true, `dnf`, `clean` ,`-y`)
	return err
}

func (d *dnf) Search(pkg string) (*string, error) {
	return osutil.RunCmd(true, `dnf`, `search`, pkg)
}

func (d *dnf) Install(pkgs ...string) error {
	flags := append([]string{"install", "-y"}, pkgs...)
	_, err := osutil.SudoRunCmd(true, `dnf`, flags...)
	return err
}

func (d *dnf) Uninstall(pkgs ...string) error {
	flags := append([]string{"purge", "-y"}, pkgs...)
	_, err := osutil.SudoRunCmd(true, `dnf`, flags...)
	return err
}

func (d *dnf) Upgrade() error {
	_, err := osutil.SudoRunCmd(
		true,
		`dnf`, `--refresh`, `upgrade`, `--best`, `--allowerasing`, `-y`,
	)
	return err
}

func (d *dnf) Update(pkg string) error {
	return d.Upgrade()
}

func (d *dnf) String() string {
	return "dnf"
}
