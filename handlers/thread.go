package handlers

import (
	"DBCourse/config"
	"DBCourse/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ThreadCreateHandler(c *gin.Context, slug string) {
	// slug forums
	// select by title

	db := config.Database()
	defer db.Close()

	var thread models.Thread
	var err error

	_ = c.BindJSON(&thread)

	var user models.User
	err = user.Get(db, thread.Author)
	if err != nil {
		fmt.Println("thread ThreadCreateHandler 1 ",err)
		c.JSON(http.StatusNotFound, models.Error{"Can't find user with nickname: " + thread.Author})
		return
	}

	var forum models.Forum
	err = forum.Get(db, slug)
	if err != nil {
		fmt.Println("thread ThreadCreateHandler 2 ",err)
		c.JSON(http.StatusNotFound, models.Error{"Can't find forum with slug: " + slug})
		return
	}

	thread.Author = user.Nickname
	thread.Forum = forum.Slug
	//thread.Slug = slug

	err = thread.Create(db)

	if err != nil {
		fmt.Println("thread ThreadCreateHandler 3 ",err)
		err = thread.Get(db, thread.Slug, thread.Id, thread.Title)
		thread.SetVotesCount(db)
		c.JSON(http.StatusConflict, thread)
		return
	}

	err = thread.Get(db, thread.Slug, thread.Id, thread.Title)

	if err != nil {
		fmt.Println("thread ThreadCreateHandler 4 ",err)
	}

	thread.SetVotesCount(db)
	c.JSON(http.StatusCreated, thread)
}


func DetailsHandler(c *gin.Context) {
	db := config.Database()
	defer db.Close()

	var thread models.Thread
	var err error

	id, err := strconv.Atoi(c.Param("identifier"))
	if err == nil {
		fmt.Println("thread DetailsHandler 1 ",err)
		err = thread.Get(db, "", int32(id), "")
		if err != nil {
			fmt.Println("thread DetailsHandler 2 ",err)
			c.JSON(http.StatusNotFound, models.Error{"Can't find thread with id: " + strconv.Itoa(id)})
			return
		}
	} else {
		slug := c.Param("identifier")
		err = thread.Get(db, slug, 0, "")
		fmt.Println("thread DetailsHandler 3 ",err)
		if err != nil {
			fmt.Println("thread DetailsHandler 4 ",err)
			c.JSON(http.StatusNotFound, models.Error{"Can't find thread with slug: " + slug})
			return
		}
	}

	method := c.Request.Method

	if method == "GET" {
		thread.SetVotesCount(db)
		c.JSON(http.StatusOK, thread)
		return
	} else {
		var threadUpdate models.ThreadUpdate
		_ = c.BindJSON(&threadUpdate)
		if threadUpdate.Title != "" || threadUpdate.Message != "" {
			if threadUpdate.Title != "" {
				thread.Title = threadUpdate.Title
			}
			if threadUpdate.Message != "" {
				thread.Message = threadUpdate.Message
			}
			err = thread.Update(db)
			if err != nil {
				fmt.Println("thread Update 1 ",err)
			}
		}
		c.JSON(http.StatusOK, thread)
		return
	}

}