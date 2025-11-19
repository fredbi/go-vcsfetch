package git

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/go-git/go-billy/v5"
	// "github.com/go-git/go-billy/v5/helper/iofs"
)

var _ fs.FS = &fsWrapper{}

type fsWrapper struct { // this is needed until go-billy/v6 and go-git/v6 are released
	billy.Filesystem
}

func (f *fsWrapper) Open(path string) (fs.File, error) {
	info, _ := f.Filesystem.Stat(path)
	if info.IsDir() {
		dir, err := f.Filesystem.ReadDir(path)
		if err != nil {
			return nil, err
		}
		return &dirWrapper{name: path, entries: dir, fs: f.Filesystem}, nil
	}

	file, err := f.Filesystem.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %q: %w", path, err)
	}

	return &fileWrapper{File: file, fs: f.Filesystem}, nil
}

var _ fs.File = &fileWrapper{}

type fileWrapper struct {
	billy.File

	fs billy.Filesystem
}

func (f *fileWrapper) Stat() (fs.FileInfo, error) {
	return f.fs.Stat(f.File.Name())
}

var _ fs.DirEntry = &dirEntryWrapper{}

type dirEntryWrapper struct {
	fs.FileInfo

	fs billy.Filesystem
}

func (d *dirEntryWrapper) Info() (fs.FileInfo, error) {
	return d.FileInfo, nil
}

func (d *dirEntryWrapper) Type() fs.FileMode {
	return d.FileInfo.Mode()
}

var _ fs.ReadDirFile = &dirWrapper{}

type dirWrapper struct {
	name    string
	entries []os.FileInfo
	offset  int
	fs      billy.Filesystem
}

func (f *dirWrapper) Stat() (fs.FileInfo, error) {
	return f.fs.Stat(f.name)
}

func (f *dirWrapper) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (f *dirWrapper) Close() error {
	f.offset = 0

	return nil
}

// ReadDir(n int) ([]DirEntry, error)
func (f *dirWrapper) ReadDir(n int) ([]fs.DirEntry, error) {
	l := len(f.entries)

	if l == 0 || f.offset >= l {
		return []fs.DirEntry{}, io.EOF
	}

	var m int
	if n <= 0 {
		m = l
	} else {
		m = min(n, l-f.offset)
	}

	entries := make([]fs.DirEntry, 0, m)
	for _, entry := range f.entries[f.offset : f.offset+m] {
		entries = append(entries, &dirEntryWrapper{
			FileInfo: entry,
		})
	}
	f.offset += m

	return entries, nil
}
