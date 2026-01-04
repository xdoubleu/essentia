package httptools

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
)

// WriteCSV writes the provided data as a CSV file with
// the provided filename to a [http.ResponseWriter].
func WriteCSV(
	w http.ResponseWriter,
	filename string,
	headers []string,
	data [][]string,
) error {
	output := [][]string{}
	output = append(output, headers)
	output = append(output, data...)

	w.Header().Set("content-type", "text/csv")
	w.Header().
		Set("content-disposition", fmt.Sprintf("attachment;filename=%s.csv", filename))

	csvWriter := csv.NewWriter(w)

	return csvWriter.WriteAll(output)
}

// ReadCSV reads the returned CSV file from a [http.Response.Body].
func ReadCSV(body io.Reader) ([][]string, error) {
	csvReader := csv.NewReader(body)
	csvData, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return csvData, nil
}
