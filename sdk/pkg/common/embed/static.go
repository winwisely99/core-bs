package embed

import (
	"net/http"
)

// AssetEmbedder is the common interface for all static files
// included in the bs tool.
// all of the embedded files has to be namespaced for cleaner approach.
type AssetEmbedder interface {
	GetFS() http.FileSystem
	WriteAllFiles(outputPath string) error // writes populated map[string][]byte to user's filesystem
	ReadSingleFile(filename string) ([]byte, error) // read single file from the namespace (without populating)
}

