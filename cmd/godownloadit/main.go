package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/melbahja/got"
)

var (
	DEFAULT_SIZE      = fyne.NewSize(520, 175)
	OPENING_FILE_SIZE = fyne.NewSize(620, 504)
)

type MyApp struct {
	app         fyne.App
	window      fyne.Window
	mainTitle   *widget.Label
	urlEntry    *widget.Entry
	pBar        *widget.ProgressBar
	downloadBtn *widget.Button
}

func main() {
	app := NewApp()

	app.setupContent()
	app.run()
}

func NewApp() *MyApp {
	app := app.New()
	window := app.NewWindow("Download It!")
	mainTitle := widget.NewLabel("Download It!")
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter URL...")

	pBar := widget.NewProgressBar()
	pBar.Max = 99.99

	return &MyApp{
		app:       app,
		window:    window,
		mainTitle: mainTitle,
		urlEntry:  urlEntry,
		pBar:      pBar,
	}
}

func (app *MyApp) resizeWindow(size fyne.Size) {
	app.window.Resize(fyne.NewSize(size.Width, size.Height))
}

func (app *MyApp) run() {
	app.window.ShowAndRun()
}

func (app *MyApp) setupContent() {
	app.resizeWindow(DEFAULT_SIZE)
	app.setupDownload()
	app.window.SetContent(container.NewVBox(app.mainTitle, app.urlEntry, app.pBar, app.downloadBtn))
}

func (app *MyApp) setupDownload() {
	app.downloadBtn = widget.NewButton("Download", func() {
		app.downloadFile()
	})
}

func (app *MyApp) downloadFile() {
	app.mainTitle.SetText("Downloading...")
	app.resizeWindow(OPENING_FILE_SIZE)
	app.updateDownloadBtnText("Downloading...", true)
	defer app.updateDownloadBtnText("Download", false)

	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		defer app.resizeWindow(DEFAULT_SIZE)

		if err != nil || writer == nil {
			app.mainTitle.SetText("Download It!")
			return
		}
		defer writer.Close()

		fileToSave := writer.URI().Path()
		fmt.Printf("File to save: %s - %s\n", app.urlEntry.Text, fileToSave)

		go func() {
			g := got.New()
			g.ProgressFunc = func(d *got.Download) {
				progress := float64(d.Size()) / float64(d.TotalSize()) * 100
				fyne.Do(func() {
					fmt.Println(progress)
					app.pBar.SetValue(progress)
				})
			}

			err := g.Download(app.urlEntry.Text, fileToSave)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error downloading file: %s", err), app.window)
				return
			}

			dialog.ShowInformation("Download completed", "File downloaded successfully!", app.window)
		}()
	}, app.window)
}

func (app *MyApp) updateDownloadBtnText(text string, disable bool) {
	fyne.Do(func() {
		app.downloadBtn.SetText(text)

		if disable {
			app.downloadBtn.Disable()
			return
		}

		app.downloadBtn.Enable()
	})
}
