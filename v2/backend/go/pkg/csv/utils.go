// backend/go/pkg/csv/utils.go
package csv

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"
)

type CSVInfo struct {
    Headers     []string
    SampleData  []map[string]string
    ColumnTypes map[string]string
}

func inferType(value string) string {
    // Try integer
    if _, err := strconv.ParseInt(value, 10, 64); err == nil {
        return "integer"
    }

    // Try float
    if _, err := strconv.ParseFloat(value, 64); err == nil {
        return "float"
    }

    // Try date formats
    dateFormats := []string{
        "2006-01-02",
        "2006/01/02",
        "02-01-2006",
        "02/01/2006",
        "2006-01-02 15:04:05",
        time.RFC3339,
    }

    for _, format := range dateFormats {
        if _, err := time.Parse(format, value); err == nil {
            return "date"
        }
    }

    // Try boolean
    lower := strings.ToLower(value)
    if lower == "true" || lower == "false" || 
       lower == "yes" || lower == "no" ||
       lower == "1" || lower == "0" {
        return "boolean"
    }

    // Default to string
    return "string"
}

func GetCSVInfo(filepath string, sampleSize int) (*CSVInfo, error) {
    file, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    
    // Read headers
    headers, err := reader.Read()
    if err != nil {
        return nil, err
    }

    // Initialize result
    info := &CSVInfo{
        Headers:     headers,
        SampleData:  make([]map[string]string, 0, sampleSize),
        ColumnTypes: make(map[string]string),
    }

    // Read sample rows
    for i := 0; i < sampleSize; i++ {
        record, err := reader.Read()
        if err != nil {
            break
        }

        row := make(map[string]string)
        for j, value := range record {
            row[headers[j]] = value
        }
        info.SampleData = append(info.SampleData, row)
    }

    // Infer column types from first non-empty value
    for _, header := range headers {
        var typeFound bool
        for _, row := range info.SampleData {
            if value := row[header]; value != "" {
                info.ColumnTypes[header] = inferType(value)
                typeFound = true
                break
            }
        }
        if !typeFound {
            info.ColumnTypes[header] = "string"
        }
    }

    return info, nil
}