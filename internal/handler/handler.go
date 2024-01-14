package handler

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/ayush5588/FileScope/model"
	"github.com/google/go-github/v58/github"
	"go.uber.org/zap"
)

// GitHubClient ...
type GitHubClient struct {
	*github.Client
}

// NewGitHubClient ...
func NewGitHubClient(authToken string) GitHubClient {
	gc := github.NewClient(nil).WithAuthToken(authToken)

	return GitHubClient{gc}
}

func (gh GitHubClient) getAllOpenPRs(f model.FileInfo) ([]model.PR, error) {
	var PRList []model.PR
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var prList []*github.PullRequest

	// For pagination support
	for {
		pullRequests, resp, err := gh.PullRequests.List(context.Background(), f.Owner, f.Repo, opt)
		if err != nil {
			return nil, err
		}

		prList = append(prList, pullRequests...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	for _, pr := range prList {
		year, month, day := pr.CreatedAt.Date()
		date := model.CreationDate{
			Day:      fmt.Sprint(day),
			Month:    month.String(),
			Year:     fmt.Sprint(year),
			FullDate: fmt.Sprintf("%d/%d/%d", day, month, year),
		}
		customPR := model.PR{
			Number:    fmt.Sprint(*pr.Number),
			URL:       *pr.HTMLURL,
			State:     *pr.State,
			Title:     *pr.Title,
			CreatedBy: pr.Head.User.GetLogin(),
			Branch:    pr.Head.GetLabel(),
			CreatedOn: date,
		}

		PRList = append(PRList, customPR)
	}

	return PRList, nil
}

func (gh GitHubClient) getModifiedFiles(pr model.PR, fileInfo model.FileInfo) ([]model.File, error) {
	var files []model.File

	prNumber, _ := strconv.Atoi(pr.Number)

	commitFiles, _, err := gh.PullRequests.ListFiles(context.Background(), fileInfo.Owner, fileInfo.Repo, prNumber, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	for i := range commitFiles {
		cf := commitFiles[i]
		file := model.File{
			Name:   *cf.Filename,
			Status: *cf.Status,
			SHA:    *cf.SHA,
		}
		files = append(files, file)
	}

	return files, nil

}

// GetFileModifyingPRs returns PRs (open PRs) that are modifying the given file
func GetFileModifyingPRs(logger *zap.SugaredLogger, fileInfo model.FileInfo) ([]model.PR, error) {
	logger.Info("inside GetFileModifyingPRs...")

	token := ""

	ghClient := NewGitHubClient(token)

	finalPRList := make([]model.PR, 0)

	// Get all the Open PRs for the given repo
	openPRs, err := ghClient.getAllOpenPRs(fileInfo)
	if err != nil {
		return nil, err
	}

	/* For each entry in openPRs, call the File_Changed_In_PR api
	   Then check if our file is in those changed files
	   		If yes then add it to the list otherwise continue
	*/
	path := fileInfo.Path

	var wg sync.WaitGroup
	wg.Add(len(openPRs))

	for _, pr := range openPRs {
		go func(fileInfo model.FileInfo, pr model.PR, path string) {
			defer wg.Done()
			// Call the api to get the files modified
			files, err := ghClient.getModifiedFiles(pr, fileInfo)
			if err != nil {
				return
			}
			logger.Info("modified files recieved")
			for _, file := range files {
				if file.Name == path {
					logger.Infof("found PR match - %s", pr.Number)
					finalPRList = append(finalPRList, pr)
					break
				}
			}
		}(fileInfo, pr, path)
	}

	wg.Wait()

	return finalPRList, nil

}
