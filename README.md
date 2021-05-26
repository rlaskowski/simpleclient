# simpleclient

### Install
```
go get github.com/rlaskowski/simpleclient
```

### Example

This example show how to download file from url

```
fs := client.NewFileStream("test.zip", "http://ipv4.download.thinkbroadband.com/200MB.zip")

err := fs.Download(func(fileinfo client.FileInfo) error {
  if fileinfo.Size > 0 {
		log.Printf("Written: %v bytes", fileinfo.Size)
	}

	if fileinfo.Complete {
		log.Print("Complete")
	}

	return nil
})

if err != nil {
    return err
}

