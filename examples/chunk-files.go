package main

import "fmt"

type Position struct {
	Start int
	End   int
}

const MinBatchSize = 10

func main() {
	rowCount := 100
	tps := 4 // default - 1000

	batchSize := (rowCount / tps) + 1
	if batchSize < MinBatchSize {
		batchSize = MinBatchSize
	}

	var positions []Position
	start := 0

	for start < rowCount {
		end := start + batchSize
		if end > rowCount {
			end = rowCount
		}

		positions = append(positions, Position{Start: start + 1, End: end})
		start = end
	}

	fmt.Println(positions)
}
