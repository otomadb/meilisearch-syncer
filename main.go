package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/meilisearch/meilisearch-go"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), `SELECT array_agg("name") AS names, "tagId" FROM "TagName" GROUP BY "tagId"`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	var tagDocuments []map[string]interface{}
	defer rows.Close()
	for rows.Next() {
		var names []string
		var tagId string

		err = rows.Scan(&names, &tagId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan the row: %v\n", err)
			os.Exit(1)
		}

		tagDocuments = append(tagDocuments, map[string]interface{}{"id": tagId, "names": names})
	}

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host: os.Getenv("MEILISEARCH_URL"),
	})
	index := client.Index("tags")

	task, err := index.AddDocuments(tagDocuments)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(task.TaskUID)
}
