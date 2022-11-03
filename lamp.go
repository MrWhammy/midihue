package main

import (
	"github.com/amimof/huego"
)

type lamp struct {
	name    string
	channel chan uint8
	light   *huego.Light
}

func Lamp(light *huego.Light) *lamp {
	var lamp lamp
	lamp.name = light.Name
	lamp.channel = make(chan uint8, 256)
	lamp.light = light
	return &lamp
}

func (lamp *lamp) SetVal(newVal uint8) {
	lamp.channel <- newVal
}

func (lamp *lamp) Run() {
	newVal, ok := lamp.ReadValue()
	for ok {
		if newVal > 0 {
			lamp.light.Bri(newVal)
		} else {
			lamp.light.Off()
		}
		newVal, ok = lamp.ReadValue()
	}
}

func (lamp *lamp) ReadValue() (uint8, bool) {
	for len(lamp.channel) > 1 {
		_, ok := <-lamp.channel
		if !ok {
			return 0, false
		}
	}
	lastVal, ok := <-lamp.channel
	return lastVal, ok
}

func (lamp *lamp) Close() {
	close(lamp.channel)
}

func (lamp *lamp) String() string {
	return lamp.name
}
