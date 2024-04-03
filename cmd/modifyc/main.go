package main

import "github.com/alecthomas/kong"

type Cli struct {
	Local struct {
		Commit    CommitCli    `cmd:""`
		SetStatus SetStatusCli `cmd:""`
		PatchYaml PatchCli     `cmd:""`
	} `cmd:""`
}

func main() {
	ctx := kong.Parse(&Cli{})
	ctx.FatalIfErrorf(ctx.Run())
}
