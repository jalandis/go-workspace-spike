package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

func main() {
	var teamName string
	flag.StringVar(&team, "teamName", "", "team name of available reviewers")

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

	token := os.Getenv("GITHUB_OAUTH_TOKEN")
	if token == "" {
		log.Fatal("Missing GITHUB_OAUTH_TOKEN")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	members, _, err := client.Teams.ListTeamMembersBySlug(ctx, "b-lab-org", teamName, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to list team members: %w", err))
	}

	for _, member := range members {
		if member.Login == nil {
			continue
		}

		if *member.Login == "the-impact-bot" {
			continue
		}

		client.PullRequests.RequestReviewers(ctx, "jalandis", "go-workspace-spike", pullRequestID, github.ReviewersRequest{
			Reviewers: []string{*member.Login},
		})
		return
	}
}
