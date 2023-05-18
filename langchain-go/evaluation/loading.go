package evaluation

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Dataset struct {
	Train []map[string]interface{} `json:"train"`
}

func LoadDataset(uri string) ([]map[string]interface{}, error) {
	filePath := fmt.Sprintf("LangChainDatasets/%s", uri)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(filePath)
	switch strings.ToLower(ext) {
	case ".json":
		var dataset Dataset
		err = json.Unmarshal(data, &dataset)
		if err != nil {
			return nil, err
		}
		return dataset.Train, nil
	case ".csv":
		return loadCSVData(data), nil
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

func loadCSVData(data []byte) []map[string]interface{} {
	r := csv.NewReader(strings.NewReader(string(data)))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var result []map[string]interface{}
	headers := records[0]
	for _, record := range records[1:] {
		item := make(map[string]interface{})
		for i, value := range record {
			item[headers[i]] = value
		}
		result = append(result, item)
	}

	return result
}
