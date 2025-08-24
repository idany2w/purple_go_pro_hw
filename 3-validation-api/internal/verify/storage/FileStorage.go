package storage

import (
	"fmt"
	"os"
)

const (
	hashPath = "./data/hashes"
)

type FileStorage struct{}

func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

func (s *FileStorage) SaveHash(hash string) error {

	file, err := os.Create(fmt.Sprintf("%s/%s", hashPath, hash))

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write([]byte(hash))

	return err
}

func (s *FileStorage) RemoveHash(hash string) error {
	return os.Remove(fmt.Sprintf("%s/%s", hashPath, hash))
}

func (s *FileStorage) GetHash(hash string) (string, error) {
	content, err := os.ReadFile(fmt.Sprintf("%s/%s", hashPath, hash))

	if err != nil {
		return "", err
	}

	return string(content), nil
}
