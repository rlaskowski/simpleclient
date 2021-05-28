# simpleclient

### Install
```
go get github.com/rlaskowski/simpleclient
```

### Example

This example show how to download file from url

```go
fs := simpleclient.NewFileStream("/home/examplepath")

url := "http://ipv4.download.thinkbroadband.com/200MB.zip"

_, err := fs.Download(url, func(fileinfo simpleclient.FileInfo) error {
  if fileinfo.WrittenBytes > 0 {
		log.Printf("Written: %v bytes", fileinfo.WrittenBytes)
	}

	if fileinfo.Complete {
		log.Print("Complete")
	}

	return nil
})

if err != nil {
    return err
}

