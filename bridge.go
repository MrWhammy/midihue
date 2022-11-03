package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/amimof/huego"
)

func findBridge() *huego.Bridge {
	conf, err := loadConf()
	if err == nil {
		return huego.New(conf.Ip, conf.User)
	}
	log.Println(err)
	bridge := discoverBridge()
	storeConf(bridge)
	return bridge
}

func discoverBridge() *huego.Bridge {
	bridge, err := huego.Discover()
	if err != nil {
		log.Fatal(err)
	}
	user, err := bridge.CreateUser("midihue") // Link button needs to be pressed
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("New user: " + user)
	}
	return bridge.Login(user)
}

type Conf struct {
	Ip   string
	User string
}

const confFile = "bridge.json"

func storeConf(bridge *huego.Bridge) {
	conf := Conf{Ip: bridge.Host, User: bridge.User}
	serConf, err := json.Marshal(conf)
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile(confFile, serConf, 0644)
}

func loadConf() (Conf, error) {
	var conf Conf
	serConf, err := os.ReadFile(confFile)
	if err != nil {
		return conf, err
	}
	json.Unmarshal(serConf, &conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}
