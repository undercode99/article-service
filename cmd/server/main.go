package main

import (
	"context"

	"github.com/undercode99/article_service/cmd/server/runner"
)

func main() {
	ctx := context.Background()
	app := runner.InitializeApp(ctx)
	app.Run(ctx)
}
