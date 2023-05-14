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

func startBTRead(deviceName string, ServiceUUID [16]byte, CharacteristicUUID [16]byte) {
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

func startBTWrite(deviceName string, ServiceUUID [16]byte, CharacteristicUUID [16]byte) {
	if err := bt.Enable(); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := bt.Connect(ctx, deviceName); err != nil {
		panic(err)
	}

	if err := bt.Write(ServiceUUID, CharacteristicUUID); err != nil {
		panic(err)
	}
}

func startParser(m multimeter.Multimeter, printArray bool) {
	go func(printArray bool) {
		for bt.Connected() {
			val, unit, flags := m.ProccessArray(<-bt.ChReceived, printArray)
			if unit != "" && !printArray {
				log.Printf("%v %v [%v]\n", val, unit, strings.Join(flags, ", "))
			}
		}
	}(printArray)
}

func write(m multimeter.MultimeterButtons) {
	bt.ChWrite <- m.Select()
	time.Sleep(2 * time.Second)

	bt.ChWrite <- m.Auto()
	time.Sleep(2 * time.Second)

	bt.ChWrite <- m.Range()
	time.Sleep(2 * time.Second)

	bt.ChWrite <- m.Range()
	time.Sleep(2 * time.Second)

	bt.ChWrite <- m.Range()
	time.Sleep(2 * time.Second)

	bt.ChWrite <- m.Range()
	time.Sleep(2 * time.Second)
}

func fs9721lp3(printArray bool) {
	startBTRead(fs9721.DeviceName, fs9721.ServiceUUID, fs9721.CharacteristicNotifyUUID)
	log.Println("[FS9721_LP3] Ready!")
	startParser(&fs9721.Fs9721{}, printArray)
}

func ow18e(printArray bool) {
	startBTRead(owon.DeviceName, owon.ServiceUUID, owon.CharacteristicNotifyUUID)
	log.Println("[OW18E] Ready!")
	startParser(&owon.OW18E{}, printArray)

	// Example to send command to multimeter
	// startBTWrite(owon.DeviceName, owon.ServiceUUID, owon.CharacteristicWriteUUID)
	// write(&owon.OW18E{})
}

func main() {
	if len(os.Args) == 1 {
		log.Println("Required argument: \"fs9721\" or \"ow18e\"")
		return
	}

	printArray := false
	if len(os.Args) > 2 {
		printArray = strings.ToLower(os.Args[2]) == "true"
	}

	switch strings.TrimSpace(strings.ToLower(os.Args[1])) {
	case "fs9721":
		fs9721lp3(printArray)
	case "ow18e":
		ow18e(printArray)
	default:
		log.Println("Invalid argument! valid arguments: \"fs9721\" or \"ow18e\"")
	}

	if bt.Connected() {
		log.Println("Press <ENTER> to exit")
		bufio.NewScanner(os.Stdin).Scan()
		bt.Disconnect()
		log.Println("Disconnected!!")
	}
}
