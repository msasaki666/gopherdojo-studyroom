package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/gif"
	_ "image/gif"
	"image/jpeg"
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
	from := flag.String("from", "jpeg", "")
	to := flag.String("to", "png", "")
	flag.Parse()

	if !isValidFlag(from, []string{"jpeg", "png", "gif"}) {
		panic("invalid flag")
	}

	if !isValidFlag(to, []string{"jpeg", "png", "gif"}) {
		panic("invalid flag")
	}

	paths := getfilePaths("./src")

	for _, p := range paths {
		i, err := extractImage(p, from)
		if err != nil {
			fmt.Println("skip this file")
			continue
		}

		b, err := convert(i, to)
		if err != nil {
			fmt.Println(err)
			panic("failed to convert")
		}
		filename := strings.Split(filepath.Base(p), ".")[0] + "." + *to
		path, err := save(b, filename)
		if err != nil {
			fmt.Println(err)
			panic("failed to save")
		}
		fmt.Println(path, "is created")
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

func extractImage(path string, from *string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	i, format, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	if format != *from {
		return nil, fmt.Errorf("format error: expected %v, actual %v", *from, format)
	}

	return i, nil
}

func convert(i image.Image, to *string) (*bytes.Buffer, error) {
	var b bytes.Buffer

	switch *to {
	case "jpeg":
		if err := jpeg.Encode(&b, i, nil); err != nil {
			return nil, err
		}
	case "png":
		if err := png.Encode(&b, i); err != nil {
			return nil, err
		}
	case "gif":
		if err := gif.Encode(&b, i, nil); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid format")
	}

	return &b, nil
}

func save(b *bytes.Buffer, name string) (string, error) {
	dst := filepath.Join("./dst", name)
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

func isValidFlag(flag *string, validList []string) bool {
	var valid bool
	for _, v := range validList {
		if *flag == v {
			valid = true
			break
		}
	}
	return valid
}
