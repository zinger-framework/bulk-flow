package models

import (
	"bufio"
	"examples/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"regexp"
	"strconv"
	"strings"
)

const RegexPattern = `(?:,|\n|^)("(?:(?:"")*[^"]*)*"|[^",\n]*|(?:\n|$))`

var ImportResults = []string{"import_success", "import_error_type", "import_error_reason"}
var Regex = regexp.MustCompile(RegexPattern)

type CsvReader struct {
	Scanner           *bufio.Scanner
	FileName          string
	RawText           string
	Header            map[string]int
	RowNumber         int
	Row               []interface{}
	ImportSuccess     bool
	ImportErrorType   string
	ImportErrorReason string
}

func (csvReader *CsvReader) Scan() bool {
	next := csvReader.Scanner.Scan()
	csvReader.SetRawText(csvReader.Scanner.Text())
	csvReader.AssignRowNumber()
	csvReader.PreprocessText()
	return next
}

func (csvReader *CsvReader) PreprocessText() {
	columns := Regex.FindAllString(csvReader.RawText, -1)
	for index, value := range columns {
		value = strings.TrimSpace(value)
		value = strings.Replace(value, ",", "", 1)
		value = strings.ReplaceAll(value, "\"", "")
		value = strings.ReplaceAll(value, "'", "")

		if strings.Contains(value, ";") {
			csvReader.SetCsvColumn(index, strings.Split(value, ";"))
		} else if govalidator.IsNull(value) {
			csvReader.SetCsvColumn(index, nil)
		} else if govalidator.IsNumeric(value) {
			intValue, _ := strconv.ParseInt(value, 10, 64)
			csvReader.SetCsvColumn(index, intValue)
		} else if govalidator.IsFloat(value) {
			floatValue, _ := strconv.ParseFloat(value, 64)
			csvReader.SetCsvColumn(index, floatValue)
		} else {
			csvReader.SetCsvColumn(index, value)
		}
	}
}

func (csvReader *CsvReader) SetRawText(rawText string) {
	csvReader.RawText = rawText
}

func (csvReader *CsvReader) AssignRowNumber() {
	csvReader.RowNumber += 1
}

func (csvReader *CsvReader) SetCsvColumn(index int, value interface{}) {
	csvReader.Row[index] = value
}

func (csvReader *CsvReader) SetImportResults(success bool, errorType string, errorReason string) {
	csvReader.ImportSuccess = success
	csvReader.ImportErrorType = errorType
	csvReader.ImportErrorReason = errorReason

	// Append import response
	csvReader.SetResults()
}

func (csvReader *CsvReader) CsvColumn(key string) interface{} {
	return csvReader.Row[csvReader.Header[key]]
}

func (csvReader *CsvReader) ToJson() map[string]interface{} {
	json := make(map[string]interface{})
	for key, value := range csvReader.Header {
		json[key] = csvReader.Row[value]
	}
	return json
}

func (csvReader *CsvReader) SetResults(results ...string) {
	if len(results) < 1 {
		results = []string{
			strconv.FormatBool(csvReader.ImportSuccess),
			csvReader.ImportErrorType,
			csvReader.ImportErrorReason,
		}
	}
	result := fmt.Sprintf("%ds/$/,%s/", csvReader.RowNumber+1, strings.Join(results, ","))
	utils.ExecSed(csvReader.FileName, result)
}

func (csvReader *CsvReader) SetResultHeaders() {
	csvReader.SetResults(ImportResults...)
}
