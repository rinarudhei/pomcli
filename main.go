package main

import (
	"os"
	"time"

	"github.com/rinarudhei/pomcli/app"
	"github.com/rinarudhei/pomcli/session"
	"github.com/rinarudhei/pomcli/session/repository"
)

func main() {
	repo := repository.NewRepository()
	s := session.NewSession(repo, 25*time.Minute, 5*time.Minute, 15*time.Minute)
	a, err := app.NewApp(s)
	if err != nil {
		os.Exit(1)
	}
	if err := a.Run(); err != nil {
		os.Exit(1)
	}
}
