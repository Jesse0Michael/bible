package bible

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func NewFile(dir, name string) (*os.File, error) {
	p := newFilePath(dir, name, 0)
	return os.Create(p)
}

func newFilePath(dir, name string, index int) string {
	path := filepath.Join(dir, fmt.Sprintf("%s.yaml", strings.ToLower(name)))
	if index != 0 {
		path = filepath.Join(dir, fmt.Sprintf("%s%d.yaml", strings.ToLower(name), index+1))
	}

	if _, err := os.Open(path); err != nil {
		return path
	}

	return newFilePath(dir, name, index+1)
}

func FindReference(dir, name string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	var matches []os.FileInfo
	for _, f := range files {
		if strings.HasPrefix(f.Name(), strings.ToLower(name)) {
			matches = append(matches, f)
		}
	}
	count := len(matches)
	if count == 1 {
		return matches[0].Name(), nil
	}
	if count > 1 {
		return "", fmt.Errorf("ambiguous references found: %s/%s*", dir, name)
	}

	return "", fmt.Errorf("no reference found: %s/%s*", dir, name)

}
