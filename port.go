package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

func readPort() drivers.In {
	var port drivers.In
	for port == nil {
		fmt.Printf("MIDI IN Ports\n")
		for _, port = range midi.GetInPorts() {

			fmt.Printf("%d %s\n", port.Number(), port)
		}
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
