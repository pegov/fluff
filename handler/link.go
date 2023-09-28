package handler

import (
	"encoding/json"
	"io"
	"log"

	"github.com/pegov/fluff/db"
	"github.com/pegov/fluff/model"

	"github.com/gin-gonic/gin"
)

func CreateLink(repo *db.LinkRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		buffer, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"detail": "Server error"})
			return
		}

		payload := model.CreateLinkRequest{}
		err = json.Unmarshal(buffer, &payload)
		if err != nil {
			log.Println(err)
			c.JSON(422, gin.H{"detail": "Unprocessable entity"})
			return
		}

		err = payload.Validate()
		if err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"detail": err.Error()})
			return
		}

		short, err := repo.Create(payload.Long, 8)

		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"detail": "Server error: can't create link"})
			return
		}

		c.JSON(200, model.CreateLinkResponse{
			Short: short,
		})
	}
}

func GetAllLinks(repo *db.LinkRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		links, err := repo.GetAllLinks()
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(500)
		}

		c.JSON(200, links)
	}
}

func GetLink(repo *db.LinkRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		short := c.Param("short")

		link, err := repo.GetByShort(short)
		if link == nil && err == nil {
			c.AbortWithStatus(404)
		} else if err != nil {
			log.Println(err)
			c.AbortWithError(500, err)
		}

		c.JSON(200, link)
	}
}

func DeleteLink(repo *db.LinkRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		short := c.Param("short")

		link, err := repo.GetByShort(short)
		if link == nil && err == nil {
			c.AbortWithStatus(404)
		}

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(500)
		}

		repo.DeleteByShort(short)

		c.Status(200)
	}
}

func RedirectToLink(repo *db.LinkRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		short := c.Param("short")
		link, err := repo.GetByShort(short)
		if link == nil && err == nil {
			c.AbortWithStatus(404)
		} else if err != nil {
			log.Println(err)
			c.AbortWithError(500, err)
		}

		c.Redirect(301, link.Long)
	}
}
