package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

var lamps []*lamp

func controlLamps(msg midi.Message, timestampms int32) {
	var channel, cc, val uint8
	switch {
	case msg.GetControlChange(&channel, &cc, &val):
		switch cc {
		case 17:
			lamps[0].SetVal(val)
		case 18:
			lamps[1].SetVal(val)
		case 19:
			lamps[2].SetVal(val)
		case 20:
			lamps[3].SetVal(val)
		}
	}
}

func main() {
	defer midi.CloseDriver()

	flag.Parse()

	in := readPort()

	bridge := findBridge()
	lights, err := bridge.GetLights()
	if err != nil {
		log.Fatal(err)
	}

	createMapping(in, lights)

	lamps = []*lamp{Lamp(&lights[0]), Lamp(&lights[1]), Lamp(&lights[2]), Lamp(&lights[3])}

	for _, lamp := range lamps {
		go lamp.Run()
	}

	stopMidi, err := midi.ListenTo(in, controlLamps)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Press ENTER to exit")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	stopMidi()

	for _, lamp := range lamps {
		go lamp.Close()
	}

	for _, lamp := range lamps {
		fmt.Println(lamp)
	}
}
