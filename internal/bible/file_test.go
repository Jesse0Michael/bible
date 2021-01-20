package bible

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFile(t *testing.T) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	got, err := NewFile(dir, "new")

	assert.NotNil(t, got)
	assert.NoError(t, err)
}

func Test_newFilePath(t *testing.T) {
	dir, _ := ioutil.TempDir("", "")
	_, _ = os.Create(filepath.Join(dir, "exists.yaml"))
	_, _ = os.Create(filepath.Join(dir, "many.yaml"))
	_, _ = os.Create(filepath.Join(dir, "many2.yaml"))
	defer os.RemoveAll(dir)
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "return new filepath",
			arg:  "new",
			want: filepath.Join(dir, "new.yaml"),
		},
		{
			name: "return second index",
			arg:  "exists",
			want: filepath.Join(dir, "exists2.yaml"),
		},
		{
			name: "return third index",
			arg:  "many",
			want: filepath.Join(dir, "many3.yaml"),
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			got := newFilePath(dir, tt.arg, 0)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFindReference(t *testing.T) {
	dir, _ := ioutil.TempDir("", "")
	_, _ = os.Create(filepath.Join(dir, "exists.yaml"))
	_, _ = os.Create(filepath.Join(dir, "many.yaml"))
	_, _ = os.Create(filepath.Join(dir, "many2.yaml"))
	defer os.RemoveAll(dir)
	tests := []struct {
		name    string
		dir     string
		arg     string
		want    string
		wantErr string
	}{
		{
			name:    "invalid directory",
			dir:     "not_here",
			wantErr: "open not_here: no such file or directory",
		},
		{
			name: "found reference",
			dir:  dir,
			arg:  "Exists",
			want: "exists.yaml",
		},
		{
			name:    "no reference",
			dir:     dir,
			arg:     "other",
			wantErr: fmt.Sprintf("no reference found: %s/other*", dir),
		},
		{
			name:    "too many references",
			dir:     dir,
			arg:     "many",
			wantErr: fmt.Sprintf("ambiguous references found: %s/many*", dir),
		},
		{
			name: "specify reference",
			dir:  dir,
			arg:  "many.yaml",
			want: "many.yaml",
		},
		{
			name: "specify reference",
			dir:  dir,
			arg:  "many2.yaml",
			want: "many2.yaml",
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindReference(tt.dir, tt.arg)

			assert.Equal(t, tt.want, got)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
