package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/hayashiki/ghutils/notify"

	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hayashiki/ghutils/csv"
	"github.com/hayashiki/ghutils/v4client"
)

var (
	token   string
	project string
	owner   string
	repo    string
)

func main() {
	fs := flag.NewFlagSet("exporter", flag.ExitOnError)
	fs.StringVar(&project, "project", "", "github project name (required)")
	fs.StringVar(&owner, "owner", "", "github owner (required)")
	fs.StringVar(&repo, "repo", "", "github repository name (required)")
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
		if token != "" {
			fmt.Fprintln(os.Stderr, errors.New("token args are required or GH_SECRET_TOKEN env is required"))
			os.Exit(1)
		}
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
	var outs []csv.Output
	message := project + "の経過レポートです\n"
	if len(result.Repository.Projects.Edges) == 0 {
		return err
	}
	for _, col := range result.Repository.Projects.Edges[0].Node.Columns.Edges {
		message = fmt.Sprintf("%s[%s]\n", message, col.Node.Name)
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

			message = fmt.Sprintf("%s*<%v|%v> * assignee: %v\n", message, out.URL, out.Title, out.Assignees)
		}
	}
	//log.Printf("message %s", message)
	if os.Getenv("SLACK_INCOMING_WEBHOOK") != "" {
		notify := notify.New(os.Getenv("SLACK_INCOMING_WEBHOOK"))
		if err := notify.Do(message); err != nil {
			return err
		}
	}
	return nil
}