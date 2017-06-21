package filesys

import (
	"io/ioutil"
	"os"
	"time"

	"crypto/md5"
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/ghongli/cache"
)

const (
	defaultFilePerm          os.FileMode = os.ModePerm
	defaultDirectoryFilePerm             = 0755
)

// a cache for file storage.
type Cache struct {
	baseDir string
	s       cache.Serializer
}

// MustNewFileSysStore an initialized Filesystem Cache
// If a non-existent directory is passed, it would be created automatically.
// Panics if the directory could not be created
func MustNewFileSysStore(baseDir string, interval time.Duration) *Cache {
	_, err := os.Stat(baseDir)
	if err != nil {
		// Directory baseDir does not exist, create it
		if err := createDirectory(baseDir); err != nil {
			panic(fmt.Errorf("Base directory[%s] could not be created : %v", baseDir, err))
		}

	}

	fs := &Cache{baseDir: baseDir, s: cache.NewCacheSerializer()}

	fs.TrashGc(interval)

	return fs
}

// Put value into file cache.
// expire means how long to keep this file, unit of ms.
// if expire equals 0, cache this item forever.
func (f *Cache) Put(key string, data []byte, expire time.Duration) error {
	path := f.filePathForKey(key)

	if err := createDirectory(filepath.Dir(path)); err != nil {
		return err
	}

	item := cache.CacheItem{
		CreatedTime: time.Now(),
		Expired:     expire,
		Data:        data,
	}

	b, err := f.s.Serialized(item)
	if err != nil {
		return err
	}

	return writeFile(path, b)
}

// Get value from file cache.
// if non-exist or expired, return nil, error.
// run garbage collection with the key if necessary
func (f *Cache) Get(key string) ([]byte, error) {
	b, err := readFile(f.filePathForKey(key))
	if err != nil {
		return nil, err
	}

	i := new(cache.CacheItem)
	err = f.s.DeSerialized(b, i)
	if err != nil {
		return nil, err
	}

	if i.IsExpired() {
		f.Delete(key)
		return nil, cache.ErrCacheMiss
	}

	return i.Data, nil
}

// Delete file cache value.
func (f *Cache) Delete(key string) error {
	fileName := f.filePathForKey(key)
	if ok, _ := exists(fileName); ok {
		return os.Remove(fileName)
	}
	return nil
}

// ClearAll will clean cached files.
func (f *Cache) ClearAll() error {
	return os.RemoveAll(f.baseDir)
}

// IsExist check value is exist.
func (f *Cache) IsExist(key string) bool {
	ret, _ := exists(f.filePathForKey(key))
	return ret
}

func (f *Cache) TrashGc(interval time.Duration) {
	filepath.Walk(
		f.baseDir,
		func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if fileInfo.IsDir() {
				return nil
			}

			b, err := readFile(path)
			if err != nil {
				return err
			}

			currentItem := new(cache.CacheItem)
			if err := f.s.DeSerialized(b, currentItem); err != nil {
				return err
			}

			if currentItem.IsExpired() {
				if err := os.Remove(path); !os.IsExist(err) {
					return err
				}
			}

			return nil
		},
	)

	time.AfterFunc(interval, func() {
		f.TrashGc(interval)
	})
}

// Gets a unique path for a cache key.
// to a directory 3 level deep. Something like "basedir/33/rr/33/hash"
func (f *Cache) filePathForKey(key string) string {
	hashSum := md5.Sum([]byte(key))
	hashSumAsString := hex.EncodeToString(hashSum[:])

	return filepath.Join(f.baseDir,
		string(hashSumAsString[0:2]),
		string(hashSumAsString[2:4]),
		string(hashSumAsString[4:6]), hashSumAsString)
}

// get bytes to file.
func readFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

// put bytes to file.
// if non-exist, create this file.
func writeFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, defaultFilePerm)
}

// check file exist.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func createDirectory(dir string) error {
	return os.MkdirAll(dir, defaultDirectoryFilePerm)
}
