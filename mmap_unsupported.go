//go:build !unix

package libjson

import (
	"io"
	"os"
)

func mmapFile(f *os.File) ([]byte, func() error, error) {
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}
	return data, nil, nil
}
