package file

import (
	"encoding/json"
	"errors"
	"go-backend/storage"
	"io"
	"os"
	"sync"
)

// FileStore implements DataStore using a JSON file for persistence.
type FileStore struct {
	mu       sync.RWMutex
	filePath string
	data     storage.Data
}

// NewFileStore initializes a FileStore and loads data from disk.
func NewFileStore(path string) (*FileStore, error) {
	fs := &FileStore{filePath: path}
	if err := fs.load(); err != nil {
		return nil, err
	}
	return fs, nil
}

func (fs *FileStore) load() error {
	file, err := os.Open(fs.filePath)
	if errors.Is(err, os.ErrNotExist) {
		fs.data = storage.Data{}
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&fs.data); err != nil {
		if errors.Is(err, io.EOF) {
			// Empty file, then start with empty data
			fs.data = storage.Data{}
			return nil
		}
		return err
	}

	return nil
}

func (fs *FileStore) save() error {
	tmp := fs.filePath + ".tmp"

	file, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(file).Encode(fs.data); err != nil {
		file.Close()
		return err
	}
	file.Close()

	return os.Rename(tmp, fs.filePath)
}

func (fs *FileStore) IsHealthy() bool {
	_, err := os.Stat(fs.filePath)
	return err == nil || os.IsNotExist(err)
}
