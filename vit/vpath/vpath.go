package vpath

import (
	"fmt"
	"io/fs"
	"os"
	"path"
)

type Path interface {
	fs.ReadDirFS
	OpenFile() (fs.File, error)
	Path() string
	Dir() Path
}

type fsPath struct {
	fs   fs.ReadDirFS
	path string
}

// FS returns a new path that is contained in the given filesystem.
func FS(fs fs.ReadDirFS, p string) Path {
	return fsPath{fs: fs, path: p}
}

func (p fsPath) OpenFile() (fs.File, error) {
	return p.fs.Open(p.path)
}

func (p fsPath) Path() string {
	return p.path
}

func (p fsPath) Open(name string) (fs.File, error) {
	return p.fs.Open(path.Join(p.path, name))
}

func (p fsPath) ReadDir(name string) ([]fs.DirEntry, error) {
	return p.fs.ReadDir(path.Join(p.path, name))
}

func (p fsPath) Dir() Path {
	return FS(p.fs, path.Dir(p.path))
}

func (p fsPath) String() string {
	return fmt.Sprintf("FS://%s", p.path)
}

type localPath string

// Local returns a new path that is contained in the local filesystem.
func Local(p string) Path {
	return localPath(p)
}

func (p localPath) OpenFile() (fs.File, error) {
	return os.Open(string(p))
}

func (p localPath) Path() string {
	return string(p)
}

func (p localPath) Open(name string) (fs.File, error) {
	return os.Open(path.Join(string(p), name))
}

func (p localPath) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(path.Join(string(p), name))
}

func (p localPath) Dir() Path {
	return Local(path.Dir(string(p)))
}

func (p localPath) String() string {
	return string(p)
}

type virtualPath string

// Virtual returns a new path that does not actually exist.
// It is used in generated code to refer to source files that are no longer present during execution.
// Most operations on a virtual path will fail with fs.ErrNotExist.
func Virtual(p string) Path {
	return virtualPath(p)
}

func (p virtualPath) OpenFile() (fs.File, error) {
	return nil, fs.ErrNotExist
}

func (p virtualPath) Path() string {
	return string(p)
}

func (p virtualPath) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

func (p virtualPath) ReadDir(name string) ([]fs.DirEntry, error) {
	return nil, fs.ErrNotExist
}

func (p virtualPath) Dir() Path {
	return Virtual(path.Dir(string(p)))
}

func (p virtualPath) String() string {
	return fmt.Sprintf("VRT://%s", string(p))
}
