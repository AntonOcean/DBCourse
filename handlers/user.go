package handlers

import (
	"DBCourse/config"
	"DBCourse/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserCreateHandler(c *gin.Context) {
	db := config.Database()
	defer db.Close()

	nickname := c.Param("nickname")

	var update models.UserUpdateInsert
	var err error

	_ = c.BindJSON(&update)

	err = update.Validate()

	if err != nil || nickname == "" {
		fmt.Println("user UserCreateHandle 1 ",err)
		fmt.Println("un valid data user create")
	}

	user := models.User{
		About: update.About,
		Email: update.Email,
		Fullname: update.Fullname,
		Nickname: nickname,
	}

	err = user.Create(db)

	if err != nil {
		fmt.Println("user UserCreateHandle 2 ",err)
		var duplicate []models.User
		duplicate, err = user.GetDuplicate(db)
		c.JSON(http.StatusConflict, duplicate)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func UserProfileHandler(c *gin.Context) {
	db := config.Database()
	defer db.Close()

	method := c.Request.Method
	nickname := c.Param("nickname")
	var user models.User

	if method == "GET" {
		err := user.Get(db, nickname)
		if err != nil {
			fmt.Println("user UserProfileHandler 1 ",err)
			c.JSON(http.StatusNotFound, models.Error{"Can't find user with nickname: " + nickname})
			return
		}
	} else {
		var update models.UserUpdateInsert
		_ = c.BindJSON(&update)

		err := user.Get(db, nickname)
		if err != nil {
			fmt.Println("user UserProfileHandler 2 ",err)
			c.JSON(http.StatusNotFound, models.Error{"Can't find user with nickname: " + nickname})
			return
		}

		if update.About == "" {
			update.About = user.About
		}

		if update.Email == "" {
			update.Email = user.Email
		}

		if update.Fullname == "" {
			update.Fullname = user.Fullname
		}

		user = models.User{
			About: update.About,
			Email: update.Email,
			Fullname: update.Fullname,
			Nickname: nickname,
		}
		count, err := user.Update(db)
		if err != nil {
			fmt.Println("user UserProfileHandler 3 ",err)
			c.JSON(http.StatusConflict, models.Error{"Can't find user with nickname: " + nickname})
			return
		}
		if count == 0 {
			c.JSON(http.StatusNotFound, models.Error{"Can't find user with nickname: " + nickname})
			return
		}
	}
	c.JSON(http.StatusOK, user)
	return

}