package main

import (
	"bufio"
	"encoding/json"
	"examples/models"
	"examples/utils"
	"fmt"
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

	const url = "https://webhook.site/c5918e84-4de2-4f02-9c97-ffbd908b6dd7"
	const contentType = "application/json"

	for scanner.Scan() {
		body, err := json.MarshalIndent(scanner.ToJson(), "", "  ")
		utils.PanicError(err)
		fmt.Printf("%s\n", body)

		//http.Post(url, contentType, body)
		scanner.SetImportResults(false, "ratelimit", "429 Too Many Requests")
	}
}
