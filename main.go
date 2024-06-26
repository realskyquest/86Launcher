package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	_app := app.NewWithID("com.github.realskyquest_86launcher")

	_app.SetIcon(resourceIconPng)

	l := NewLauncher()
	l.LoadUI(_app)
	_app.Run()
}
