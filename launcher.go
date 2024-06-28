package main

import (
	"fmt"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Launcher struct {
	banner      *canvas.Image
	title       *widget.Label
	tagSelect   *widget.Select
	assetSelect *widget.Select
	assetInfo   [3]*widget.Label

	tags     []string
	assets   []string
	assetMap map[string]map[string]interface{}

	assetDownloadTag  string
	assetDownloadName string
	assetDownloadURL  string
	assetDownloadSize int64

	assetDownloadDest     fyne.ListableURI
	assetDownloadButton   *widget.Button
	assetDownloadProgress *widget.Label

	isDownloading bool

	window fyne.Window
}

func (l *Launcher) LoadUI(_app fyne.App) {
	l.banner = canvas.NewImageFromResource(resourceBannerJpg)
	l.banner.FillMode = canvas.ImageFillContain
	l.banner.SetMinSize(fyne.Size{Width: 400, Height: 100})

	l.title = widget.NewLabel("86 Game Launcher")
	l.title.Alignment = fyne.TextAlignCenter
	l.title.TextStyle = fyne.TextStyle{Bold: true}

	// Tags
	tags, err := GetGithubReleasesTags()
	if err != nil {
		fmt.Println(err)
		_app.SendNotification(fyne.NewNotification("86 Game Launcher: Tags Error", err.Error()))
		_app.Quit()
	}
	l.tags = tags
	l.tagSelect = widget.NewSelect(l.tags, nil)
	l.tagSelect.SetSelectedIndex(0)
	l.tagSelect.OnChanged = func(s string) {
		// tag selector changed
		go l.GetGithubReleasesAssetsThread(_app, s)
	}

	// Assets
	assets, assetMap, err := GetGithubReleasesAssets(l.tags[0])
	if err != nil {
		fmt.Println(err)
		_app.SendNotification(fyne.NewNotification("86 Game Launcher: Assets Error", err.Error()))
		_app.Quit()
	}
	l.assets = assets
	l.assetMap = assetMap
	l.assetSelect = widget.NewSelect(l.assets, func(s string) {
		// asset selector changed
		if s != "" {
			assetSize := l.assetMap[l.assets[l.assetSelect.SelectedIndex()]]["size"]
			assetSizeMB := float64(assetSize.(int64)) / (1024 * 1024)
			assetSizeMBRounded := math.Round(assetSizeMB*1000) / 1000
			l.assetDownloadProgress.SetText(fmt.Sprintf("Progress: 0 Mb / %0.3f Mb", assetSizeMBRounded))

			l.assetInfo[0].SetText(fmt.Sprintf("Name: %s", l.assets[l.assetSelect.SelectedIndex()]))
			l.assetInfo[1].SetText(fmt.Sprintf("URL: %s", l.assetMap[s]["url"]))
			l.assetInfo[2].SetText(fmt.Sprintf("Size: %d", l.assetMap[s]["size"]))
		} else {
			l.assetDownloadProgress.SetText("Progress:")

			l.assetInfo[0].SetText("Name:")
			l.assetInfo[1].SetText("URL:")
			l.assetInfo[2].SetText("Size:")
		}
	})
	l.assetSelect.PlaceHolder = "Select asset file"

	// Asset Info init
	l.assetInfo[0] = widget.NewLabel("Name:")
	l.assetInfo[1] = widget.NewLabel("URL:")
	l.assetInfo[2] = widget.NewLabel("Size:")

	// Asset Download init
	l.assetDownloadButton = widget.NewButton("Download", func() {
		// Runs a download thread
		if !l.isDownloading {
			l.isDownloading = true
			dialog.ShowFolderOpen(func(lu fyne.ListableURI, err error) {
				if err != nil {
					fmt.Println(err)
					_app.SendNotification(fyne.NewNotification("86 Game Launcher: Download Path Error", err.Error()))
				}
				if lu == nil {
					fmt.Println("No dest selected")
				} else {
					fmt.Println("Selected dest =", lu)
					l.assetDownloadDest = lu

					l.assetDownloadTag = l.tags[l.tagSelect.SelectedIndex()]
					l.assetDownloadName = l.assets[l.assetSelect.SelectedIndex()]
					l.assetDownloadURL = l.assetMap[l.assets[l.assetSelect.SelectedIndex()]]["url"].(string)
					l.assetDownloadSize = l.assetMap[l.assets[l.assetSelect.SelectedIndex()]]["size"].(int64)

					go l.DownloadFileThread(_app)
				}
			}, l.window)
		} else {
			_app.SendNotification(fyne.NewNotification("86 Game Launcher: Download Status", "Download is already in progress!"))
		}
	})
	l.assetDownloadProgress = widget.NewLabel("Progress:")

	// Window setup
	l.window = _app.NewWindow("86 Game Launcher")
	l.window.SetContent(container.NewVBox(
		l.banner,
		l.title,
		l.tagSelect,
		l.assetSelect,
		container.NewHScroll(
			container.NewVBox(
				l.assetInfo[0],
				l.assetInfo[1],
				l.assetInfo[2],
			),
		),
		container.NewHBox(
			l.assetDownloadButton,
			l.assetDownloadProgress,
		),
	))

	// Resize and center the window
	l.window.Resize(fyne.Size{Width: 400, Height: 450})
	l.window.CenterOnScreen()
	l.window.Show()
}

func NewLauncher() *Launcher {
	return &Launcher{}
}
