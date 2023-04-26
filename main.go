package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/raulaguila/go-multimeter/bluetooth"
	"github.com/raulaguila/go-multimeter/multimeter"
	"github.com/raulaguila/go-multimeter/multimeter/fs9721"
	"github.com/raulaguila/go-multimeter/multimeter/owon"
)

var bt = bluetooth.Bluetooth{}

func startBT(deviceName string, ServiceUUID [16]byte, CharacteristicUUID [16]byte) {
	if err := bt.Enable(); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := bt.Connect(ctx, deviceName); err != nil {
		panic(err)
	}

	if err := bt.Read(ServiceUUID, CharacteristicUUID); err != nil {
		panic(err)
	}
}

func startParser(m multimeter.Multimeter) {
	go func(printArray bool) {
		for bt.Connected() {
			if printArray {
				go m.ProccessArray(<-bt.ChReceived, printArray)
				continue
			}

			val, unit, flags := m.ProccessArray(<-bt.ChReceived, printArray)
			if unit != "" {
				log.Printf("%v %v %v\n", val, unit, flags)
			}
		}
	}(false)
}

func fs9721lp3() {
	startBT(fs9721.DeviceName, fs9721.ServiceUUID, fs9721.CharacteristicUUID)
	log.Println("[FS9721_LP3] Ready!")
	startParser(&fs9721.Fs9721{})
}

func ow18e() {
	startBT(owon.DeviceName, owon.ServiceUUID, owon.CharacteristicUUID)
	log.Println("[OW18E] Ready!")
	startParser(&owon.OW18E{})
}

func main() {
	if len(os.Args) == 1 {
		log.Println("Pass argument: \"fs9721\" or \"ow18e\"")
		return
	}

	switch strings.TrimSpace(strings.ToLower(os.Args[1])) {
	case "fs9721":
		fs9721lp3()
	case "ow18e":
		ow18e()
	default:
		log.Println("Invalid argument! valid argument: \"fs9721\" or \"ow18e\"")
	}

	if bt.Connected() {
		log.Println("Press <ENTER> to exit")
		bufio.NewScanner(os.Stdin).Scan()
		bt.Disconnect()
		log.Println("Desconnected!!")
	}
}
