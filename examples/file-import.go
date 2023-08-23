package main

import (
	"bufio"
	"examples/models"
	"examples/utils"
	"os"
	"strings"
)

func CsvScanner(fileName string) models.CsvReader {
	// Can be client side call - Convert file from DOS to Unix format
	utils.ExecSed(fileName, "s/\r$//")

	file, err := os.Open(fileName)
	utils.PanicError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	text := scanner.Text()
	headers := strings.Split(text, ",")
	headerMap := map[string]int{}
	for index, header := range headers {
		headerMap[header] = index
	}

	return models.CsvReader{
		Scanner:   scanner,
		FileName:  fileName,
		RawText:   text,
		Header:    headerMap,
		RowNumber: 0,
		Row:       make([]interface{}, len(headerMap)),
	}
}

func main() {
	scanner := CsvScanner("sample.csv")

	// Append import result headers
	scanner.SetResultHeaders()

	for scanner.Scan() {
		response := scanner.SendRequest()
		scanner.HandleResponse(response)
	}
}
