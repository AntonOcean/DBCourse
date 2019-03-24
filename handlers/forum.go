package handlers

import (
	"DBCourse/config"
	"DBCourse/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ForumCreateHandler(c *gin.Context) {
	db := config.Database()
	defer db.Close()

	var forum models.Forum
	var err error

	_ = c.BindJSON(&forum)

	err = forum.Validate()

	if err != nil {
		fmt.Println("forum ForumCreateHandler 1 ",err)
		fmt.Println("un valid data forum create")
	}

	var user models.User
	err = user.Get(db, forum.User)
	if err != nil {
		fmt.Println("forum ForumCreateHandler 2 ",err)
		c.JSON(http.StatusNotFound, models.Error{"Can't find user with nickname: " + forum.User})
		return
	}

	forum.User = user.Nickname

	err = forum.Create(db)

	if err != nil {
		fmt.Println("forum ForumCreateHandler 3 ",err)
		err = forum.Get(db, forum.Slug)
		c.JSON(http.StatusConflict, forum)
		return
	}
	c.JSON(http.StatusCreated, forum)
}

func ForumDetailHandler(c *gin.Context) {
	db := config.Database()
	defer db.Close()

	slug := c.Param("slug")

	var forum models.Forum
	var err error

	err = forum.Get(db, slug)
	if err != nil {
		fmt.Println("forum ForumDetailHandler 1 ",err)
		c.JSON(http.StatusNotFound, models.Error{"Can't find forum with slug: " + slug})
		return
	}
	forum.SetPostsCount(db)
	forum.SetThreadsCount(db)

	c.JSON(http.StatusOK, forum)
}

func ForumThreadList(c *gin.Context) {
	db := config.Database()
	defer db.Close()


	slug := c.Param("slug")
	queryString := c.Request.URL.Query()
	limit, err := strconv.Atoi(queryString.Get("limit"))
	if err != nil {
		fmt.Println("forum ForumThreadList 1 ",err)
		limit = 501
	}

	since := queryString.Get("since")
	//if since == "" {
	//	//since = "1900-01-01T00:00:00.000Z"
	//	since = "9999-12-11 23:59:59.997"
	//}

	desc := queryString.Get("desc")

	var forum models.Forum
	//var err error

	err = forum.Get(db, slug)
	if err != nil {
		fmt.Println("forum ForumThreadList 2 ",err)
		c.JSON(http.StatusNotFound, models.Error{"Can't find forum with slug: " + slug})
		return
	}
	forum.SetPostsCount(db)
	forum.SetThreadsCount(db)

	threads, _ := forum.GetThreadList(db, limit, since, desc)

	if threads == nil {
		c.JSON(http.StatusOK, []int{})
		return
	}

	c.JSON(http.StatusOK, threads)
}


func ForumUserList(c *gin.Context) {
	db := config.Database()
	defer db.Close()


	slug := c.Param("slug")
	queryString := c.Request.URL.Query()

	limit, err := strconv.Atoi(queryString.Get("limit"))
	if err != nil {
		fmt.Println("forum ForumThreadList 1 ",err)
		limit = 501
	}

	since := queryString.Get("since")

	desc := queryString.Get("desc")

	if desc == "" {
		desc = "false"
	}

	var forum models.Forum

	err = forum.Get(db, slug)
	if err != nil {
		fmt.Println("forum ForumThreadList 2 ",err)
		c.JSON(http.StatusNotFound, models.Error{"Can't find forum with slug: " + slug})
		return
	}

	users, _ := forum.GetUserList(db, limit, since, desc)

	if users == nil {
		c.JSON(http.StatusOK, []int{})
		return
	}

	c.JSON(http.StatusOK, users)
}