package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

/*
 * @author leig
 * @date 2022/02/11/2:23 PM
 */

func getWindows() fyne.Window {

	wpmApp := app.New()
	wpmWindows := wpmApp.NewWindow("wpm")

	rootPathEntry := widget.NewEntry()
	rootPathEntry.Text = wpmConfig.RootPath
	destPathEntry := widget.NewEntry()
	destPathEntry.Text = wpmConfig.DestPath
	schemaEntry := widget.NewCheck(" all", func(b bool) {
		if b {
			wpmConfig.Schema = All
		} else {
			wpmConfig.Schema = Express
		}
	})
	if wpmConfig.Schema == All {
		schemaEntry.Checked = true
	}
	backupsButton := widget.NewButton("backups", func() {
		wpmConfig.RootPath = rootPathEntry.Text
		wpmConfig.DestPath = destPathEntry.Text
		backups(getPaths(wpmConfig))
		saveConfig(wpmConfig, cfg)
	})
	recoverButton := widget.NewButton("recover", func() {
		wpmConfig.RootPath = rootPathEntry.Text
		wpmConfig.DestPath = destPathEntry.Text
		recover(getPaths(wpmConfig), wpmConfig.RootPath)
		saveConfig(wpmConfig, cfg)
	})

	form := widget.NewForm(
		&widget.FormItem{Text: "WOW root path: ", Widget: rootPathEntry},
		&widget.FormItem{Text: "WOW backups path: ", Widget: destPathEntry},
		&widget.FormItem{Text: "schema mode: ", Widget: schemaEntry},
		&widget.FormItem{Widget: backupsButton},
		&widget.FormItem{Widget: recoverButton},
	)

	wpmWindows.SetContent(form)
	wpmWindows.Resize(fyne.NewSize(650, 200))

	return wpmWindows
}
