package documentLoaders

import (
	"context"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/bigquery"
)

type BigQueryLoader struct {
	BaseLoaderImpl
	Query              string
	Project            string
	PageContentColumns []string
	MetadataColumns    []string
}

func NewBigQueryLoader(query string, project string, pageContentColumns []string, metadataColumns []string) *BigQueryLoader {
	return &BigQueryLoader{
		Query:              query,
		Project:            project,
		PageContentColumns: pageContentColumns,
		MetadataColumns:    metadataColumns,
	}
}

func (l *BigQueryLoader) Load() []*documentSchema.Document {
	ctx := context.Background()

	bqClient, err := bigquery.NewClient(ctx, l.Project)
	if err != nil {
		fmt.Println("Error creating BigQuery client:", err)
		return nil
	}

	query := bqClient.Query(l.Query)
	queryResult, err := query.Read(ctx)
	if err != nil {
		fmt.Println("Error executing BigQuery query:", err)
		return nil
	}

	docs := []*documentSchema.Document{}
	pageContentColumns := l.PageContentColumns
	metadataColumns := l.MetadataColumns

	if pageContentColumns == nil {
		pageContentColumns = make([]string, len(queryResult.Schema))
		for i, column := range queryResult.Schema {
			pageContentColumns[i] = column.Name
		}
	}

	if metadataColumns == nil {
		metadataColumns = []string{}
	}

	for {
		var row map[string]bigquery.Value
		err := queryResult.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println("Error iterating over BigQuery results:", err)
			return nil
		}

		pageContent := ""
		for key, value := range row {
			if contains(pageContentColumns, key) {
				pageContent += fmt.Sprintf("%s: %v\n", key, value)
			}
		}

		metadata := map[string]interface{}{}
		for key, value := range row {
			if contains(metadataColumns, key) {
				metadata[key] = value
			}
		}

		doc := &documentSchema.Document{
			PageContent: pageContent,
			Metadata:    metadata,
		}
		docs = append(docs, doc)
	}

	return docs
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
