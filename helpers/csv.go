package helpers

import (
	"encoding/csv"
	"io"
	"os"
)

// ParseCSV reads CSV data from an io.Reader and returns a channel of string slices.
// Each slice represents one row of the CSV file.
//
// The function uses a channel for memory-efficient processing of large CSV files,
// allowing you to process rows one at a time.
//
// Example:
//
//	file, _ := os.Open("data.csv")
//	defer file.Close()
//
//	records, err := helpers.ParseCSV(file)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for record := range records {
//		fmt.Println("Row:", record)
//	}
//
//	// With HTTP upload
//	func UploadHandler(w http.ResponseWriter, r *http.Request) {
//		file, _, _ := r.FormFile("csv")
//		records, _ := helpers.ParseCSV(file)
//		for row := range records {
//			processRow(row)
//		}
//	}
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

// ReadFile reads the entire contents of a file and returns it as a byte slice.
//
// Example:
//
//	content, err := helpers.ReadFile("config.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	var config Config
//	json.Unmarshal(content, &config)
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
