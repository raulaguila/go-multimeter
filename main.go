package main

import (
	"log"
	"time"

	"github.com/raulaguila/go-multimeter/bluetooth"
	fs9721lp3 "github.com/raulaguila/go-multimeter/multimeter/FS9721-LP3"
)

var (
	bt              = bluetooth.Bluetooth{}
	multi_fs9721lp3 = fs9721lp3.Fs9721lp3{}
)

func Fs9721lp3() {
	if err := bt.Enable(); err != nil {
		panic(err)
	}

	log.Println("connecting...")
	if err := bt.Connect(fs9721lp3.DeviceName); err != nil {
		panic(err)
	}

	log.Println("preparing to read...")
	if err := bt.Read(fs9721lp3.ServiceUUID, fs9721lp3.CharacteristicUUID); err != nil {
		panic(err)
	}

	log.Println("ready!")
	go func() {
		for bt.Connected() {
			val, unit, flags := multi_fs9721lp3.AddToByteArray(<-bt.ChReceived)
			if unit != "" {
				log.Printf("%v %v %v\n", val, unit, flags)
			}
		}
	}()

	time.Sleep(50 * time.Second)
	log.Println("disconnecting...")
	bt.Disconnect()
}

func main() {
	Fs9721lp3()
}
