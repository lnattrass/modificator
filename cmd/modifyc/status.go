package main

import (
	"context"

	"github.com/google/go-github/v60/github"
)

type SetStatusCli struct {
	Token       string `env:"GITHUB_TOKEN" required:""`
	Owner       string `short:"o" default:"ed448io"`
	Repository  string `short:"r" default:"code"`
	Commit      string `short:"c"`
	Status      string `help:"one of error, failure, pending, success"`
	Description string `default:"auto-status"`
	Context     string `default:"modificator-api"`
}

func (c *SetStatusCli) Run() error {
	api := github.NewClient(nil).WithAuthToken(c.Token)

	ctx := context.Background()

	_, _, err := api.Repositories.CreateStatus(ctx, c.Owner, c.Repository, c.Commit, &github.RepoStatus{
		State:   &c.Status,
		Context: &c.Context,
	})
	if err != nil {
		return err
	}
	return nil
}
