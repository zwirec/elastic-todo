package main

import (
	"context"
	"github.com/olivere/elastic"
	"log"
)

type ES struct {
	client *elastic.Client
}

func (es *ES) CreateTask(t Task) (*elastic.IndexResponse, error) {
	res, err := es.client.Index().Index("todo-list").Type("task").BodyJson(t).Do(context.Background())
	return res, err
}

func (es *ES) CreateIndex() error {
	_, err := es.client.CreateIndex("todo-list").BodyString(mapping).Do(context.Background())
	return err
}

func (es *ES) SearchByID(id string) (*elastic.GetResult, error) {
	return es.client.
		Get().
		Index("todo-list").
		Type("task").
		Id(id).
		Do(context.Background())
}

func (es *ES) SearchByTitle(title string) (*elastic.SearchResult, error) {
	return es.client.
		Search().
		Index("todo-list").
		Query(elastic.NewMatchQuery("title", title)).
		Do(context.Background())
}

func (es *ES) DeleteByID(id string) (*elastic.DeleteResponse, error) {
	return es.client.
		Delete().
		Index("todo-list").
		Id(id).
		Do(context.Background())
}

func (es *ES) UpdateByID(id string, task Task) (*elastic.UpdateResponse, error) {
	return es.client.
		Update().
		Index("todo-list").
		Type("task").
		Id(id).
		Doc(task).
		DocAsUpsert(true).
		FetchSource(true).
		Do(context.Background())
}

func NewES() *ES {
	esClient, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		log.Fatal(err)
	}
	return &ES{client: esClient}
}
