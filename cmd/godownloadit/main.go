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
	DEFAULT_SIZE                 = fyne.NewSize(520, 175)
	OPENING_FILE_SIZE            = fyne.NewSize(620, 504)
	APP_LABEL                    = "Download it!"
	DEFAULT_DOWNLOAD_SPEED_LABEL = "Avg speed: 0"
	DOWNLOADING_LABEL            = "Downloading..."
)

type DownloadIt struct {
	app                fyne.App
	window             fyne.Window
	mainTitle          *widget.Label
	urlEntry           *widget.Entry
	pBar               *widget.ProgressBar
	downloadBtn        *widget.Button
	downloadSpeedLabel *widget.Label
}

func main() {
	app := NewApp()

	app.setupContent()
	app.run()
}

func NewApp() *DownloadIt {
	app := app.New()
	window := app.NewWindow(APP_LABEL)
	mainTitle := widget.NewLabel(APP_LABEL)
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter URL...")

	pBar := widget.NewProgressBar()
	pBar.Max = 99.99

	downloadSpeedLabel := widget.NewLabel(DEFAULT_DOWNLOAD_SPEED_LABEL)

	return &DownloadIt{
		app:                app,
		window:             window,
		mainTitle:          mainTitle,
		urlEntry:           urlEntry,
		pBar:               pBar,
		downloadSpeedLabel: downloadSpeedLabel,
	}
}

func (app *DownloadIt) resizeWindow(size fyne.Size) {
	app.window.Resize(fyne.NewSize(size.Width, size.Height))
}

func (app *DownloadIt) run() {
	app.window.ShowAndRun()
}

func (app *DownloadIt) setupContent() {
	app.resizeWindow(DEFAULT_SIZE)
	app.setupDownload()
	app.window.SetContent(container.NewVBox(
		app.mainTitle,
		app.urlEntry,
		app.pBar,
		app.downloadSpeedLabel,
		app.downloadBtn,
	))
}

func (app *DownloadIt) setupDownload() {
	app.downloadBtn = widget.NewButton("Download", func() {
		app.downloadFile()
	})
}

func (app *DownloadIt) downloadFile() {
	app.updateStartAndFinishedState(true)
	app.resizeWindow(OPENING_FILE_SIZE)

	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		defer app.resizeWindow(DEFAULT_SIZE)

		if err != nil || writer == nil {
			app.mainTitle.SetText(APP_LABEL)
			return
		}
		defer writer.Close()

		fileToSave := writer.URI().Path()
		fmt.Printf("File to save: %s - %s\n", app.urlEntry.Text, fileToSave)
		downloadSpeedStr := DEFAULT_DOWNLOAD_SPEED_LABEL

		go func() {
			g := got.New()
			g.ProgressFunc = func(d *got.Download) {
				progress := float64(d.Size()) / float64(d.TotalSize()) * 100
				fyne.Do(func() {

					if d.Speed() > 0 {
						downloadSpeedStr = fmt.Sprintf("Avg speed: %s", formatDownloadSpeed(d.Speed()))
					}

					app.pBar.SetValue(progress)
					app.downloadSpeedLabel.SetText(downloadSpeedStr)
				})
			}

			err := g.Download(app.urlEntry.Text, fileToSave)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error downloading file: %s", err), app.window)
				app.updateStartAndFinishedState(false)
				return
			}

			app.updateStartAndFinishedState(false)
		}()
	}, app.window)
}

func (app *DownloadIt) updateDownloadBtnText(text string, disable bool) {
	fyne.Do(func() {
		app.downloadBtn.SetText(text)

		if disable {
			app.downloadBtn.Disable()
			return
		}

		app.downloadBtn.Enable()
	})
}

func (app *DownloadIt) updateStartAndFinishedState(starting bool) {
	if starting {
		app.mainTitle.SetText(DOWNLOADING_LABEL)
		app.updateDownloadBtnText(DOWNLOADING_LABEL, true)
		return
	}

	app.mainTitle.SetText(APP_LABEL)
	dialog.ShowInformation("Download completed", "File downloaded successfully!", app.window)
	app.updateDownloadBtnText("Download", false)
	app.pBar.SetValue(0.0)
	app.downloadSpeedLabel.SetText(DEFAULT_DOWNLOAD_SPEED_LABEL)
}

func formatDownloadSpeed(bps uint64) string {
	const unit = 1024
	if bps < unit {
		return fmt.Sprintf("%d B/s", bps)
	}
	div, exp := unit, 0
	for n := bps / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB/s", float64(bps)/float64(div), "KMGTPE"[exp])
}
