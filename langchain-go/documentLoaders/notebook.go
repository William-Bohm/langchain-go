package documentLoaders

import (
	"encoding/json"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"io/ioutil"
	"os"
	"strings"
)

type NotebookLoader struct {
	FilePath        string
	IncludeOutputs  bool
	MaxOutputLength int
	RemoveNewline   bool
	Traceback       bool
}

type Cell struct {
	CellType string   `json:"cell_type"`
	Source   []string `json:"source"`
	Outputs  []Output `json:"outputs"`
}

type Output struct {
	Ename      string   `json:"ename"`
	Evalue     string   `json:"evalue"`
	Traceback  []string `json:"traceback"`
	OutputType string   `json:"output_type"`
	Text       []string `json:"text"`
}

func concatenateCells(cell Cell, includeOutputs bool, maxOutputLength int, traceback bool) string {
	if includeOutputs && cell.CellType == "code" && len(cell.Outputs) > 0 {
		if cell.Outputs[0].Ename != "" {
			if traceback {
				return "'" + cell.CellType + "' cell: '" + strings.Join(cell.Source, "") + "', gives error '" + cell.Outputs[0].Ename + "', with description '" + cell.Outputs[0].Evalue + "' and traceback '" + strings.Join(cell.Outputs[0].Traceback, "") + "'\n\n"
			} else {
				return "'" + cell.CellType + "' cell: '" + strings.Join(cell.Source, "") + "', gives error '" + cell.Outputs[0].Ename + "', with description '" + cell.Outputs[0].Evalue + "'\n\n"
			}
		} else if cell.Outputs[0].OutputType == "stream" {
			output := strings.Join(cell.Outputs[0].Text, "")
			minOutput := len(output)
			if minOutput > maxOutputLength {
				minOutput = maxOutputLength
			}
			return "'" + cell.CellType + "' cell: '" + strings.Join(cell.Source, "") + "' with output: '" + output[:minOutput] + "'\n\n"
		}
	} else {
		return "'" + cell.CellType + "' cell: '" + strings.Join(cell.Source, "") + "'\n\n"
	}

	return ""
}

func removeNewlines(x interface{}) interface{} {
	switch v := x.(type) {
	case string:
		return strings.ReplaceAll(v, "\n", "")
	case []string:
		newArr := make([]string, len(v))
		for i, s := range v {
			newArr[i] = removeNewlines(s).(string)
		}
		return newArr
	case []interface{}:
		newArr := make([]interface{}, len(v))
		for i, s := range v {
			newArr[i] = removeNewlines(s)
		}
		return newArr
	default:
		return v
	}
}

func (n *NotebookLoader) Load() []documentSchema.Document {
	jsonFile, _ := os.Open(n.FilePath)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var notebook map[string]interface{}
	json.Unmarshal(byteValue, &notebook)

	cells := notebook["cells"].([]interface{})

	text := ""

	for _, c := range cells {
		cellMap := c.(map[string]interface{})
		cell := Cell{
			CellType: cellMap["cell_type"].(string),
			Source:   cellMap["source"].([]string),
			Outputs:  []Output{},
		}

		for _, o := range cellMap["outputs"].([]interface{}) {
			outputMap := o.(map[string]interface{})
			output := Output{
				Ename:      outputMap["ename"].(string),
				Evalue:     outputMap["evalue"].(string),
				Traceback:  outputMap["traceback"].([]string),
				OutputType: outputMap["output_type"].(string),
				Text:       outputMap["text"].([]string),
			}
			cell.Outputs = append(cell.Outputs, output)
		}

		if n.RemoveNewline {
			cell.CellType = removeNewlines(cell.CellType).(string)
			cell.Source = removeNewlines(cell.Source).([]string)
			for i, _ := range cell.Outputs {
				cell.Outputs[i].Ename = removeNewlines(cell.Outputs[i].Ename).(string)
				cell.Outputs[i].Evalue = removeNewlines(cell.Outputs[i].Evalue).(string)
				cell.Outputs[i].Traceback = removeNewlines(cell.Outputs[i].Traceback).([]string)
				cell.Outputs[i].OutputType = removeNewlines(cell.Outputs[i].OutputType).(string)
				cell.Outputs[i].Text = removeNewlines(cell.Outputs[i].Text).([]string)
			}
		}

		text += concatenateCells(cell, n.IncludeOutputs, n.MaxOutputLength, n.Traceback)
	}

	metadata := map[string]interface{}{
		"source": n.FilePath,
	}

	return []documentSchema.Document{
		{
			PageContent: text,
			Metadata:    metadata,
		},
	}
}
