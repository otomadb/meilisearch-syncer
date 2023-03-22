package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/meilisearch/meilisearch-go"
	"golang.org/x/exp/slog"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		slog.Error("Unable to connect to database.", slog.Any("error", err))
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	ms := meilisearch.NewClient(meilisearch.ClientConfig{Host: os.Getenv("MEILISEARCH_URL")})

	update_tags(conn, ms)
	update_videos(conn, ms)

	slog.Info("Succeedded to sync.")
}

func update_tags(conn *pgx.Conn, ms *meilisearch.Client) {
	rows, err := conn.Query(context.Background(), `SELECT "id","name","tagId" FROM "TagName"`)
	defer rows.Close()
	if err != nil {
		slog.Error("Unable to query tags.", slog.Any("error", err))
		os.Exit(1)
	}

	var docs []map[string]interface{}
	for rows.Next() {
		var id string
		var name string
		var tagId string

		err = rows.Scan(&id, &name, &tagId)
		if err != nil {
			slog.Error("Unable to scan the row.", slog.Any("error", err))
			os.Exit(1)
		}

		docs = append(docs, map[string]interface{}{"id": id, "name": name, "tag_id": tagId})
	}
	slog.Info("Fetched tags.", slog.Int("count", len(docs)))

	index := ms.Index("tags")

	task, err := index.UpdateIndex("id")
	if err != nil {
		slog.Error("Unable to set primary key in tags index.", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Set primary key in tags index", slog.Int64("task_id", task.TaskUID))

	task, err = index.UpdateDistinctAttribute("tag_id")
	if err != nil {
		slog.Error("Unable to set distinct attribute in tags index.", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Update distinct attribute in tags index", slog.Int64("task_id", task.TaskUID))

	/*
		_, err = index.DeleteAllDocuments()
		if err != nil {
			slog.Error("Unable to delete all documents in tag index", err)
			os.Exit(1)
		}
	*/

	task, err = index.AddDocuments(docs)
	if err != nil {
		slog.Error("Unable to add documents in tags index.", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Add tags documents", slog.Int64("task_id", task.TaskUID))

	slog.Info("Succeedded to update tags.")
}

func update_videos(conn *pgx.Conn, ms *meilisearch.Client) {
	rows, err := conn.Query(context.Background(), `SELECT "id","title","videoId" FROM "VideoTitle"`)
	defer rows.Close()
	if err != nil {
		slog.Error("Unable to query tags.", slog.Any("error", err))
		os.Exit(1)
	}

	var docs []map[string]interface{}
	for rows.Next() {
		var id string
		var title string
		var videoId string

		err = rows.Scan(&id, &title, &videoId)
		if err != nil {
			slog.Error("Unable to scan the row.", slog.Any("error", err))
			os.Exit(1)
		}

		docs = append(docs, map[string]interface{}{"id": id, "title": title, "video_id": videoId})
	}
	slog.Info("Fetched videos.", slog.Int("count", len(docs)))

	index := ms.Index("videos")

	task, err := index.UpdateIndex("id")
	if err != nil {
		slog.Error("Unable to set primary key in videos index.", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Set primary key in video index", slog.Int64("task_id", task.TaskUID))

	task, err = index.UpdateDistinctAttribute("video_id")
	if err != nil {
		slog.Error("Unable to set distinct attribute in videos index.", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Update distinct attribute in video index", slog.Int64("task_id", task.TaskUID))

	/*
		_, err = index.DeleteAllDocuments()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	*/

	task, err = index.AddDocuments(docs)
	if err != nil {
		slog.Error("Unable to add documents in videos index.", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Add videos documents", slog.Int64("task_id", task.TaskUID))

	slog.Info("Succeedded to update videos.")
}
