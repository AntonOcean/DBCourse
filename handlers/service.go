package handlers

import (
	"DBCourse/config"
	"DBCourse/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ServiceStatusHandler(c *gin.Context) {
	db := config.Database()
	defer db.Close()

	var info models.InfoDB

	info.Get(db)

	c.JSON(http.StatusOK, info)

}

func ServiceClearHandler(c *gin.Context) {
	db := config.Database()
	defer db.Close()

	models.DropAll(db)

	c.Status(http.StatusOK)
}