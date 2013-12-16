package embedfs

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func DirAlloc(name string) *_dir {
	return &_dir{
		name:    name,
		modTime: time.Now(),
		dirs:    make(map[string]*_dir),
		files:   make(map[string]*EmbedFile),
	}
}

// http.FileSystem
// type FileSystem interface {
// 	Open(name string) (File, error)
// }
//
// http.File
// type File interface {
// 	Close() error
// 	Stat() (os.FileInfo, error)
// 	Readdir(count int) ([]os.FileInfo, error)
// 	Read([]byte) (int, error)
// 	Seek(offset int64, whence int) (int64, error)
// }
//
// os.FileInfo
// type FileInfo interface {
// 	Name() string       // base name of the file
// 	Size() int64        // length in bytes for regular files; system-dependent for others
// 	Mode() FileMode     // file mode bits
// 	ModTime() time.Time // modification time
// 	IsDir() bool        // abbreviation for Mode().IsDir()
// 	Sys() interface{}   // underlying data source (can return nil)
// }

// Ensures proper implementation of interfaces
var _ http.FileSystem = (*_dirHandle)(nil)
var _ http.File = (*fileHandle)(nil)
var _ os.FileInfo = (*_dir)(nil)
var _ os.FileInfo = (*EmbedFile)(nil)

////////////////////////////////////////////////////////////////////////
// DIRECTORY

type _dir struct {
	name    string
	modTime time.Time
	files   map[string]*EmbedFile
	dirs    map[string]*_dir
	sync    sync.Mutex
}

func (d *_dir) Name() string {
	return d.name
}
func (d *_dir) Size() int64 {
	return 0
}
func (d *_dir) Mode() os.FileMode {
	return 0444 | os.ModeDir
}
func (d *_dir) ModTime() time.Time {
	return d.modTime
}
func (d *_dir) IsDir() bool {
	return true
}
func (d *_dir) Sys() interface{} {
	return nil
}
func (d *_dir) Open() (*_dirHandle, error) {
	files := make([]os.FileInfo, 0)
	for _, dir := range d.dirs {
		if dir.name != d.name {
			files = append(files, dir)
		}
	}
	for _, file := range d.files {
		files = append(files, file)
	}
	return &_dirHandle{
		stat:  d,
		files: files,
	}, nil
}

func (dir *_dir) AddFile(file *EmbedFile) {
	dir.sync.Lock()
	dir.files[file.FileName] = file
	dir.sync.Unlock()
}

func (dir *_dir) AddDir(subdir *_dir) {
	dir.sync.Lock()
	dir.dirs[subdir.name] = subdir
	dir.sync.Unlock()
}

type _dirHandle struct {
	stat   *_dir
	offset int
	files  []os.FileInfo // for implementing Readdir
}

func (d *_dirHandle) Open(name string) (handle http.File, err error) {
	name = filepath.Clean(name)
	if filepath.IsAbs(name) {
		name, err = filepath.Rel("/", name)
		if err != nil {
			return
		}
	}

	if name == "." {
		return d, nil
	}

	next := strings.Split(name, string(filepath.Separator))[0]

	if dir, exists := d.stat.dirs[next]; exists {
		if dirHandle, err := dir.Open(); err == nil {
			if p, err := filepath.Rel(next, name); err == nil {
				return dirHandle.Open(p)
			}
		}
		return
	}
	if file, exists := d.stat.files[next]; exists {
		h := &fileHandle{
			stat: file,
			open: true,
		}
		if file.Compressed {
			h.inflater, err = zlib.NewReader(bytes.NewBuffer(h.stat.Data))
		}
		handle = h
		return
	}

	err = errors.New("not found: " + name)
	return
}

func (d *_dirHandle) Readdir(count int) ([]os.FileInfo, error) {
	if count <= 0 {
		return d.files, nil
	}
	if d.offset >= len(d.files) {
		return []os.FileInfo{}, io.EOF
	}
	if d.offset+count > len(d.files) {
		count = len(d.files) - d.offset
	}
	result := d.files[d.offset : d.offset+count]
	d.offset += count

	var err error
	if d.offset > len(d.files) {
		err = io.EOF
	}
	return result, err
}

func (d *_dirHandle) Close() error {
	return nil
}
func (d *_dirHandle) Read(p []byte) (int, error) {
	return 0, errors.New("not file")
}
func (d *_dirHandle) Seek(int64, int) (int64, error) {
	return 0, os.ErrInvalid
}
func (d *_dirHandle) Stat() (os.FileInfo, error) {
	return d.stat, nil
}

////////////////////////////////////////////////////////////////////////
// REGULAR FILE

type EmbedFile struct {
	FileName         string
	Original         string
	Compressed       bool
	Data             []byte
	OriginalSize     int64
	ModificationTime time.Time
}

type fileHandle struct {
	stat     *EmbedFile
	offset   int64
	open     bool
	inflater io.ReadCloser
}

func (f *EmbedFile) Name() string {
	return f.FileName
}

func (f *EmbedFile) Size() int64 {
	return f.OriginalSize
}

func (f *EmbedFile) Mode() os.FileMode {
	return 0444
}

func (f *EmbedFile) ModTime() time.Time {
	return f.ModificationTime
}

func (f *EmbedFile) IsDir() bool {
	return false
}

func (f *EmbedFile) Sys() interface{} {
	return nil
}

func (h *fileHandle) Close() error {
	if h.inflater != nil {
		return h.inflater.Close()
	}
	return nil
}

func (h *fileHandle) Stat() (os.FileInfo, error) {
	return h.stat, nil
}

func (h *fileHandle) Readdir(count int) ([]os.FileInfo, error) {
	return nil, errors.New("not a directory")
}

func (h *fileHandle) Read(buff []byte) (int, error) {
	if h.inflater != nil {
		return h.inflater.Read(buff)
	} else {
		if h.offset >= int64(len(h.stat.Data)) {
			return 0, io.EOF
		}
		n := copy(buff, h.stat.Data[h.offset:])
		h.offset += int64(n)
		return n, nil
	}
}

func (h *fileHandle) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case os.SEEK_SET:
		h.offset = offset
	case os.SEEK_CUR:
		h.offset += offset
	case os.SEEK_END:
		h.offset = h.stat.OriginalSize + offset
	default:
		return 0, os.ErrInvalid
	}
	if h.offset < 0 {
		h.offset = 0
	}
	return h.offset, nil
}
