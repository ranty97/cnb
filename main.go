package main

import (
	"log"
	"strconv"
	"time"

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
	mainWindow := App.NewWindow("Computer Networks Basics")
	mainWindow.Resize(fyne.Size{Width: 500, Height: 265})

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

	sDropdownWidget.SetSelectedIndex(0)
	rDropdownWidget.SetSelectedIndex(1)

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

	sSpeedDropdownWidget.SetSelectedIndex(0)
	rSpeedDropdownWidget.SetSelectedIndex(1)

	var parityMode = serial.EvenParity

	parityDropdownWidget := widget.NewSelect(
		com.GetParities(com.ParityMap),
		func(s string) {
			parityMode = com.ParityMap[s]
			print(parityMode)
			log.Printf("Parity mode changed: %s", s)
		},
	)

	parityDropdownWidget.SetSelectedIndex(0)

	input := widget.NewEntry()
	var inputMessage string
	inputMessageBinding := binding.BindString(&inputMessage)
	input.Bind(inputMessageBinding)
	input.SetPlaceHolder("Type text to transfer")

	receivedMessageLabel := widget.NewLabel("Received Message")

	updateLabel := func(newMessage string) {
		receivedMessageLabel.SetText(newMessage)
	}

	sendButton := widget.NewButton("Send", func() {
		msg, _ := inputMessageBinding.Get()
		var packet [][]byte
		var err error
		// change logic of sending messages
		for i := 0; i < 10; i++ {
			com.Port{
				Name:   sName,
				Speed:  sSpeed,
				Parity: parityMode,
			}.SendData([]byte(msg))
			packet, err = com.Port{
				Name:   rName,
				Speed:  rSpeed,
				Parity: parityMode,
			}.ReceivePacket()
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Print("waiting for deadline")
		time.Sleep(time.Second * 1)

		print(string(utils.ConcatenateByteSlices(packet)))

		updateLabel(string(utils.ConcatenateByteSlices(packet)))
	})

	l := container.NewVBox(
		container.NewHBox(widget.NewLabel("Tx Name:"), sDropdownWidget),
		container.NewHBox(widget.NewLabel("Rx Name:"), rDropdownWidget),
		container.NewHBox(widget.NewLabel("Tx Speed:"), sSpeedDropdownWidget),
		container.NewHBox(widget.NewLabel("Rx Speed:"), rSpeedDropdownWidget),
		container.NewHBox(widget.NewLabel("Parity mode:"), parityDropdownWidget),
		input,
		sendButton,
		receivedMessageLabel,
	)

	mainWindow.SetContent(l)
	mainWindow.ShowAndRun()
}
