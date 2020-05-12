package osutil

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	CheckMark = "\u2713"
	CrossMark = "\u274c"
)

// BinExists checks binary name in the user's PATH
func BinExists(binname string) bool {
	p, err := PrintBinaryPath(binname)
	if err != nil {
		return false
	}
	log.Println(p)
	return true
}

// PrintBinaryPath prints binary location in the user's PATH
func PrintBinaryPath(binname string) (string, error) {
	return exec.LookPath(binname)
}

// blanket implementation for Unices / *nix-like OSes
func RunUnixCmd(displayOutput bool, cmdName string, flags ...string) (*string, error) {
	cmd := exec.Command(cmdName, flags...)
	var stdoutBuf, stderrBuf bytes.Buffer

	if displayOutput {
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	} else {
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf
	}
	// we will ignore error from cmd run
	// and only assign error when it is from stderr
	cmd.Run()
	buf := strings.TrimSpace(string(stdoutBuf.Bytes()))
	err := strings.TrimSpace(string(stderrBuf.Bytes()))
	if err != "" {
		return nil, errors.New(err)
	}
	return &buf, nil
}

// WalkPath gets all the directory names in the specified directory
// returns all the directory paths that contains a Makefile.
//func WalkPath(l *logger.Logger, dir string) (dirs map[string]string, err error) {
//	dirs = map[string]string{}
//	l.Infof("Scanning : %s\n", dir)
//	var allDirs []string
//	if err := godirwalk.Walk(dir, &godirwalk.Options{
//		Callback: func(pathname string, de *godirwalk.Dirent) error {
//			if de.IsDir() {
//				allDirs = append(allDirs, filepath.FromSlash(pathname))
//			}
//			return nil
//		},
//		ErrorCallback: func(pathname string, err error) godirwalk.ErrorAction {
//			log.Errorf("ERROR: %v\n", err)
//			return godirwalk.SkipNode
//		},
//	}); err != nil {
//		log.Errorf("%v\n", err)
//		return dirs, err
//	}
//	var scanner *godirwalk.Scanner
//	for _, dir := range allDirs {
//		scanner, err = godirwalk.NewScanner(dir)
//		if err != nil {
//			log.Errorf("cannot lazily scan dir %s: %v\n", dir, err)
//			return nil, err
//		}
//		for scanner.Scan() {
//			dirent, err := scanner.Dirent()
//			if err != nil {
//				log.Warnf("cannot get directory content: %v", err)
//				continue
//			}
//			if dirent.Name() == "Makefile" || dirent.Name() == "statik.go" {
//				splitted := strings.Split(dir, "/")
//				dirs[splitted[len(splitted)-1]] = dir
//			}
//		}
//	}
//	return dirs, err
//}

func ExtractTarGz(gzipStream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(uncompressedStream)
	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			outFile.Close()

		default:
			return errors.New(fmt.Sprintf(
				"ExtractTarGz: uknown type: %s in %s",
				header.Typeflag,
				header.Name))
		}
	}
	return nil
}
