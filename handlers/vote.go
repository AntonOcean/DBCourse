package handlers

import (
	"DBCourse/config"
	"DBCourse/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)


func VoteCreateHandler(c *gin.Context) {
	db := config.Database()
	defer db.Close()

	var thread models.Thread
	var err error

	id, err := strconv.Atoi(c.Param("identifier"))
	if err == nil {
		fmt.Println("vote VoteCreateHandler 1 ",err)
		err = thread.Get(db, "", int32(id), "")
		if err != nil {
			fmt.Println("vote VoteCreateHandler 2 ",err)
			c.JSON(http.StatusNotFound, models.Error{"Can't find thread with id: " + strconv.Itoa(id)})
			return
		}
	} else {
		slug := c.Param("identifier")
		err = thread.Get(db, slug, 0, "")
		fmt.Println("vote VoteCreateHandler 3 ",err)
		if err != nil {
			fmt.Println("vote VoteCreateHandler 4 ",err)
			c.JSON(http.StatusNotFound, models.Error{"Can't find thread with slug: " + slug})
			return
		}
	}

	var vote models.Vote

	_ = c.BindJSON(&vote)

	var user models.User

	err = user.Get(db, vote.Nickname)
	if err != nil {
		fmt.Println("vote VoteCreateHandler 5 ",err)
		c.JSON(http.StatusNotFound, models.Error{"Can't find user with nickname: " + vote.Nickname})
		return
	}


	// TODO не опитимально
	err = vote.Create(db, thread.Id)

	if err != nil {
		fmt.Println("vote VoteCreateHandler 6 ",err)
	}

	//err =thread.Get(db, thread.Slug, thread.Id, thread.Title)
	thread.SetVotesCount(db)

	c.JSON(http.StatusOK, thread)

}