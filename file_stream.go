package simpleclient

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultWriteBuffer int = 32 * 1024
)

type StreamFunc func(fileinfo StreamInfo) error

type StreamInfo struct {
	Name         string
	TotalSize    int64
	WrittenBytes int64
	Complete     bool
}

type FileStream struct {
	path        string
	headers     http.Header
	writeBuffer int
	Client      *Client
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

	return res, nil
}

func (f *FileStream) WriteBuffer(buff int) {
	f.writeBuffer = buff
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

	sinfo := StreamInfo{
		Name:      path,
		TotalSize: res.ContentLength(),
	}

	src := res.Body()

	if !(f.writeBuffer > 0) {
		f.WriteBuffer(defaultWriteBuffer)
	}

	buff := make([]byte, f.writeBuffer)

	for {
		n, err := src.Read(buff)

		if n > 0 {
			file.Write(buff[0:n])

			sinfo.WrittenBytes += int64(n)
			fn(sinfo)
		}

		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
	}

	sinfo.Complete = true

	return fn(sinfo)
}

func (fi StreamInfo) Progress() float64 {
	if !(fi.TotalSize > 0) {
		return 0
	}

	return float64(fi.WrittenBytes) / float64(fi.TotalSize)
}

func (fi StreamInfo) ProgressInPercent() float64 {
	return fi.Progress() * 100
}
