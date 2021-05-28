package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type StreamFunc func(fileinfo FileInfo) error

type FileInfo struct {
	Name         string
	TotalSize    int64
	WrittenBytes int64
	Complete     bool
}

type FileStream struct {
	path    string
	headers http.Header
	client  *Client
}

func NewFileStream(path string) *FileStream {
	return &FileStream{
		path:    path,
		headers: make(http.Header),
		client:  NewClient(),
	}
}

func (f *FileStream) Download(url string, fn StreamFunc) (*Response, error) {
	req, err := f.client.GetRequest(url)
	if err != nil {
		return nil, err
	}

	req.Header = f.headers

	res, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = f.copy(res, fn)
	if err != nil {
		return nil, err
	}

	fn(FileInfo{
		Complete: true,
	})

	return res, nil
}

func (f *FileStream) SetHeader(key, value string) {
	f.headers.Set(key, value)
}

func (f *FileStream) stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (f *FileStream) filepath(url *url.URL) (string, error) {
	if !strings.Contains(url.Path, "/") {
		return "", errors.New("Can't find filename from url path")
	}

	filename := filepath.Base(url.Path)

	stat, err := f.stat(f.path)
	if err != nil || !stat.IsDir() {
		return "", errors.New("Incorrect file path")
	}

	return filepath.Join(f.path, filename), nil
}

func (f *FileStream) copy(res *Response, fn StreamFunc) error {
	url := res.URL()

	path, err := f.filepath(url)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0700)
	if err != nil {
		return fmt.Errorf("Could not create file %s due to %s", filepath.Base(path), err)
	}

	defer file.Close()

	finfo := FileInfo{
		Name:      path,
		TotalSize: res.ContentLength(),
	}

	src := res.Body()

	buff := make([]byte, 32*1024)

	for {
		n, err := src.Read(buff)

		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		if n > 0 {
			file.Write(buff[0:n])

			finfo.WrittenBytes += int64(n)
			fn(finfo)
		}
	}

	return nil
}
