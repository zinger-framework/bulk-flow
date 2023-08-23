package models

import (
	"bufio"
	"encoding/json"
	"examples/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const RegexPattern = `(?:,|\n|^)("(?:(?:"")*[^"]*)*"|[^",\n]*|(?:\n|$))`
const (
	VALIDATION_ERROR  = "validation"
	RATELIMIT_ERROR   = "ratelimit"
	DESTINATION_ERROR = "destination"
)

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

	// Append import response to file
	results := []string{
		strconv.FormatBool(csvReader.ImportSuccess),
		csvReader.ImportErrorType,
		csvReader.ImportErrorReason,
	}
	result := fmt.Sprintf("%ds/$/,%s/", csvReader.RowNumber+1, strings.Join(results, ","))
	utils.ExecSed(csvReader.FileName, result)
}

func (csvReader *CsvReader) CsvColumn(key string) interface{} {
	return csvReader.Row[csvReader.Header[key]]
}

func (csvReader *CsvReader) ToJson() map[string]interface{} {
	jsonData := make(map[string]interface{})
	for key, value := range csvReader.Header {
		jsonData[key] = csvReader.Row[value]
	}
	return jsonData
}

func (csvReader *CsvReader) SetResultHeaders() {
	result := fmt.Sprintf("%ds/$/,%s/", csvReader.RowNumber+1, strings.Join(ImportResults, ","))
	utils.ExecSed(csvReader.FileName, result)
}

func (csvReader *CsvReader) SendRequest() *http.Response {
	const url = "https://webhook.site/c5918e84-4de2-4f02-9c97-ffbd908b6dd7"
	const contentType = "application/json"
	const reqMethod = "POST"

	body, err := json.Marshal(csvReader.ToJson())
	utils.PanicError(err)
	fmt.Printf("%s\n", body)

	request, err := http.NewRequest(reqMethod, url, strings.NewReader(string(body)))
	request.Header.Set("Content-Type", contentType)
	utils.PanicError(err)

	response, err := http.DefaultClient.Do(request)
	utils.PanicError(err)

	return response
}

func (csvReader *CsvReader) HandleResponse(response *http.Response) {
	if response.StatusCode/100 == 2 {
		// success
		csvReader.SetImportResults(true, "", "")
	} else {
		// failure
		csvReader.SetImportResults(false, RATELIMIT_ERROR, response.Status)
	}
}
