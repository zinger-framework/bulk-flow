package main

import (
	"examples/models"
	"examples/utils"
	"fmt"
	"strings"
)

const MinBatchSize = 10

func main() {
	fileName := "sample.csv"

	// Can be client side call - Convert file from DOS to Unix format
	utils.ReplacePattern(fileName, "s/\r$//")

	_ = utils.FetchHeaderFromFile(fileName)

	// Append import result headers
	result := fmt.Sprintf("%ds/$/,%s/", 1, strings.Join(models.ResultHeaders, ","))
	utils.ReplacePattern(fileName, result)

	rowCount := utils.FetchFileRowCount(fileName)
	tps := 4 // default - 1000

	batchSize := (rowCount / tps) + 1
	if batchSize < MinBatchSize {
		batchSize = MinBatchSize
	}

	var positions []models.Position
	start := 1

	for start < rowCount {
		end := start + batchSize
		if end > rowCount {
			end = rowCount
		}

		positions = append(positions, models.Position{Start: start + 1, End: end})
		start = end
	}

	for _, position := range positions {
		// Trigger lambda in async mode
		fmt.Println(position)
	}
}
