package util
import (
	"path/filepath"
	"os"
	"io"
	"strings"
)

func CopyUnderFolder(srcDir, dstDir string) {
	filepath.Walk(srcDir, func(path string, f os.FileInfo, err error) error {
		if (f.IsDir()) {
			os.MkdirAll(dstDir + strings.TrimPrefix(path, srcDir), f.Mode())
		} else {
			copyFileContents(path, dstDir + strings.TrimPrefix(path, srcDir))
		}
		return nil
	})
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}