package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/isacikgoz/ghm/internal/fetcher"
	"github.com/isacikgoz/ghm/internal/render"
	"github.com/rodaine/table"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "Expected at least 4 arguments, but got %d\n", len(os.Args)-1)
		printUsage()
		os.Exit(1)
	}

	if os.Args[1] != "pr" {
		os.Exit(0)
	}

	accessToken := os.Args[4]
	owner := os.Args[2]
	project := os.Args[3]

	f := fetcher.New(accessToken)
	ctx := context.Background()

	cancel := render.StartSimpleProgress(ctx, os.Stdout, "fetching pull requests of "+strings.Join(os.Args[2:4], "/")+"...")
	defer cancel()

	prs, err := f.GetOpenPullRequests(ctx, owner, project)
	checkError(err)

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	tbl := table.New("Owner", "Last Update", "URL", "Reviewers")
	tbl.WithHeaderFormatter(headerFmt)

	for _, pr := range prs {
		reviewers, err := f.GetPullRequestReviewers(ctx, owner, project, pr.GetUser().GetLogin(), pr.GetNumber())
		checkError(err)

		tbl.AddRow(pr.GetUser().GetLogin(), &render.LastUpdate{Date: pr.GetUpdatedAt()}, pr.GetHTMLURL(), &render.Reviewers{Reviewers: reviewers})
	}
	cancel()

	fmt.Println()
	tbl.Print()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Program exited with error: %s\n", err.Error())
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: ghm pr <repository owner> <repository name> <github access token>")
}
