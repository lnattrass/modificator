package main

import (
	"context"

	"github.com/alecthomas/kong"
	"github.com/lnattrass/modificator/pkg/api"
)

type PatchCli struct {
	Owner      string `short:"o" default:"ed448io"`
	Repository string `short:"r" default:"test"`
	Branch     string `short:"b" required:""`
	Path       string `short:"p" required:""`
	Message    string `short:"m" help:"Commit message" default:"Automated change"`

	Patch kong.FileContentFlag `short:"f" required:""`

	Token string `env:"GITHUB_TOKEN" required:""`

	CreatePR bool `short:"c" help:"create a PR too"`
	MergePR  bool `help:"merge the PR"`
}

func (c *PatchCli) Run() error {
	ctx := context.Background()

	return api.Patch(ctx, c.Token, c.Owner, c.Repository, c.Branch, c.Path, c.Message, c.Patch, c.CreatePR, c.MergePR)
}
