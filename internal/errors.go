package internal

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// ErrOwnerNotFound ...
	ErrOwnerNotFound = errors.New("owner not found in the given path url")
	// ErrRepoNotFound ...
	ErrRepoNotFound = errors.New("repo not found in the given path url")
	// ErrFilePathCannotBeDetected ...
	ErrFilePathCannotBeDetected = errors.New("file path not found in the given path url")
	// ErrInvalidURL ...
	ErrInvalidURL = errors.New("invalid url")
	// ErrNoValidToken ...
	ErrNoValidToken = errors.New("rate limit hit for all the tokens. Wait for an hour")
)

// HandleError ...
func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrOwnerNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Owner cannot be detected in the given URL."})
		break
	case errors.Is(err, ErrRepoNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Repository cannot be detected in the given URL."})
		break
	case errors.Is(err, ErrFilePathCannotBeDetected):
		c.JSON(http.StatusBadRequest, gin.H{"msg": "File path cannot be detected in the given URL."})
		break
	case errors.Is(err, ErrInvalidURL):
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid URL. Please provide a valid URL."})
		break
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Please try again after some time"})
		break
	}
	return
}
