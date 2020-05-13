package statics

/*
BPAsset satisfies AssetEmbedder interface
This one is namespaced for root of the boilerplate directory
see boilerplate directory in bootstrap to see the content.
*/

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/common/embed"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/logger"
	_ "github.com/getcouragenow/core-bs/statiks/statik"
	"github.com/rakyll/statik/fs"
	"net/http"
)

type BPAsset struct {
	fsys http.FileSystem // the rakyll fs
	l    *logger.Logger
}

// NewBPAsset function to filter valid namespace
// for now this will be hardcoded, later down the line,
// it will be generated.
// func NewBPAsset(namespaces []string, namespaceArg string) (embed.AssetEmbedder, error) {
func NewBPAsset(l *logger.Logger, namespaceArg string) (embed.AssetEmbedder, error) {
	bfs, err := fs.NewWithNamespace(namespaceArg)
	if err != nil {
		return nil, err
	}
	return &BPAsset{
		fsys: bfs,
		l:    l,
	}, nil
}

func (r *BPAsset) GetFS() http.FileSystem { return r.fsys }
func (r *BPAsset) WriteAllFiles(outputPath string) error {
	return writeAllFiles(r.fsys, outputPath)
}
func (r *BPAsset) ReadSingleFile(name string) ([]byte, error) {
	return readSingleFile(r.fsys, name)
}
