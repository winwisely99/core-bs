package oses

import (
	"errors"
	"fmt"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/termutil"
	m "github.com/pbnjay/memory"
	"os"
	"runtime"
)

// Blanket os info getter for all OSes (including windows)
type OSInfoGetter interface {
	GetOsName() string
	GetKernel() string
	GetPlatform() string
	GetHostName() string
	GetMemory() float64
	GetCores() int
	String() string
	ToContent() termutil.Contents
}

func getOsInfoGetter() (OSInfoGetter, error) {
	switch runtime.GOOS {
	case "windows":
		return getWindowsOsInfo()
	case "darwin":
		return getDarwinOsInfo()
	case "linux":
		return getLinuxOsInfo()
	default:
		return nil, errors.New("unknown / unsupported OS")
	}
}

// Darwin
type DarwinOSInfo struct {
	osName   string
	kernel   string
	platform string
	hostName string
	memory   float64
	cores    int
}

func getDarwinOsInfo() (*DarwinOSInfo, error) {
	var osName, kernel, platform, hostname *string
	var err error
	if osName, err = getUnixOSName(); err != nil {
		return nil, err
	}
	if kernel, err = getUnixKernel(); err != nil {
		return nil, err
	}
	if platform, err = getUnixPlatform(); err != nil {
		return nil, err
	}
	if hostname, err = getUnixHostname(); err != nil {
		return nil, err
	}
	mem := getMemory()
	core := getCPUCore()
	return &DarwinOSInfo{
		osName:   *osName,
		kernel:   *kernel,
		platform: *platform,
		hostName: *hostname,
		memory:   mem,
		cores:    core,
	}, nil
}
func (d *DarwinOSInfo) GetOsName() string   { return d.osName }
func (d *DarwinOSInfo) GetKernel() string   { return d.kernel }
func (d *DarwinOSInfo) GetPlatform() string { return d.platform }
func (d *DarwinOSInfo) GetHostName() string { return d.hostName }
func (d *DarwinOSInfo) GetMemory() float64  { return d.memory }
func (d *DarwinOSInfo) GetCores() int       { return d.cores }
func (d *DarwinOSInfo) String() string {
	return fmt.Sprintf("OS: %s, Kernel: %s, Platform: %s, Hostname: %s, Cores: %d, Memory: %d",
		d.osName, d.kernel, d.platform, d.hostName, d.cores, d.memory)
}
func (d *DarwinOSInfo) ToContent() termutil.Contents { return toContent(d) }

// Linux
type LinuxOSInfo struct {
	osName   string
	kernel   string
	platform string
	hostName string
	memory   float64
	cores    int
}

func getLinuxOsInfo() (*LinuxOSInfo, error) {
	var osName, kernel, platform, hostname *string
	var err error
	if osName, err = getUnixOSName(); err != nil {
		return nil, err
	}
	if kernel, err = getUnixKernel(); err != nil {
		return nil, err
	}
	if platform, err = getUnixPlatform(); err != nil {
		return nil, err
	}
	if hostname, err = getUnixHostname(); err != nil {
		return nil, err
	}
	mem := getMemory()
	return &LinuxOSInfo{
		osName:   *osName,
		kernel:   *kernel,
		platform: *platform,
		hostName: *hostname,
		memory:   mem,
		cores:    getCPUCore(),
	}, nil
}
func (l *LinuxOSInfo) GetOsName() string   { return l.osName }
func (l *LinuxOSInfo) GetKernel() string   { return l.kernel }
func (l *LinuxOSInfo) GetPlatform() string { return l.platform }
func (l *LinuxOSInfo) GetHostName() string { return l.hostName }
func (l *LinuxOSInfo) GetMemory() float64  { return l.memory }
func (l *LinuxOSInfo) GetCores() int       { return l.cores }
func (l *LinuxOSInfo) String() string {
	return fmt.Sprintf("OS: %s, Kernel: %s, Platform: %s, Hostname: %s, Cores: %d, Memory: %d",
		l.osName, l.kernel, l.platform, l.hostName, l.cores, l.memory)
}
func (l *LinuxOSInfo) ToContent() termutil.Contents { return toContent(l) }

// Windows
type WindowsOSInfo struct {
	osName   string
	kernel   string
	platform string
	hostName string
	memory   float64
	cores    int
}

func getWindowsOsInfo() (*WindowsOSInfo, error) {
	var osName, platform, hostname string
	var err error
	osName = runtime.GOOS
	hostname, err = os.Hostname()
	if err != nil {
		return nil, err
	}
	platform = runtime.GOARCH
	mem := getMemory()
	return &WindowsOSInfo{
		osName:   osName,
		kernel:   "Windows",
		platform: platform,
		hostName: hostname,
		memory:   mem,
		cores:    getCPUCore(),
	}, nil
}
func (w *WindowsOSInfo) GetOsName() string   { return w.osName }
func (w *WindowsOSInfo) GetKernel() string   { return w.kernel }
func (w *WindowsOSInfo) GetPlatform() string { return w.platform }
func (w *WindowsOSInfo) GetHostName() string { return w.hostName }
func (w *WindowsOSInfo) GetMemory() float64  { return w.memory }
func (w *WindowsOSInfo) GetCores() int       { return w.cores }
func (w *WindowsOSInfo) String() string {
	return fmt.Sprintf("OS: %s, Kernel: %s, Platform: %s, Hostname: %s, Cores: %d, Memory: %d",
		w.osName, w.kernel, w.platform, w.hostName, w.cores, w.memory)
}
func (w *WindowsOSInfo) ToContent() termutil.Contents { return toContent(w) }

// blanket implementation for Unices / *nix-like OSes 
func RunCmd(cmdName string, flags ...string) (*string, error) {
	return osutil.RunCmd(false, cmdName, flags...)
}

func getUnixUname(flag string) (*string, error) {
	return RunCmd("uname", flag)
}

func getUnixPlatform() (*string, error) {
	return getUnixUname("-m")
}

func getUnixHostname() (*string, error) {
	return getUnixUname("-n")
}

func getUnixKernel() (*string, error) {
	return getUnixUname("-r")
}

func getUnixOSName() (*string, error) {
	return getUnixUname("-s")
}

func getCPUCore() int {
	return runtime.NumCPU()
}

func toContent(o OSInfoGetter) termutil.Contents {
	ms := map[string][]string{}
	ms["OS"] = []string{o.GetOsName()}
	ms["Kernel / Version"] = []string{o.GetKernel()}
	ms["Platform"] = []string{o.GetPlatform()}
	ms["Cores"] = []string{fmt.Sprintf("%d", o.GetCores())}
	ms["Memory"] = []string{fmt.Sprintf("%.2f MiB", o.GetMemory())}
	ms["Hostname"] = []string{o.GetHostName()}
	return ms
}

func getMemory() float64 {
	mem := m.TotalMemory()
	return float64(mem) / 1000000
}
