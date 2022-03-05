package main

import (
	"bytes"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"golang.org/x/sync/errgroup"
)

const tmpDir = "tmp"
const dstDir = "dst"

// hash/maphashを使って、リクエスト先のURL文字列からハッシュ値を計算、ハッシュ値-indexみたいなファイル名にしたら将来途中から開始機能が実装される時にやりやすそうでいいかも
// 該当のファイルが存在して、かつファイルサイズが適切であれば再ダウンロードしない
// https://pkg.go.dev/hash/maphash
func main() {
	url := flag.String("url", "", "")
	flag.Parse()

	prefix := createHash(*url)

	resp, err := http.Head(*url)
	if err != nil {
		panic("http.NewRequest")
	}
	ar := resp.Header.Get("Accept-Ranges")
	if ar == "" {
		panic("not accept range access")
	}

	cl, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		panic("strconv.Atoi")
	}
	nc := runtime.NumCPU()
	var chunkBytes int
	if cl%nc == 0 {
		chunkBytes = cl / nc
	} else {
		chunkBytes = cl/nc + 1
	}

	var eg errgroup.Group
	for i := 0; i < nc; i++ {
		offset := chunkBytes * i
		var limit int
		if cl <= chunkBytes {
			limit = cl
		} else {
			limit = chunkBytes * (i + 1)
		}
		// goroutineとforの罠回避変数定義
		// https://qiita.com/sudix/items/67d4cad08fe88dcb9a6d
		index := i
		eg.Go(func() error {
			b, err := download(url, offset, limit)
			if err != nil {
				return err
			}

			err = saveToTmpFile(b, index, tmpDir, prefix)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		fmt.Println(err)
		panic("eg.Wait()")
	}
	mergeTmpFiles(dstDir, tmpDir, prefix)
}

func download(url *string, offset, limit int) ([]byte, error) {
	client := &http.Client{}
	r, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		panic("http.NewRequest")
	}
	r.Header.Set("Range", ("bytes=" + strconv.Itoa(offset) + "-" + strconv.Itoa(limit)))
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return b, nil
}

func saveToTmpFile(b []byte, index int, tmpDir, prefix string) error {
	bb := bytes.NewReader(b)

	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return err
	}
	f, err := os.CreateTemp(tmpDir, createTmpFilename(prefix, index))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, bb)
	if err != nil {
		return err
	}
	return nil
}

func createHash(seed string) string {
	h := md5.New()
	io.WriteString(h, seed)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func createTmpFilename(prefix string, index int) string {
	return fmt.Sprintf("%v-%d", prefix, index)
}

func mergeTmpFiles(dstDir, tmpDir, prefix string) (string, error) {
	paths, err := filepath.Glob(filepath.Join(tmpDir, prefix))
	if err != nil {
		return "", err
	}
	var files []io.Reader
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return "", err
		}

		defer f.Close()

		files = append(files, f)
	}
	b, err := io.ReadAll(io.MultiReader(files...))
	if err != nil {
		return "", err
	}
	f, err := os.Create(createFilename())
	if err != nil {
		return "", err
	}

	if _, err = io.Copy(f, bytes.NewReader(b)); err != nil {
		return "", err
	}
	return f.Name(), nil
}

func createFilename() string {
	return "example-filename"
}
