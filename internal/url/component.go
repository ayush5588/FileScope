package url

import (
	"fmt"
	"regexp"

	"github.com/ayush5588/FileScope/internal"
	"github.com/ayush5588/FileScope/model"
)

// ExtractComponentsFromURL will extract OWNER, REPONAME & Relative FILEPATH from the valid URL
func ExtractComponentsFromURL(url string) (model.FileInfo, error) {
	var filePath model.FileInfo

	re := regexp.MustCompile(`https://github.com/([^/]+)/([^/]+)/blob/([^/]+)/(.+)`)

	// Find matches in the URL
	matches := re.FindStringSubmatch(url)
	fmt.Println(matches)
	// Check if there are enough matches
	if len(matches) != 5 {
		return filePath, internal.ErrInvalidURL
	}

	// Extract owner, repo, and filepath
	filePath.Owner = matches[1]
	filePath.Repo = matches[2]
	filePath.Branch = matches[3]
	filePath.Path = matches[4]

	return filePath, nil
}
