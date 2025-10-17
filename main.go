package main

import (
	"os"
	"time"

	"github.com/rinarudhei/pomcli/app"
)

func main() {
	a, err := app.NewApp(50*time.Minute, 10*time.Minute, 60*time.Minute)
	if err != nil {
		os.Exit(1)
	}
	if err := a.Run(); err != nil {
		os.Exit(1)
	}
}
