package models

import (
	"bufio"
	"github.com/asaskevich/govalidator"
	"regexp"
	"strconv"
	"strings"
)

const RegexPattern = `(?:,|\n|^)("(?:(?:"")*[^"]*)*"|[^",\n]*|(?:\n|$))`

var Regex = regexp.MustCompile(RegexPattern)

type CsvReader struct {
	Scanner *bufio.Scanner
	Header  map[string]int
	Row     []interface{}
}

func (csvReader *CsvReader) Scan() bool {
	next := csvReader.Scanner.Scan()
	csvReader.PreprocessRow(csvReader.Scanner.Text())
	return next
}

func (csvReader CsvReader) PreprocessRow(row string) {
	columns := Regex.FindAllString(row, -1)
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

func (csvReader CsvReader) SetCsvColumn(index int, value interface{}) {
	csvReader.Row[index] = value
}

func (csvReader CsvReader) CsvColumn(key string) interface{} {
	return csvReader.Row[csvReader.Header[key]]
}

func (csvReader CsvReader) ToJson() map[string]interface{} {
	json := make(map[string]interface{})
	for key, value := range csvReader.Header {
		json[key] = csvReader.Row[value]
	}
	return json
}
