package router

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ayush5588/FileScope/internal"
	"github.com/ayush5588/FileScope/internal/handler"
	"github.com/ayush5588/FileScope/internal/url"
	"github.com/ayush5588/FileScope/model"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
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
			2.3 If <50, assign the next token value to the GITHUB_TOKEN env var
			2.4 Repeat the step 2.3 until the valid token found
	*/

	myenv := make(map[string]string)

	id := 1
	for {
		tokenName := fmt.Sprintf("GITHUB_TOKEN_%d", id)
		tokenVal := os.Getenv(tokenName)
		if tokenVal == "" {
			break
		}
		myenv[tokenName] = tokenVal
		id += 1
	}

	var githubToken string

	// Iterate over the env variables of pattern GITHUB_TOKEN_<%d>
	id = 1
	for {
		tokenName := fmt.Sprintf("GITHUB_TOKEN_%d", id)
		tokenVal := os.Getenv(tokenName)
		if tokenVal == "" {
			break
		}
		if githubToken == "" {
			githubToken = tokenVal
		}
		myenv[tokenName] = tokenVal
		id += 1
	}
	logger.Info("Making api call for the current token")
	// Make a API call to get the request lefts for the current token
	reqLeft, err := makeRateLimitAPICall(githubToken)
	if err != nil {
		logger.Errorw("error in making rate-limit api acall for current token", "err", err)
		return err
	}

	if reqLeft > 50 {
		os.Setenv("GH_TOKEN", githubToken)
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

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func comparePR(pr1, pr2 model.PR) int {
	d1 := pr1.CreatedOn.FullDate
	d2 := pr2.CreatedOn.FullDate

	date1, _ := time.Parse("2/1/2006", d1)
	date2, _ := time.Parse("2/1/2006", d2)

	if !date1.Equal(date2) {
		if date1.After(date2) {
			return -1
		}
		return 1
	}

	return cmp.Compare(pr2.Number, pr1.Number)
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

	router.Use(corsMiddleware())

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
		token := c.PostForm("token")
		// if user has provided their GitHub Token, then we will not set / read the GitHub Token from env value.
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

		prs, err := handler.GetFileModifyingPRs(logger, urlComponent, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
			return
		}

		slices.SortFunc(prs, comparePR)
		fmt.Println(prs[0])
		//c.HTML(http.StatusOK, "index.html", gin.H{"prs": prs})
		c.JSON(http.StatusOK, gin.H{"prs": prs})

		return

	})

	return router
}
