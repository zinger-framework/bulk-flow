package main

import (
	"bufio"
	"examples/models"
	"examples/utils"
	"os"
)

const MaxBatchSize = 50

func ProcessLambda(srcFileName string, position models.Position, headerMap map[string]int, rowCount int) {
	destFileName := "tmp/output.csv"
	endPosition := position.End
	if endPosition > MaxBatchSize {
		endPosition = MaxBatchSize
	}

	utils.CreateTempFile(srcFileName, destFileName, position.Start, endPosition)
	file, err := os.Open(destFileName)
	utils.PanicError(err)
	defer os.Remove(destFileName)
	defer file.Close()

	scanner := models.CsvReader{
		Scanner:   bufio.NewScanner(file),
		FileName:  srcFileName,
		Header:    headerMap,
		RowNumber: position.Start - 1,
		Row:       make([]interface{}, len(headerMap)),
	}

	for scanner.Scan() {
		response := scanner.SendRequest()
		scanner.HandleResponse(response)
	}

	if scanner.RowNumber < position.End {
		position.Start = scanner.RowNumber + 1
		// Trigger next set of lambda in async mode
	} else if utils.FetchRowCount(srcFileName, models.ImportProcessed) == rowCount {
		// Upload srcFile to S3.
		// Delete srcFile from server if upload success.
	}
}

func main() {
	fileName := "sample.csv"
	position := models.Position{Start: 6, End: 10}
	headerMap := utils.FetchHeaderFromFile(fileName)
	rowCount := utils.FetchFileRowCount(fileName)

	ProcessLambda(fileName, position, headerMap, rowCount)
}
