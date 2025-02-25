package test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// CopyDirectory copies all files in the given src directory into the given dst directory.
// Both the source and the destination directory must exist.
func CopyDirectory(src, dst string) error {
	return filepath.Walk(src, func(srcPath string, fileInfo os.FileInfo, e error) error {
		dstPath := strings.Replace(srcPath, src, dst, 1)
		if fileInfo.IsDir() {
			err := os.Mkdir(dstPath, fileInfo.Mode())
			if err != nil {
				return fmt.Errorf("cannot create target directory: %w", err)
			}
			return nil
		}
		sourceFile, err := os.Open(srcPath)
		if err != nil {
			return fmt.Errorf("cannot open source file %q: %w", srcPath, err)
		}
		destFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY, fileInfo.Mode())
		if err != nil {
			return fmt.Errorf("cannot create target file %q: %w", srcPath, err)
		}
		_, err = io.Copy(destFile, sourceFile)
		if err != nil {
			return fmt.Errorf("cannot copy %q into %q: %w", srcPath, dstPath, err)
		}
		err = sourceFile.Close()
		if err != nil {
			return fmt.Errorf("cannot close source file %q: %w", srcPath, err)
		}
		err = destFile.Close()
		return err
	})
}

// createFile creates a file with the given filename in the given directory.
func createFile(t *testing.T, dir, filename string) {
	t.Helper()
	filePath := filepath.Join(dir, filename)
	err := os.MkdirAll(filepath.Dir(filePath), 0o744)
	assert.NoError(t, err)
	err = os.WriteFile(filePath, []byte(filename+" content"), 0o500)
	assert.NoError(t, err)
}

// CreateTempDir creates a new empty directory in the system's temp directory and provides the path to it.
func CreateTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "")
	assert.Nil(t, err, "cannot create TempDir")
	evalDir, err := filepath.EvalSymlinks(dir) // Evaluate symlinks as Mac temp dir is symlinked
	assert.Nil(t, err, "cannot evaluate symlinks of TempDir")
	return evalDir
}
