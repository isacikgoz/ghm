package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/isacikgoz/ghm/internal/fetcher"
	"github.com/isacikgoz/ghm/internal/printer"
	"github.com/rodaine/table"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Expected at least 3 arguments, but got %d\n", len(os.Args)-1)
		printUsage()
		os.Exit(1)
	}

	accessToken := os.Args[3]
	owner := os.Args[1]
	project := os.Args[2]

	f := fetcher.New(accessToken)
	ctx := context.Background()

	cancel := printer.StartSimpleProgress(ctx, os.Stdout, "fetching "+strings.Join(os.Args[1:3], "/")+"...")
	defer cancel()

	prs, err := f.GetOpenPullRequests(ctx, owner, project)
	checkError(err)

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	tbl := table.New("Owner", "Last Update", "URL", "Reviewers")
	tbl.WithHeaderFormatter(headerFmt)

	for _, pr := range prs {
		reviewers, err := f.GetPullRequestReviewers(ctx, owner, project, pr.GetUser().GetLogin(), pr.GetNumber())
		checkError(err)

		tbl.AddRow(pr.GetUser().GetLogin(), &printer.LastUpdate{Date: pr.GetUpdatedAt()}, pr.GetHTMLURL(), &printer.Reviewers{Reviewers: reviewers})
	}

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
	fmt.Println("Usage: ghm <repository owner> <repository name> <github access token>")
}
