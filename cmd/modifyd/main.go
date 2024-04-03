package main

import (
	"github.com/alecthomas/kong"
	"github.com/gin-gonic/gin"
	"github.com/lnattrass/modificator/pkg/api"
)

type Cli struct {
	ListenAddr string `default:":8000"`
}

func (c *Cli) Run() error {
	router := gin.Default()

	g := router.Group("/api/v1")
	api.ConfigureRoutes(g)

	router.Run(c.ListenAddr)
	return nil
}

func main() {
	ctx := kong.Parse(&Cli{})
	ctx.FatalIfErrorf(ctx.Run())
}
