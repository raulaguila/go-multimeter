package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/raulaguila/go-multimeter/bluetooth"
	fs9721 "github.com/raulaguila/go-multimeter/multimeter/FS9721-LP3"
	owon "github.com/raulaguila/go-multimeter/multimeter/OWON"
)

var bt = bluetooth.Bluetooth{}

func startBT(deviceName string, ServiceUUID [16]byte, CharacteristicUUID [16]byte) {
	if err := bt.Enable(); err != nil {
		panic(err)
	}

	if err := bt.Connect(deviceName); err != nil {
		panic(err)
	}

	if err := bt.Read(ServiceUUID, CharacteristicUUID); err != nil {
		panic(err)
	}
}

func Fs9721lp3() {
	startBT(fs9721.DeviceName, fs9721.ServiceUUID, fs9721.CharacteristicUUID)

	go func() {
		fs9721 := fs9721.Fs9721lp3{}
		log.Println("[FS9721_LP3] Ready!")
		for bt.Connected() {
			val, unit, flags := fs9721.AddToByteArray(<-bt.ChReceived)
			if unit != "" {
				log.Printf("%v %v %v\n", val, unit, flags)
			}
		}
	}()
}

func Ow18e() {
	startBT(owon.DeviceName, owon.ServiceUUID, owon.CharacteristicUUID)

	go func() {
		ow18 := owon.OW18E{}
		log.Println("[OW18E] Ready!")
		for bt.Connected() {
			val, unit, flags := ow18.ProccessArray(<-bt.ChReceived)
			log.Printf("%v %v %v\n", val, unit, flags)
		}
	}()
}

func main() {
	switch strings.TrimSpace(strings.ToLower(os.Args[1])) {
	case "fs9721":
		Fs9721lp3()
	case "ow18e":
		Ow18e()
	}

	if bt.Connected() {
		log.Println("Press <ENTER> to exit")
		bufio.NewScanner(os.Stdin).Scan()
		log.Println("disconnecting...")
		bt.Disconnect()
	}
}
