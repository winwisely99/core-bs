package statics

import (
	"fmt"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
	"github.com/rakyll/statik/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func checkAndMakeDir(path string) error {
	exists, err := osutil.Exists(path)
	if err != nil {
		return err
	}
	if !exists {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}

func writeAllFiles(fsys http.FileSystem, outputPath string) error {
	if err := checkAndMakeDir(outputPath); err != nil {
		return err
	}
	if err := fs.Walk(fsys, "/", func(filePath string, fileInfo os.FileInfo, err error) error {
		newPath := path.Join(outputPath, filePath)
		if fileInfo.IsDir() {
			if err := checkAndMakeDir(newPath); err != nil {
				return fmt.Errorf("creating directory %q: %w", newPath, err)
			}
		} else {
			file, err := fsys.Open(filePath)
			if err != nil {
				return fmt.Errorf("opening %q in embedded filesystem: %w", filePath, err)
			}

			buf, err := ioutil.ReadAll(file)
			if err != nil {
				return fmt.Errorf("reading %q in embedded filesystem: %w", filePath, err)
			}

			if err := ioutil.WriteFile(newPath, buf, 0664); err != nil {
				return fmt.Errorf("writing %q to %q: %w", filePath, newPath, err)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	s, err := filepath.Abs(outputPath)
	if err != nil {
		return err
	}

	log.Printf("Successfully exported boilerplates to %s", s)
	return nil
}

func readSingleFile(fsys http.FileSystem, name string) ([]byte, error) {
	f, err := fsys.Open(fmt.Sprintf("/%s", name))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
