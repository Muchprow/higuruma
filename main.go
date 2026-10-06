package main

import (
	"embed"
	"fmt"

	"github.com/Muchprow/hakari"
)

//go:embed assets
var assetsFS embed.FS

func main() {
	fmt.Println(">>> 1. Program started")

	app := hakari.New("Higuruma")
	fmt.Println(">>> 2. App created")

	app.Screen("home", func(s *hakari.Screen) {
		fmt.Println(">>> 3. Registering home")
		s.HTML("assets/screens/home/home.html")

		s.On("gotoMods", func() { s.GoTo("mods") })
		s.On("gotoInstalled", func() { s.GoTo("installed") })
		s.On("gotoProfiles", func() { s.GoTo("profiles") })
	})

	app.Screen("mods", func(s *hakari.Screen) {
		fmt.Println(">>> 4. Registering mods")
		s.HTML("assets/screens/mods/mods.html")
		s.On("back", func() { s.GoTo("home") })
	})

	app.Screen("installed", func(s *hakari.Screen) {
		fmt.Println(">>> 5. Registering installed")
		s.HTML("assets/screens/installed/installed.html")
		s.On("back", func() { s.GoTo("home") })
	})

	app.Screen("profiles", func(s *hakari.Screen) {
		fmt.Println(">>> 6. Registering profiles")
		s.HTML("assets/screens/profiles/profile.html")
		s.On("back", func() { s.GoTo("home") })
	})

	fmt.Println(">>> 7. Loading index.html")
	if err := app.LoadHTMLFromFS(assetsFS, "assets/index.html"); err != nil {
		fmt.Println(">>> ERROR loading HTML:", err)
		panic(err)
	}
	fmt.Println(">>> 8. Index loaded")

	fmt.Println(">>> 9. Starting app with screen 'home'")
	app.Start("home")
	fmt.Println(">>> 10. App finished")
}
