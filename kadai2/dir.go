package conv

import (
	"fmt"
	"os"
	"path/filepath"
)

// TraverseDir does DFS and return the file list.
// Please pass an empty string slice to `fileList`.
func TraverseDir(fileList []string, recursive bool, dir string, opts *Options) error {
	es, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("conv.DoDir %v (%w)", err, ErrFileAccess)
	}
	for _, e := range es {
		filePath := filepath.Join(dir, e.Name())
		if recursive && e.IsDir() {
			// Ignores errors during traversing as e.IsDir must be true here.
			TraverseDir(fileList, true, filePath, opts)
		} else {
			fileList = append(fileList, filePath)
		}
	}
	return nil
}
