package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"examples/models"
	"examples/utils"
	"fmt"
	"os"
	"strings"
)

func CsvScanner(scanner *bufio.Scanner) (models.CsvReader, error) {
	next := scanner.Scan()
	if !next {
		return models.CsvReader{}, errors.New("reached EOF")
	}

	headers := strings.Split(scanner.Text(), ",")
	headerMap := map[string]int{}
	for index, header := range headers {
		headerMap[header] = index
	}

	return models.CsvReader{
		Scanner: scanner,
		Header:  headerMap,
		Row:     make([]interface{}, len(headerMap)),
	}, nil
}

func main() {
	file, err := os.Open("sample.csv")
	utils.PanicError(err)

	defer file.Close()

	scanner, err := CsvScanner(bufio.NewScanner(file))
	utils.PanicError(err)

	const url = "https://webhook.site/c5918e84-4de2-4f02-9c97-ffbd908b6dd7"
	const contentType = "application/json"

	for scanner.Scan() {
		//record := scanner.CsvColumn("name")
		//fmt.Println(record)

		body, err := json.MarshalIndent(scanner.ToJson(), "", "  ")
		utils.PanicError(err)
		fmt.Printf("%s\n", body)

		//http.Post(url, contentType, body)
	}

	err = scanner.Scanner.Err()
	utils.PanicError(err)
}
