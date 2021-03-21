package restaurants

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// reads a CSV give file by a given number of columns
func readCSVFile(file string, numOfColumns int) ([][]string, error) {
	csvfile, err := os.Open(filepath.Clean(file))
	if err != nil {
		return nil, fmt.Errorf("failed to open csv file '%v': [%v]", file, err)
	}

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = numOfColumns

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv file '%v': [%v]", file, err)
	}

	err = csvfile.Close()
	if err != nil {
		return nil, err
	}
	return rawCSVdata, nil
}

// strToInt parses an integer from a given string
func strToInt(s string) (int, error) {
	num, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("failed converting ID into an integer: '%v'", err)
	}
	return int(num), err
}

// will remove common typos and unwanted text that is usually added to the csv by mistake
func cleanCommonTypos(inputStr string) string {
	commonTypos := []string{
		"Click to check domain availability.",
	}

	for _, v := range commonTypos {
		return strings.Replace(inputStr, v, "", -1)
	}

	return ""
}
