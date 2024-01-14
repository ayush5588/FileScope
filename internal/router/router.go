package router

import (
	"net/http"
	"sort"

	"github.com/ayush5588/FileScope/internal"
	"github.com/ayush5588/FileScope/internal/handler"
	"github.com/ayush5588/FileScope/internal/url"
	"github.com/ayush5588/FileScope/model"
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

	/*
		Method: GET
		Path: /
		Definition: Serves the home page
	*/
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.POST("/getPR", func(c *gin.Context) {

		logger.Info("Serving POST request...")
		userInputURL := c.PostForm("filePath")
		//var userInputURL reqBody

		// err := c.ShouldBindJSON(&userInputURL)
		// if err != nil {
		// 	internal.HandleError(c, err)
		// 	return
		// }

		err := url.ValidateFilePath(&userInputURL)
		if err != nil {
			logger.Errorw("error in validating filepath", "error", err)
			internal.HandleError(c, internal.ErrInvalidURL)
			return
		}

		urlComponent, err := url.ExtractComponentsFromURL(userInputURL)
		if err != nil {
			logger.Errorw("error in extracting details from URL", "error", err)
			internal.HandleError(c, err)
			return
		}

		urlComponent.URL = userInputURL

		prs, err := handler.GetFileModifyingPRs(logger, urlComponent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
			return
		}

		sortByCreatedOn := func(pr1 model.PR, pr2 model.PR) bool {
			if pr1.CreatedOn.Year > pr2.CreatedOn.Year {
				return true
			} else if pr1.CreatedOn.Year < pr2.CreatedOn.Year {
				return false
			} else {
				if pr1.CreatedOn.Month > pr2.CreatedOn.Month {
					return true
				} else if pr1.CreatedOn.Month < pr2.CreatedOn.Month {
					return false
				} else {
					if pr1.CreatedOn.Day > pr2.CreatedOn.Day {
						return true
					} else if pr1.CreatedOn.Day < pr2.CreatedOn.Day {
						return false
					}
				}

			}

			return pr1.Number > pr2.Number
		}

		sort.Slice(prs, func(i, j int) bool {
			return sortByCreatedOn(prs[i], prs[j])
		})

		//c.HTML(http.StatusOK, "index.html", gin.H{"prs": prs})
		c.JSON(http.StatusOK, gin.H{"prs": prs})

		return

	})

	return router
}
