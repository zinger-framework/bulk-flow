package main

import (
	"examples/models"
	"examples/utils"
	"fmt"
)

func main() {
	fileName := "sample.csv"
	status := map[string]int{}

	status["row_count"] = utils.FetchFileRowCount(fileName)
	status["processed"] = utils.FetchRowCount(fileName, models.ImportProcessed)
	status["succeeded"] = utils.FetchRowCount(fileName, models.ImportSucceeded)
	status["failed"] = utils.FetchRowCount(fileName, models.ImportStopped)

	fmt.Println(status)
}
