package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/ayush5588/FileScope/internal"
	"github.com/ayush5588/FileScope/internal/handler"
	"github.com/ayush5588/FileScope/internal/url"
	"github.com/ayush5588/FileScope/model"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type reqBody struct {
	URL string `json:"url"`
}

func makeRateLimitAPICall(token string) (int, error) {

	url := "https://api.github.com/rate_limit"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	bearerToken := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Authorization", bearerToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	var resBody model.RateLimitResBody

	err = json.Unmarshal(body, &resBody)
	if err != nil {
		return 0, err
	}

	return resBody.Rate.Remaining, nil

}

func manageToken(logger *zap.SugaredLogger) error {
	logger.Info("inside manageToken ...")
	/*
		1. Run a cronjob every 5 minutes
		2. In each iteration do the following:
			2.1 Get the current GITHUB_TOKEN env var value
			2.2 Make a rate limit call to check number of reqs left
			2.3 If <20, assign the next token value to the GITHUB_TOKEN env var
			2.4 Repeat the step 2.3 until the valid token found
	*/
	myenv, err := godotenv.Read("token.env")
	if err != nil {
		logger.Errorw("error in reading token.env file", "err", err)
		return err
	}

	var GitHubToken string
	if token, ok := myenv["GH_TOKEN"]; ok {
		GitHubToken = token
	}

	if GitHubToken == "" {
		GitHubToken = myenv["GH_TOKEN_1"]
	}
	logger.Info("Making api call for the current token")
	// Make a API call to get the request lefts for the current token
	reqLeft, err := makeRateLimitAPICall(GitHubToken)
	if err != nil {
		logger.Errorw("error in making rate-limit api acall for current token", "err", err)
		return err
	}

	if reqLeft > 50 {
		os.Setenv("GH_TOKEN", GitHubToken)
		return nil
	}

	for tokenName, token := range myenv {
		// Make an API call to get the request limits left for the token
		logger.Info("Making api call for the each token")
		reqLeft, err := makeRateLimitAPICall(token)
		if err != nil {
			logger.Errorf("error in making rate-limit api call for token %s", tokenName, "err", err)
			return err
		}
		if reqLeft > 50 {
			os.Setenv("GH_TOKEN", token)
			logger.Infof("Current token: %s", tokenName)
			return nil
		}
	}

	os.Setenv("GH_TOKEN", "")

	return internal.ErrNoValidToken
}

// SetupRouter initalizes the router
func SetupRouter() *gin.Engine {

	logger := internal.GetLogger()

	go func(logger *zap.SugaredLogger) {
		s := gocron.NewScheduler(time.UTC)
		s.Every(5).Minutes().Do(func() {
			logger.Info("Inside cronjob")
			manageToken(logger)
		})
		s.StartBlocking()
	}(logger)

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
