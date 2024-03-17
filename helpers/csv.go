package helpers

import (
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/ochom/gutils/gttp"
)

// ParseCSV reads csv from io.Reader and returns a chan of slice of strings
func ParseCSV(file io.Reader) (chan []string, error) {
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	result := make(chan []string)
	go func() {
		defer close(result)
		for _, record := range records {
			result <- record
		}
	}()

	return result, nil
}

// ReadFile ...
func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func DownloadFile(url string) ([]byte, error) {
	resp, err := gttp.Get(url, nil, 5*time.Minute)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
