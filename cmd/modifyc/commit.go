package main

import (
	"context"

	"github.com/alecthomas/kong"
	"github.com/lnattrass/modificator/pkg/api"
)

type CommitCli struct {
	Owner      string               `short:"o" default:"ed448io"`
	Repository string               `short:"r" default:"code"`
	Branch     string               `short:"b" required:""`
	Path       string               `short:"p" required:""`
	File       kong.FileContentFlag `short:"f" required:""`
	Message    string               `short:"m" help:"Commit message" default:"Automated change"`
	Token      string               `env:"GITHUB_TOKEN" required:""`
	CreatePR   bool                 `short:"c" help:"create a PR too"`
	MergePR    bool                 `help:"merge the PR"`
}

func (c *CommitCli) Run() error {
	ctx := context.Background()
	return api.Commit(ctx, c.Token, c.Owner, c.Repository, c.Branch, c.Path, c.Message, c.File, c.CreatePR, c.MergePR)
}
