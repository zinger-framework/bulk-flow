package utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

func ReplacePattern(fileName string, result string) {
	err := exec.Command("sed", "-i", "", result, fileName).Start()
	PanicError(err)
}

func CreateTempFile(srcFileName string, destFileName string, start int, end int) {
	sed := fmt.Sprintf("sed -n %d,%dp %s > %s", start, end, srcFileName, destFileName)
	_, err := exec.Command("bash", "-c", sed).Output()
	PanicError(err)
}

func FetchHeaderFromFile(fileName string) map[string]int {
	headerText, err := exec.Command("head", "-1", fileName).Output()
	PanicError(err)

	headers := strings.Split(string(headerText), ",")
	headerMap := map[string]int{}
	for index, header := range headers {
		headerMap[header] = index
	}
	return headerMap
}

func ToInteger(result string) int {
	value, err := strconv.ParseInt(strings.TrimSpace(result), 10, 0)
	PanicError(err)

	return int(value)
}

func FetchFileRowCount(fileName string) int {
	result, err := exec.Command("sed", "-n", "$=", fileName).Output()
	PanicError(err)

	return ToInteger(string(result))
}

func IsFileFullyProcessed(fileName string, rowCount int) bool {
	result, err := exec.Command("grep", "-c", ",import_", fileName).Output()
	PanicError(err)

	return ToInteger(string(result)) == rowCount
}
