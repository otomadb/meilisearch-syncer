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

	update_tags(conn, ms)
	update_videos(conn, ms)

	fmt.Println("Synced successfully.")
}

func update_tags(conn *pgx.Conn, ms *meilisearch.Client) {
	rows, err := conn.Query(context.Background(), `SELECT "id","name","tagId" FROM "TagName"`)
	defer rows.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	var docs []map[string]interface{}
	for rows.Next() {
		var id string
		var name string
		var tagId string

		err = rows.Scan(&id, &name, &tagId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan the row: %v\n", err)
			os.Exit(1)
		}

		docs = append(docs, map[string]interface{}{"id": id, "name": name, "tag_id": tagId})
	}

	index := ms.Index("tags")

	_, err = index.UpdateIndex("id")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = index.UpdateDistinctAttribute("tag_id")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = index.DeleteAllDocuments()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = index.AddDocuments(docs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Tags sync successfully.")
}

func update_videos(conn *pgx.Conn, ms *meilisearch.Client) {
	rows, err := conn.Query(context.Background(), `SELECT "id","title","videoId" FROM "VideoTitle"`)
	defer rows.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	var docs []map[string]interface{}
	for rows.Next() {
		var id string
		var title string
		var videoId string

		err = rows.Scan(&id, &title, &videoId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan the row: %v\n", err)
			os.Exit(1)
		}

		docs = append(docs, map[string]interface{}{"id": id, "title": title, "video_id": videoId})
	}

	index := ms.Index("videos")

	_, err = index.UpdateIndex("id")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = index.UpdateDistinctAttribute("video_id")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = index.DeleteAllDocuments()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = index.AddDocuments(docs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Videos sync successfully.")
}
