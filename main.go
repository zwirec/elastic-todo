package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
)

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"task":{
			"properties":{
				"title":{
					"type":"keyword"
				},
				"description":{
					"type":"text",
				},
				"deadline":{
					"type":"date"
				},
				"created":{
					"type":"date"
				}
			}
		}
	}
}`

func main() {
	esc := NewES()
	esc.CreateIndex()
	r := gin.Default()
	g := r.Group(`/`, func(c *gin.Context) {
		c.Set(elasticServiceKey, esc)
	})
	g.POST(`/task`, CreateTask)
	g.GET(`/get/by_id`, GetByID)
	g.GET(`/get/by_title`, GetByTitle)
	g.DELETE(`/task`, DeleteByID)
	g.PUT(`/task`, UpdateByID)
	r.Run(":8000")
}

func UpdateByID(c *gin.Context) {
	esc := c.MustGet(elasticServiceKey).(*ES)
	id, ok := c.GetQuery("id")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	task := Task{}
	if err := c.BindJSON(&task); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if res, err := esc.UpdateByID(id, task); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
		return
	} else {
		c.JSON(http.StatusOK, res.GetResult.Source)
	}
}

const elasticServiceKey = "elastic"

func CreateTask(c *gin.Context) {
	esc := c.MustGet(elasticServiceKey).(*ES)
	t := Task{}
	err := c.BindJSON(&t)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}
	if _, err := esc.CreateTask(t); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
		return
	} else {
		c.Status(http.StatusCreated)
		return
	}
}

func GetByID(c *gin.Context) {
	esc := c.MustGet(elasticServiceKey).(*ES)
	id := c.Query("id")
	res, err := esc.SearchByID(id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
	}
	c.JSON(http.StatusOK, res.Source)
}

func GetByTitle(c *gin.Context) {
	esc := c.MustGet(elasticServiceKey).(*ES)
	title, ok := c.GetQuery("title")
	log.Println(title)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if res, err := esc.SearchByTitle(title); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
	} else {
		tasks := make([]Task, 0)
		task := Task{}
		for _, item := range res.Each(reflect.TypeOf(task)) {
			task := item.(Task)
			tasks = append(tasks, task)
		}
		c.JSON(http.StatusOK, tasks)
	}
}

func DeleteByID(c *gin.Context) {
	esc := c.MustGet(elasticServiceKey).(*ES)
	id, ok := c.GetQuery("id")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	_, err := esc.DeleteByID(id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.Status(http.StatusOK)
	return
}
