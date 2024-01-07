package utils

import (
	"io"
	"log"
	"os"
)

func ReadFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
