package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hayashiki/ghutils/csv"
	"github.com/hayashiki/ghutils/v4client"
)

var (
	token   string
	project string
	owner   string
	repo    string
	dir     string
)

func main() {
	fs := flag.NewFlagSet("exporter", flag.ExitOnError)
	fs.StringVar(&project, "project", "", "github project name (required)")
	fs.StringVar(&owner, "owner", "", "github owner (required)")
	fs.StringVar(&repo, "repo", "", "github repository name (required)")
	fs.StringVar(&repo, "dir", "", "output directory")
	fs.StringVar(&token, "token", "", "github personal token (required)")
	if err := fs.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if owner == "" || repo == "" {
		fs.PrintDefaults()
		fmt.Fprintln(os.Stderr, errors.New("name, email, password args are required"))
		os.Exit(1)
	}

	if token == "" {
		token = os.Getenv("GH_SECRET_TOKEN")
		if token == "" {
			fmt.Fprintln(os.Stderr, errors.New("token args are required or GH_SECRET_TOKEN env is required"))
			os.Exit(1)
		}
	}

	if dir == "" {
		dir = "."
	}

	if err := Run(project, owner, repo, token); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Run(project, owner, repo, token string) error {
	client, err := v4client.NewClient(http.DefaultClient, token)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	result, err := client.Issue.Get(owner, repo, project)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	var outs []csv.Output
	if len(result.Repository.Projects.Edges) == 0 {
		return err
	}
	for _, col := range result.Repository.Projects.Edges[0].Node.Columns.Edges {
		for _, node := range col.Node.Cards.Edges {
			card := node.Node.Content
			var out csv.Output
			var assignees []string
			for _, assignee := range card.Assignees.Edges {
				assignees = append(assignees, assignee.Node.Login)
			}
			var labels []string
			for _, label := range card.Labels.Edges {
				labels = append(labels, label.Node.Name)
			}
			assigneesStr := strings.Join(assignees, ",")
			labelsStr := strings.Join(labels, ",")
			out.Title = card.Title
			out.Number = card.Number
			out.URL = card.Url
			out.Column = col.Node.Name
			out.Assignees = assigneesStr
			out.Labels = labelsStr
			outs = append(outs, out)
		}
	}
	c := csv.New()
	res, err := c.Generate(outs)
	if err != nil {
		return err
	}
	tStr := timeToString(time.Now())
	file := fmt.Sprintf("%s/export-%s.csv", dir, tStr)
	outFile, err := os.Create(file)
	defer outFile.Close()
	_, err = io.Copy(outFile, res)
	return err
}

var layout = "20060102150405"

func timeToString(t time.Time) string {
	str := t.Format(layout)
	return str
}
