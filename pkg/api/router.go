package api

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ConfigureRoutes(r gin.IRouter) {
	r.PUT("/:owner/:repository/:branch/:path", func(c *gin.Context) {
		owner, repository, branch, path := getParams(c)
		var o PutRequestOptions
		if err := c.BindJSON(&o); err != nil {
			log.Printf("failed to handle JSON: %v", err)
			c.AbortWithStatus(400)
			return
		}

		err := Commit(c.Request.Context(), o.Token, owner, repository, branch, path, o.Message, o.File, o.CreatePR, o.MergePR)
		if err != nil {
			log.Printf("failed to commit: %v", err)
			c.JSON(400, gin.H{
				"ok":     false,
				"errors": err,
			})
			return
		}

		c.JSON(200, gin.H{
			"ok":     true,
			"errors": err,
		})

	})

	r.PATCH("/:owner/:repository/:branch/:path", func(c *gin.Context) {
		owner, repository, branch, path := getParams(c)
		var o PatchRequestOptions
		if err := c.BindJSON(&o); err != nil {
			log.Printf("failed to handle JSON: %v", err)
			c.AbortWithStatus(400)
			return
		}

		err := Patch(c.Request.Context(), o.Token, owner, repository, branch, path, o.Message, o.Patch, o.CreatePR, o.MergePR)
		if err != nil {
			log.Printf("failed to patch: %v", err)
			c.JSON(400, gin.H{
				"ok":     false,
				"errors": err,
			})
			return
		}

		c.JSON(200, gin.H{
			"ok":     true,
			"errors": err,
		})
	})
}

func getParams(c *gin.Context) (owner, repository, branch, path string) {
	owner = c.Param("owner")
	repository = c.Param("repository")
	branch = c.Param("branch")
	path = c.Param("path")
	return
}
