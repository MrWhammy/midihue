package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/amimof/huego"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

var lampen []*lamp

func onMessage(msg midi.Message, timestampms int32) {
	var channel, cc, val uint8
	switch {
	case msg.GetControlChange(&channel, &cc, &val):
		switch cc {
		case 17:
			lampen[0].SetVal(val)
		case 18:
			lampen[1].SetVal(val)
		case 19:
			lampen[2].SetVal(val)
		case 20:
			lampen[3].SetVal(val)
		}
	}
}

func main() {
	defer midi.CloseDriver()

	bridge := findBridge()
	lights, err := bridge.GetLights()
	if err != nil {
		log.Fatal(err)
	}
	for index, light := range lights {
		fmt.Println(index, ": ", light)
	}
	lampen = []*lamp{Lamp(&lights[0]), Lamp(&lights[1]), Lamp(&lights[2]), Lamp(&lights[3])}
	in := readPort()

	for _, lamp := range lampen {
		go lamp.Run()
	}

	stopMidi, err := midi.ListenTo(in, onMessage)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Press ENTER to exit")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	stopMidi()

	for _, lamp := range lampen {
		go lamp.Close()
	}

	for _, lamp := range lampen {
		fmt.Println(lamp)
	}
}

func readPort() drivers.In {
	var port drivers.In
	for port == nil {
		fmt.Printf("MIDI IN Ports\n")
		fmt.Println(midi.GetInPorts())
		fmt.Print("Enter the port to listen to: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		parsedPort, err := strconv.ParseInt(scanner.Text(), 10, 8)
		if err != nil {
			fmt.Println(err)
		} else {
			openedPort, err := midi.InPort(int(parsedPort))
			if err != nil {
				fmt.Println(err)
			} else {
				port = openedPort
			}
		}
	}
	return port
}

func findBridge() *huego.Bridge {
	bridge, err := huego.Discover()
	if err != nil {
		log.Fatal(err)
	}
	user, err := bridge.CreateUser("midihue") // Link button needs to be pressed
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("New user: " + user)
	}
	return bridge.Login(user)
}
