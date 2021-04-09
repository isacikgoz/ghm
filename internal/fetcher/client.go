package fetcher

import (
	"context"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

type Fetcher struct {
	client *github.Client
}

type Reviewer struct {
	Handle string
	Status string
}

func New(token string) *Fetcher {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return &Fetcher{client}
}

func (f *Fetcher) GetOpenPullRequests(ctx context.Context, owner, project string) ([]*github.PullRequest, error) {
	prs := make([]*github.PullRequest, 0)
	prListOpts := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 50},
	}

	for {
		ps, resp, err := f.client.PullRequests.List(ctx, owner, project, prListOpts)
		if err != nil {
			return nil, err
		}
		if resp.NextPage == 0 {
			break
		}
		prListOpts.Page = resp.NextPage
		time.Sleep(50 * time.Millisecond)

		prs = append(prs, ps...)
	}

	return prs, nil
}

func (f *Fetcher) GetPullRequestReviewers(ctx context.Context, owner, project, submitter string, number int) ([]*Reviewer, error) {
	revs, _, err := f.client.PullRequests.ListReviews(ctx, owner, project, number, &github.ListOptions{
		PerPage: 50,
	})
	if err != nil {
		return nil, err
	}
	rs, _, err := f.client.PullRequests.ListReviewers(ctx, owner, project, number, &github.ListOptions{
		PerPage: 50,
	})
	if err != nil {
		return nil, err
	}

	var requestedReviewersMap = make(map[string]string)
	for _, r := range rs.Users {
		requestedReviewersMap[r.GetLogin()] = "PENDING"
	}

	for _, review := range revs {
		if review.GetUser().GetLogin() == submitter {
			continue
		}
		requestedReviewersMap[review.User.GetLogin()] = review.GetState()
	}

	reviewers := make([]*Reviewer, 0)

	for handle, status := range requestedReviewersMap {
		reviewers = append(reviewers, &Reviewer{handle, status})
	}

	return reviewers, nil
}
