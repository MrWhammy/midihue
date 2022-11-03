package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/amimof/huego"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

var mapping map[uint8][]int

var lastTouchedControl int16

func detectControl(msg midi.Message, timestampms int32) {
	var channel, cc, val uint8
	switch {
	case msg.GetControlChange(&channel, &cc, &val):
		previousControl := lastTouchedControl
		if previousControl > 99 {
			fmt.Printf("\x08")
		}
		if previousControl > 9 {
			fmt.Printf("\x08")
		}
		if previousControl > 0 {
			fmt.Printf("\x08")
		}
		lastTouchedControl = int16(cc)
		fmt.Printf("%d", cc)
	}
}

func createMapping(in drivers.In, lights []huego.Light) {
	scanner := bufio.NewScanner(os.Stdin)
	for _, light := range lights {
		fmt.Printf("Touch the knob that will control %s or press ENTER to skip: ", light.Name)
		lastTouchedControl = -1
		stopMidi, err := midi.ListenTo(in, detectControl)
		if err != nil {
			log.Fatal(err)
		}
		scanner.Scan()
		stopMidi()
		if lastTouchedControl > -1 {
			mapping[uint8(lastTouchedControl)] = append(mapping[uint8(lastTouchedControl)], light.ID)
			fmt.Printf("Will use %d\n", lastTouchedControl)
		}
	}
}
