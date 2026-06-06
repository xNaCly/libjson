//go:build unix

package libjson

import (
	"fmt"
	"os"
	"syscall"
)

func mmapFile(f *os.File) ([]byte, func() error, error) {
	info, err := f.Stat()
	if err != nil {
		return nil, nil, err
	}

	size := info.Size()
	if size == 0 {
		return []byte{}, nil, nil
	}
	if size < 0 {
		return nil, nil, fmt.Errorf("invalid file size %d", size)
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE)
	if err != nil {
		return nil, nil, err
	}

	return data, func() error {
		return syscall.Munmap(data)
	}, nil
}
