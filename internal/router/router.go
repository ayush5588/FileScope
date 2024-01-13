package router

import (
	"net/http"

	"github.com/ayush5588/FileScope/internal"
	"github.com/ayush5588/FileScope/internal/handler"
	"github.com/ayush5588/FileScope/internal/url"
	"github.com/gin-gonic/gin"
)

type reqBody struct {
	URL string `json:"url"`
}

// SetupRouter initalizes the router
func SetupRouter() *gin.Engine {

	logger := internal.GetLogger()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/templates", "./templates/")

	/*
		Method: GET
		Path: /healthz
		Definition: Represents server health
	*/
	router.GET("/healthz", func(c *gin.Context) {
		logger.Info("Successfully served GET /healthz request")
		c.JSON(http.StatusOK, gin.H{"message": "Server is healthy"})
		return
	})

	router.POST("/getPR", func(c *gin.Context) {

		logger.Info("Serving POST request...")
		var rb reqBody
		err := c.ShouldBindJSON(&rb)

		userInputURL := &rb.URL

		err = url.ValidateFilePath(userInputURL)
		if err != nil {
			logger.Errorw("error in validating filepath", "error", err)
			internal.HandleError(c, internal.ErrInvalidURL)
			return
		}

		urlComponent, err := url.ExtractComponentsFromURL(*userInputURL)
		if err != nil {
			logger.Errorw("error in extracting details from URL", "error", err)
			internal.HandleError(c, err)
			return
		}

		urlComponent.URL = *userInputURL

		prs, err := handler.GetFileModifyingPRs(urlComponent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": prs})
		return

	})

	return router
}
