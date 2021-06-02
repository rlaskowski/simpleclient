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

_, err := fs.Download(url, func(streaminfo simpleclient.StreamInfo) error {
  if streaminfo.WrittenBytes > 0 {
		log.Printf("Written: %.2f percent", streaminfo.ProgressInPercent())
	}

	if streaminfo.Complete {
		log.Print("Complete")
	}

	return nil
})

if err != nil {
    return err
}

