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

	ms := meilisearch.NewClient(meilisearch.ClientConfig{Host: os.Getenv("MEILISEARCH_URL")})

	rows, err := conn.Query(context.Background(), `SELECT "id","name","tagId" FROM "TagName"`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	var tagDocuments []map[string]interface{}
	defer rows.Close()
	for rows.Next() {
		var id string
		var name string
		var tagId string

		err = rows.Scan(&id, &name, &tagId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan the row: %v\n", err)
			os.Exit(1)
		}

		tagDocuments = append(tagDocuments, map[string]interface{}{"id": id, "name": name, "tag_id": tagId})
	}

	tags_index := ms.Index("tags")

	_, err = tags_index.DeleteAllDocuments()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = tags_index.UpdateDistinctAttribute("tag_id")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = tags_index.AddDocuments(tagDocuments)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Synced successfully.")
}
