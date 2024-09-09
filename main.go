package main

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ranty97/cnb/internal/com"
	"github.com/ranty97/cnb/internal/utils"
	"go.bug.st/serial"
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

	var rSpeed = com.Speeds[0]
	var sSpeed = com.Speeds[0]

	rSpeedDropdownWidget := widget.NewSelect(
		utils.ItoaSlice(com.Speeds),
		func(s string) {
			r, _ := strconv.Atoi(s)
			rSpeed = r
			log.Printf("rSpeed changed: %d", rSpeed)
		},
	)

	sSpeedDropdownWidget := widget.NewSelect(
		utils.ItoaSlice(com.Speeds),
		func(s string) {
			r, _ := strconv.Atoi(s)
			sSpeed = r
			log.Printf("sSpeed changed: %d", sSpeed)
		},
	)

	var parityMode = serial.EvenParity

	parityDropdownWidget := widget.NewSelect(
		com.GetParities(com.ParityMap),
		func(s string) {
			parityMode = com.ParityMap[s]
			log.Printf("Parity mode changed: %s", string(parityMode))
		},
	)

	mainWindow.SetContent(container.NewVBox(
		sDropdownWidget,
		rDropdownWidget,
		sSpeedDropdownWidget,
		rSpeedDropdownWidget,
		parityDropdownWidget,
		widget.NewButton("Quit", func() {
			App.Quit()
		})))

	mainWindow.ShowAndRun()
}
