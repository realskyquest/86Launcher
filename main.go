package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
<<<<<<< Updated upstream
	_app := app.New()
=======
	_app := app.NewWithID("com.github.realskyquest_86launcher")
>>>>>>> Stashed changes
	_app.SetIcon(resourceIconPng)

	l := NewLauncher()
	l.LoadUI(_app)
	_app.Run()
}
