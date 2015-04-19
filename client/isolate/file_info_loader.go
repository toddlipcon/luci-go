package isolate

import (
	"crypto/sha1"
	"encoding/gob"
	"io"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"fmt"
	"syscall"
)

type FileInfoLoader struct {
	cache map[shaCacheKey]shaCacheValue
}

type FileInfo struct {
	Path string
	Hash  string
	Mode os.FileMode
	FileSize int64
	LinkDestination string
}

type shaCacheKey struct {
	Inum, Devnum uint64
}

type shaCacheValue struct {
	Mtime syscall.Timespec
	Sha1  string
}


func LoadOrCreateCache() *FileInfoLoader {
	c := newCache()
	cache_file, err := os.Open(cache_path())
	if err != nil {
		e := err.(*os.PathError).Err
		if e == syscall.ENOENT {
			return c
		}
		panic(err)
	}
	defer cache_file.Close()

	parser := gob.NewDecoder(cache_file)
	if err = parser.Decode(&c.cache); err != nil {
		fmt.Println("parsing config file", err.Error())
	}

	return c
}

func (cache *FileInfoLoader) LookupRecursive(path string) ([]*FileInfo, error) {
	type walkEntry struct {
		path string
		file_info os.FileInfo
		err error
	}

	ch := make(chan *walkEntry)
	go func() {
		filepath.Walk(path,
			func(path string, info os.FileInfo, err error) error {
				ch <- &walkEntry{path, info, err}
				return nil
			})
		close(ch)
	}()

	ret := make([]*FileInfo, 0)

	for {
		entry, available := <-ch
		if ! available {
			break
		}
		if entry.err != nil {
			return nil, entry.err
		}

		s, err := cache.LookupInfo(entry.path, entry.file_info)
		if err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
}

func (cache *FileInfoLoader) LookupInfoForPath(path string) (*FileInfo, error) {
	finfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return cache.LookupInfo(path, finfo)
}

func (cache *FileInfoLoader) LookupInfo(path string, fileinfo os.FileInfo) (*FileInfo, error) {
	stat := fileinfo.Sys().(*syscall.Stat_t)

	key := shaCacheKey{Inum: stat.Ino, Devnum: stat.Dev}
	result, found_in_cache := cache.cache[key]
	if !found_in_cache || result.Mtime != stat.Mtim {
		s := sha1_file(path)
		result = shaCacheValue{Mtime: stat.Mtim, Sha1: s}
		cache.cache[key] = result
	}

	ret := &FileInfo{
		Path: path,
		Hash: result.Sha1,
		Mode: fileinfo.Mode(),
		FileSize: fileinfo.Size()}
	// TODO: handle symlinks

	return ret, nil
}

func (c* FileInfoLoader) Save() {
	cache_file, err := os.Create(cache_path())
	if err != nil {
		panic(err)
	}
	defer cache_file.Close()

	enc := gob.NewEncoder(cache_file)
	enc.Encode(c.cache)
}


func newCache() *FileInfoLoader {
	return &FileInfoLoader{
		cache: make(map[shaCacheKey]shaCacheValue),
	}
}

func cache_path() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return path.Join(usr.HomeDir, ".isolate-sha-cache.json")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func sha1_file(path string) string {
	hash := sha1.New()
	f, err := os.Open(path)
	check(err)
	defer f.Close()

	io.Copy(hash, f)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

