package main

import (
	"AP-90EPD/bluetooth"
	"AP-90EPD/multimeter"
	"log"
	"time"
)

var (
	deviceName = "FS9721-LP3"
	bt         = bluetooth.Bluetooth{}
	multi      = multimeter.Ap90epd{}
)

func main() {
	if err := bt.Enable(); err != nil {
		panic(err)
	}

	log.Println("connecting...")
	if err := bt.Connect(deviceName); err != nil {
		panic(err)
	}

	log.Println("preparing to read...")
	if err := bt.Read(multimeter.ServiceUUID, multimeter.CharacteristicUUID); err != nil {
		panic(err)
	}

	log.Println("ready!")
	go func() {
		for bt.Connected() {
			val, unit, flags := multi.AddToByteArray(<-bt.ChReceived)
			if unit != "" {
				log.Printf("%v %v %v\n", val, unit, flags)
			}
		}
	}()

	time.Sleep(5 * time.Second)
	log.Println("disconnecting...")
	bt.Disconnect()
}
