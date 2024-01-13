package handler

import (
	"context"
	"fmt"

	"github.com/ayush5588/FileScope/model"
	"github.com/google/go-github/v58/github"
)

type ghClient struct {
	*github.Client
}

// NewGitHubClient ...
func NewGitHubClient(authToken string) ghClient {
	gc := github.NewClient(nil).WithAuthToken(authToken)

	return ghClient{gc}
}

func (gh ghClient) getAllOpenPRs(f model.FileInfo) ([]model.PR, error) {
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
		date := fmt.Sprintf("%d/%d/%d", day, month, year)
		customPR := model.PR{
			Number:    fmt.Sprint(*pr.Number),
			URL:       *pr.HTMLURL,
			State:     *pr.State,
			Title:     *pr.Title,
			CreatedAt: date,
		}

		PRList = append(PRList, customPR)
	}

	return PRList, nil
}

// func (gh ghClient) getModifiedFiles(pr model.PR) ([]string, error) {

// }

// GetFileModifyingPRs returns PRs (open PRs) that are modifying the given file
func GetFileModifyingPRs(fileInfo model.FileInfo) ([]model.PR, error) {

	token := "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"

	ghClient := NewGitHubClient(token)

	// var finalPRList []model.PR

	// Get all the Open PRs for the given repo
	openPRs, err := ghClient.getAllOpenPRs(fileInfo)
	if err != nil {
		return nil, err
	}

	return openPRs, nil

	/* For each entry in openPRs, call the File_Changed_In_PR api
	   Then check if our file is in those changed files
	   		If yes then add it to the list otherwise continue
	*/
	// path := fileInfo.Path

	// fileModifyingPRChan := make(chan model.PR, 100)

	// var wg sync.WaitGroup
	// wg.Add(len(openPRs))

	// for idx := range openPRs {
	// 	pr := openPRs[idx]
	// 	fileModifyingPRChan <- pr
	// }

	// for pr := range fileModifyingPRChan {
	// 	go func(pr model.PR, path string) {
	// 		defer wg.Done()
	// 		// Call the api to get the files modified
	// 		files, err := ghClient.getModifiedFiles(pr)
	// 		if err != nil {
	// 			return
	// 		}
	// 		for _, file := range files {
	// 			if file == path {
	// 				finalPRList = append(finalPRList, pr)
	// 				break
	// 			}
	// 		}
	// 	}(pr, path)
	// }

	// go func() {
	// 	wg.Wait()
	// 	close(fileModifyingPRChan)
	// }()
}
