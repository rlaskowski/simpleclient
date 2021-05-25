package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type StreamFunc func(fileinfo FileInfo) error

type FileInfo struct {
	Name      string
	TotalSize int64
	Size      int64
	Complete  bool
}

type FileStream struct {
	Filename string
	request  *http.Request
	client   *Client
}

func NewFileStream(filename, url string) *FileStream {
	f := &FileStream{
		Filename: filename,
		client:   NewClient(),
	}

	req, err := f.client.GetRequest(url)
	if err != nil {
		return nil
	}

	f.request = req

	return f
}

func (f *FileStream) Download(fn StreamFunc) error {
	res, err := f.client.Do(f.request)
	if err != nil {
		return err
	}

	err = f.copy(res, fn)
	if err != nil {
		return err
	}

	fn(FileInfo{
		Complete: true,
	})

	return nil
}

func (f *FileStream) AddHeader(key, value string) {
	f.request.Header.Set(key, value)
}

func (f *FileStream) copy(res *Response, fn StreamFunc) error {
	file, err := os.OpenFile(f.Filename, os.O_WRONLY|os.O_CREATE, 0700)
	if err != nil {
		return fmt.Errorf("Could not create file %s due to %s", f.Filename, err)
	}

	defer file.Close()

	finfo := FileInfo{
		Name:      f.Filename,
		TotalSize: res.ContentLength(),
	}

	src := res.Body()

	buff := make([]byte, 32*1024)

	for {
		n, err := src.Read(buff)

		finfo.Size++

		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		if n > 0 {
			file.Write(buff[0:n])

			finfo.Size++
			fn(finfo)
		}
	}

	return nil
}
