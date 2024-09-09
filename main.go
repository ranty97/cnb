package main

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/ranty97/cnb/internal/com"
	"github.com/ranty97/cnb/internal/utils"
	"go.bug.st/serial"
)

func main() {
	App := app.New()
	mainWindow := App.NewWindow("OKS")
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
			print(parityMode)
			log.Printf("Parity mode changed: %s", s)
		},
	)

	input := widget.NewEntry()
	var inputMessage string
	inputMessageBinding := binding.BindString(&inputMessage)
	input.Bind(inputMessageBinding)
	input.SetPlaceHolder("Type text to transfer")

	sendButton := widget.NewButton("Send", func() {
		msg, err := inputMessageBinding.Get()
		if err != nil {
			log.Printf("no message provided")
		}
		com.SendData(sName, &serial.Mode{BaudRate: sSpeed, Parity: parityMode}, msg+"\n")
		rMsg, err := com.RecieveData(rName, &serial.Mode{BaudRate: rSpeed, Parity: parityMode})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Received massage: %s", rMsg)
	})

	mainWindow.SetContent(container.NewVBox(
		sDropdownWidget,
		rDropdownWidget,
		sSpeedDropdownWidget,
		rSpeedDropdownWidget,
		parityDropdownWidget,
		input,
		sendButton,
		widget.NewButton("Quit", func() {
			App.Quit()
		})))

	mainWindow.ShowAndRun()
}
