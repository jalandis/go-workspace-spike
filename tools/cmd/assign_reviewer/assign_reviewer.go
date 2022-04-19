package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

const (
	ProjectOwner string = "b-lab-org"

	// TODO: Switch to real repository
	Repository string = "go-workspace-spike"
)

func main() {
	var teamName string
	flag.StringVar(&teamName, "teamName", "", "team name of available reviewers")

	var pullRequestID int
	flag.IntVar(&pullRequestID, "pullRequestID", -1, "pull request ID")
	flag.Parse()

	if flag.NArg() > 0 {
		flag.PrintDefaults()
		log.Fatal("Unknown arguments")
	}

	if teamName == "" {
		log.Fatal("teamName parameter is required")
	}

	if pullRequestID <= 0 {
		log.Fatal("pullRequestID parameter is required")
	}

	token := os.Getenv("OAUTH_TOKEN")
	if token == "" {
		log.Fatal("Missing OAUTH_TOKEN")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	members, _, err := client.Teams.ListTeamMembersBySlug(ctx, ProjectOwner, teamName, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to list team members: %w", err))
	}

	CurrentReviewers := map[string]int{}
	for _, member := range members {
		if member.Login == nil {
			continue
		}

		CurrentReviewers[*member.Login] = 0
	}

	// TODO: Replace with real ProjectOwner + Repository
	// pullRequests, _, err := client.PullRequests.List(ctx, ProjectOwner, Repository , nil)
	pullRequests, _, err := client.PullRequests.List(ctx, "jalandis", Repository, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to list open pull requests: %w", err))
	}

	for _, pr := range pullRequests {
		for _, requestedReviewer := range pr.RequestedReviewers {
			if requestedReviewer.Login == nil {
				continue
			}

			CurrentReviewers[*requestedReviewer.Login] = CurrentReviewers[*requestedReviewer.Login] + 1
		}
	}

	// TODO: Add test for random selection of user with least number of reviews assigned
	minReviews := math.MaxInt
	reviewerOptions := []string{}
	for reviewer, reviewCount := range CurrentReviewers {
		if minReviews > reviewCount {
			minReviews = reviewCount
			reviewerOptions = []string{}
		}

		if minReviews == reviewCount {
			reviewerOptions = append(reviewerOptions, reviewer)
		}
	}

	randomUser := reviewerOptions[rand.Intn(len(reviewerOptions))]
	fmt.Printf("Randomly assigned reviewer: %s\n", randomUser)

	// TODO: Fix hardcoded username
	_, _, err = client.PullRequests.RequestReviewers(ctx, "jalandis", Repository, pullRequestID, github.ReviewersRequest{
		// Reviewers: []string{randomUser},
		Reviewers: []string{"umohm1"},
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to request reviewer: %w", err))
	}
}
