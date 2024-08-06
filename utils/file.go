package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"syscall"
)

// lock file to prevent concurrent read/writes
func LoadFile(filepath string) (*os.File, error) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading")
	}

	// Exclusive lock obtained on the file descriptor
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, err
	}

	return f, nil
}

// unlock the file
func CloseFile(f *os.File) error {
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return f.Close()
}

// Read file
func ReadTasks(file *os.File) ([][]string, error) {
	// Reset file offset to the beginning
	if _, err := file.Seek(0, 0); err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	tasks, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
