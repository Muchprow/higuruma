package main

import (
	"embed"
	"fmt"

	"github.com/Muchprow/hakari"
)

//go:embed assets
var assetsFS embed.FS

func main() {
	app := hakari.New("Higuruma")

	app.Screen("home", func(s *hakari.Screen) {
		s.HTML("assets/screens/home/home.html")

		s.On("gotoMods", func() { s.GoTo("mods") })
		s.On("gotoInstalled", func() { s.GoTo("installed") })
		s.On("gotoProfiles", func() { s.GoTo("profiles") })
	})

	app.Screen("mods", func(s *hakari.Screen) {
		s.HTML("assets/screens/mods/mods.html")
		s.On("back", func() { s.GoTo("home") })
	})

	app.Screen("installed", func(s *hakari.Screen) {
		s.HTML("assets/screens/installed/installed.html")
		s.On("back", func() { s.GoTo("home") })
	})

	app.Screen("profiles", func(s *hakari.Screen) {
		s.HTML("assets/screens/profiles/profile.html")
		s.On("back", func() { s.GoTo("home") })
	})

	if err := app.LoadHTMLFromFS(assetsFS, "assets/index.html"); err != nil {
		fmt.Println("ERROR loading HTML:", err)
		panic(err)
	}

	app.Start("home")
}
