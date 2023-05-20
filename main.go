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

func connectBT(deviceName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := bt.Connect(ctx, deviceName); err != nil {
		panic(err)
	}
}

func startBTNotifier(ServiceUUID [16]byte, CharacteristicUUID [16]byte) {
	if err := bt.StartNotifier(ServiceUUID, CharacteristicUUID); err != nil {
		panic(err)
	}
}

func startBTWriter(ServiceUUID [16]byte, CharacteristicUUID [16]byte) {
	if err := bt.StartWriter(ServiceUUID, CharacteristicUUID); err != nil {
		panic(err)
	}
}

func startParser(m multimeter.Multimeter, printArray bool) {
	for bt.Connected() {
		val, unit, flags := m.ProccessArray(<-bt.ChNotifier, printArray)
		if unit != "" && !printArray {
			log.Printf("%v %v [%v]\n", val, unit, strings.Join(flags, ", "))
		}
	}
}

func write(m multimeter.MultimeterButtons) {
	bt.ChWrite <- m.Select()
	time.Sleep(2 * time.Second)

	bt.ChWrite <- m.Range()
	time.Sleep(2 * time.Second)

	bt.ChWrite <- m.Range()
	time.Sleep(2 * time.Second)

	bt.ChWrite <- m.Auto()
	time.Sleep(2 * time.Second)

	bt.ChWrite <- m.Light()
	time.Sleep(2 * time.Second)
}

func fs9721lp3(printArray bool) {
	connectBT(fs9721.DeviceName)
	startBTNotifier(fs9721.ServiceUUID, fs9721.CharacteristicNotifyUUID)
	go startParser(&fs9721.Fs9721{}, printArray)
}

func ow18e(printArray bool) {
	connectBT(owon.DeviceName)
	startBTNotifier(owon.ServiceUUID, owon.CharacteristicNotifyUUID)
	go startParser(&owon.OW18E{}, printArray)

	// Example to send command to multimeter
	// startBTWriter(owon.ServiceUUID, owon.CharacteristicWriteUUID)
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

	device := strings.TrimSpace(strings.ToLower(os.Args[1]))
	switch device {
	case "fs9721":
		fs9721lp3(printArray)
	case "ow18e":
		ow18e(printArray)
	default:
		log.Println("Invalid argument!")
		log.Println("Valid arguments: \"fs9721\" or \"ow18e\"")
	}

	if bt.Connected() {
		log.Printf("[%v] Ready!\n", device)
		log.Println("Press <ENTER> to exit")
		bufio.NewScanner(os.Stdin).Scan()
		bt.Disconnect()
		log.Println("Disconnected!!")
	}
}
