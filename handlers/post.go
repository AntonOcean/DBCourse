package handlers

import (
	"DBCourse/config"
	"DBCourse/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func PostsCreateHandler(c *gin.Context) {
	// создание bulk
	//loc, _ := time.LoadLocation("Europe/Moscow")
	str := time.Now()
	currentTime := time.Date(str.Year(), str.Month(), str.Day(), str.Hour(), str.Minute(), 0, 0, str.Location())
	//currentTime := time.Now().Format("2006-01-02 15:04:05 +0300 MSK")

	db := config.Database()
	defer db.Close()

	var thread models.Thread
	var err error

	var posts []models.Post

	err = c.BindJSON(&posts)


	id, err := strconv.Atoi(c.Param("identifier"))
	//fmt.Println(err)
	if err == nil {
		fmt.Println("post PostsCreateHandler 1 ",err)
		err = thread.Get(db, "", int32(id), "")
		if err != nil {
			fmt.Println("post PostsCreateHandler 2 ",err)
			c.JSON(http.StatusNotFound, models.Error{"Can't find post thread with id: " + strconv.Itoa(id)})
			return
		}
	} else {
		slug := c.Param("identifier")
		err = thread.Get(db, slug, 0, "")
		if err != nil {
			fmt.Println("post PostsCreateHandler 4 ",err)
			c.JSON(http.StatusNotFound, models.Error{"Can't find post thread with slug: " + slug})
			return
		}
	}

	if len(posts) ==0 {

		c.JSON(http.StatusCreated, []int{})
		return
	}

	var forum models.Forum

	err = forum.Get(db, thread.Forum)

	if err != nil {

		fmt.Println("post PostsCreateHandler 5 ",err)
	}

	var postsCreated []interface{}

	//
	for _, post := range posts {
		post.Thread = thread.Id
		post.Forum = forum.Slug

		var user models.User
		err := user.Get(db, post.Author)
		if err != nil {
			fmt.Println("user UserProfileHandler 2 ",err)
			c.JSON(http.StatusNotFound, models.Error{"Can't find author with nickname: " + post.Author})
			return
		}

		if post.Parent != 0 {
			exist, err := post.CheckParentExists(db)
			if err != nil || !exist {
				fmt.Println("post PostsCreateHandler 6 ",err)
				c.JSON(http.StatusConflict, models.Error{"Can't find post with id: " + strconv.Itoa(int(post.Parent))})
				return
			}
		}

		err = post.CreatePost(db, currentTime)
		if err != nil {
			fmt.Println("post PostsCreateHandler 7 ",err)
		}
		err = post.Get(db)
		if err != nil {
			fmt.Println("post PostsCreateHandler 8 ",err)
		}

		postsCreated = append(postsCreated, post)

	}

	if posts == nil {
		c.JSON(http.StatusCreated, []int{})
		return
	}

	c.JSON(http.StatusCreated, postsCreated)

}

func PostDetailsHandler(c *gin.Context) {
	db := config.Database()
	defer db.Close()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		fmt.Println("post PostDetailsHandler 1", err)
	}

	var post models.Post

	err = post.GetById(db, id)

	if err != nil {
		fmt.Println("post PostDetailsHandler 2", err)
		c.JSON(http.StatusNotFound, models.Error{"Can't find post with id: " + strconv.Itoa(id)})
		return
	}

	method := c.Request.Method

	if method == "GET" {

		queryString := c.Request.URL.Query()

		//"user" "forum" "thread"
		relateds := strings.Split(queryString.Get("related"), ",")

		result := map[string]interface{}{}
		result["post"] = post

		for _, obj := range relateds {
			switch obj {
			case "user":
				var user models.User
				err = user.Get(db, post.Author)
				if err != nil {
					fmt.Println("blat post 1")
				}
				result["author"] = user
			case "forum":
				var forum models.Forum
				err = forum.Get(db, post.Forum)
				if err != nil {
					fmt.Println("blat post 2")
				}
				forum.SetThreadsCount(db)
				forum.SetPostsCount(db)
				result["forum"] = forum
			case "thread":
				var thread models.Thread
				err = thread.Get(db, "", post.Thread, "")
				if err != nil {
					fmt.Println("blat post 3")
				}
				thread.SetVotesCount(db)
				result["thread"] = thread
			}
		}
		c.JSON(http.StatusOK, result)
		return

	} else {

		var postUpdate models.PostUpdate

		_ = c.BindJSON(&postUpdate)

		if postUpdate.Message != "" && post.Message != postUpdate.Message {
			post.Message = postUpdate.Message
			err = post.Update(db)
			post.IsEdited = true
			if err != nil {
				fmt.Println("post PostDetailsHandler 3")
			}
		}
		//post.IsEdited = true
		c.JSON(http.StatusOK, post)
	}

}


func GetPostsHandler(c *gin.Context) {
	db := config.Database()
	defer db.Close()

	var thread models.Thread
	var err error

	id, err := strconv.Atoi(c.Param("identifier"))
	if err == nil {
		fmt.Println("post GetPostsHandler 1 ",err)
		err = thread.Get(db, "", int32(id), "")
		if err != nil {
			fmt.Println("post GetPostsHandler 2 ",err)
			c.JSON(http.StatusNotFound, models.Error{"Can't find thread with id: " + strconv.Itoa(id)})
			return
		}
	} else {
		slug := c.Param("identifier")
		err = thread.Get(db, slug, 0, "")
		//fmt.Println("post GetPostsHandler 3 ",err)
		if err != nil {
			fmt.Println("post GetPostsHandler 4 ",err)
			c.JSON(http.StatusNotFound, models.Error{"Can't find thread with slug: " + slug})
			return
		}
	}


	queryString := c.Request.URL.Query()

	limit, err := strconv.Atoi(queryString.Get("limit"))
	if err != nil {
		fmt.Println("post GetPostsHandler 5 ",err)
		limit = 501
	}

	since, err := strconv.Atoi(queryString.Get("since"))
	if err != nil {
		fmt.Println("post GetPostsHandler 6 ",err)
		since = 0
	}

	desc := queryString.Get("desc")

	if desc == "" {
		desc = "false"
	}

	sort := queryString.Get("sort")

	posts, err := thread.GetListPost(db, int64(since), sort, desc, limit)

	if err != nil {
		fmt.Println("post GetPostsHandler 6 ",err)
	}

	if len(posts) == 0 {
		c.JSON(http.StatusOK, []int{})
		return
	}

	c.JSON(http.StatusOK, posts)

}
