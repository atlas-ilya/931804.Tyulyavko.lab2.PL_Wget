package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// WriteCounter counts the number of bytes written to it. By implementing the Write method,
// it is of the io.Writer interface and we can pass this into io.TeeReader()
type WriteCounter struct {
	Length, Total int64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {

	n := len(p)
	wc.Total += int64(n)
	return n, nil

}

func main() {

	fmt.Println("Download Started")

	url := "http://ovh.net/files/10Gb.dat"
	err := DownloadFile(url)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println("Download Finished")
}

// DownloadFile will download a url and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
// Bytes are output to the console
// We pass an io.TeeReader into Copy() to report progress on the download.
func DownloadFile(url string) error {

	out, err := os.Create(url[strings.LastIndex(url, "/")+1 : len(url)])

	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s \n", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: status: %s \n", resp.Status)
		return err
	}

	counter := WriteCounter{
		Length: resp.ContentLength,
	}

	go func() {
		i := 0
		for {
			i++
			time.Sleep(time.Second)
			fmt.Printf(
				"\r %d sec |Downloading... %2d Bytes of %d Bytes complete|",
				i, int64(counter.Total), int64(counter.Length))
		}
	}()

	_, err = io.Copy(out, io.TeeReader(resp.Body, &counter))

	if err != nil {
		return err
	}

	return nil
}
