package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("start.")
	paths := getfilePaths("./src")

	for _, p := range paths {
		i, err := extractImage(p)
		if err != nil {
			fmt.Println(err)
			fmt.Println("skip this file")
			continue
		}

		b, err := convert(i)
		if err != nil {
			fmt.Println(err)
			panic("failed to convert")
		}

		_, err = save(b, strings.Split(filepath.Base(p), ".")[0])
		if err != nil {
			fmt.Println(err)
			panic("failed to save")
		}
	}
	fmt.Println("fin.")
}

func getfilePaths(root string) []string {
	var files []string

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files
}

func extractImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	i, format, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	if format != "jpeg" {
		return nil, fmt.Errorf("format error: expected jpeg, actual %v", format)
	}

	return i, nil
}

func convert(i image.Image) (*bytes.Buffer, error) {
	var b bytes.Buffer

	if err := png.Encode(&b, i); err != nil {
		return nil, err
	}

	return &b, nil
}

func save(b *bytes.Buffer, name string) (string, error) {
	dst := filepath.Join("./dst", name) + ".png"
	f, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, b)
	if err != nil {
		return "", err
	}

	return dst, nil
}
