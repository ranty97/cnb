package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ranty97/cnb/internal/com"
)

func main() {
	App := app.New()
	mainWindow := App.NewWindow("LAB 1")
	mainWindow.Resize(fyne.Size{Width: 500, Height: 500})

	var sName, rName string

	sDropdownWidget := widget.NewSelect(
		com.GetPorts(),
		func(s string) {
			sName = s
			log.Printf("Sender changed: %s", sName)
		},
	)

	rDropdownWidget := widget.NewSelect(
		com.GetPorts(),
		func(r string) {
			rName = r
			log.Printf("Receiver changed: %s", rName)
		},
	)

	mainWindow.SetContent(container.NewVBox(
		widget.NewLabel(fmt.Sprint(com.GetPorts())),
		widget.NewButton("Quit", func() {
			App.Quit()
		})))

	mainWindow.ShowAndRun()
}
