package state

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Lock struct {
	path string
	file *os.File
}

func AcquireLock(dir string) (*Lock, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	path := filepath.Join(dir, "wsl-tunneling.lock")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("another wsl-tunneling operation appears to be running; remove %s if this is stale", path)
		}
		return nil, err
	}

	if _, err := file.WriteString(strconv.Itoa(os.Getpid()) + "\n"); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return nil, err
	}

	return &Lock{path: path, file: file}, nil
}

func (lock *Lock) Release() error {
	if lock == nil {
		return nil
	}
	closeErr := lock.file.Close()
	removeErr := os.Remove(lock.path)
	if closeErr != nil {
		return closeErr
	}
	return removeErr
}
