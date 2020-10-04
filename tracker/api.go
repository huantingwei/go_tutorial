package tracker

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const layoutISO = "2006-01-02 15:04:05"

func ListBook(c *gin.Context) {
	id := c.Query("id")
	title := c.Query("title")
	author := c.Query("author")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	filter := map[string]string{
		"id":        id,
		"title":     title,
		"author":    author,
		"startTime": startTime,
		"endTime":   endTime,
	}
	books, err := listBook(filter)
	if err != nil {
		ResponseBadRequest(c, err)
	} else {
		ResponseSuccess(c, books)
	}
}

func GetBook(c *gin.Context) {
	id := c.Param("bookid")
	book, err := getBook(id)
	if err != nil {
		ResponseBadRequest(c, err)
	} else {
		ResponseSuccess(c, book)
	}
}

func AddBook(c *gin.Context) {
	title := c.PostForm("title")
	author := c.PostForm("author")
	status, _ := strconv.Atoi(c.PostForm("status"))
	startTime, _ := time.Parse(layoutISO, c.PostForm("startTime"))
	endTime, _ := time.Parse(layoutISO, c.PostForm("endTime"))
	description := c.PostForm("description")
	book := Book{
		Title:       title,
		Author:      author,
		Status:      status,
		StartTime:   startTime,
		EndTime:     endTime,
		Description: description,
	}
	oid, err := addBook(&book)
	if err != nil {
		ResponseBadRequest(c, err)
	}
	ResponseSuccess(c, oid)
}

func DeleteBook(c *gin.Context) {
	id := c.PostFormArray("id")
	deleteCount, err := deleteBook(id)
	if err != nil {
		ResponseBadRequest(c, err)
	} else {
		ResponseSuccess(c, deleteCount)
	}
}

func EditBook(c *gin.Context) {
	fields := make(map[string]interface{})
	fields["id"], _ = primitive.ObjectIDFromHex(c.Param("bookid"))
	fields["title"] = c.PostForm("title")
	fields["author"] = c.PostForm("author")
	fields["status"], _ = strconv.Atoi(c.PostForm("status"))
	fields["startTime"], _ = time.Parse(layoutISO, c.PostForm("startTime"))
	fields["endTime"], _ = time.Parse(layoutISO, c.PostForm("endTime"))
	fields["description"] = c.PostForm("description")

	editCount, err := editBook(fields)
	if err != nil {
		ResponseBadRequest(c, err)
	} else {
		ResponseSuccess(c, editCount)
	}
}

// func ListNote {}
// func GetNote{}
// func AddNote{}
// func DeleteNote{}
// func EditNote{}