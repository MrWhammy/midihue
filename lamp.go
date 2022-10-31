package main

import (
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/amimof/huego"
)

type lamp struct {
	name    string
	channel chan uint8
	val     uint32
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
		fmt.Println("Setting ", lamp.name, " to ", newVal)
		lamp.light.Bri(newVal)
		fmt.Println("Done ", lamp.name, " to ", newVal)
		atomic.StoreUint32(&lamp.val, uint32(newVal))
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
	return lamp.name + ": " + strconv.FormatUint(uint64(lamp.val), 10)
}
