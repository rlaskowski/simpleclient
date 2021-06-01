package simpleclient

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
	Checksum     string
	TotalSize    int64
	WrittenBytes int64
	Complete     bool
}

type FileStream struct {
	path    string
	headers http.Header
	Client  *Client
}

func NewFileStream(path string) *FileStream {
	return &FileStream{
		path:    path,
		headers: make(http.Header),
		Client:  NewClient(),
	}
}

func (f *FileStream) Download(url string, fn StreamFunc) (*Response, error) {
	req, err := f.Client.GetRequest(url)
	if err != nil {
		return nil, err
	}

	req.Header = f.headers

	res, err := f.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.response.Body.Close()

	err = f.copy(res, fn)
	if err != nil {
		return nil, err
	}

	checksum := f.checksum(res.response.Request.URL)

	fn(FileInfo{
		Complete: true,
		Checksum: checksum,
	})

	return res, nil
}

func (f *FileStream) SetHeader(key, value string) {
	f.headers.Set(key, value)
}

func (f *FileStream) stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (f *FileStream) filepath(res *Response) (string, error) {
	url := res.URL()

	if !strings.Contains(url.Path, "/") {
		return "", errors.New("Can't find filename from url path")
	}

	filename := filepath.Base(url.Path)

	stat, err := f.stat(f.path)
	if err != nil || !stat.IsDir() {
		return "", errors.New("Incorrect path to store file")
	}

	return filepath.Join(f.path, filename), nil
}

func (f *FileStream) checksum(url *url.URL) string {
	return url.Query().Get("checksum")
}

func (f *FileStream) copy(res *Response, fn StreamFunc) error {
	path, err := f.filepath(res)
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

func (fi FileInfo) Progress() float64 {
	if !(fi.TotalSize > 0) {
		return 0
	}

	return float64(fi.WrittenBytes) / float64(fi.TotalSize)
}

func (fi FileInfo) ProgressInPercent() float64 {
	return fi.Progress() * 100
}
