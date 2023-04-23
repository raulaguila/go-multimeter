package bluetooth

import (
	"context"
	"errors"
	"fmt"
	"log"

	"tinygo.org/x/bluetooth"
)

type Bluetooth struct {
	ChReceived chan []byte
	connected  bool
	adapter    *bluetooth.Adapter
	chScan     chan bluetooth.ScanResult
	device     *bluetooth.Device
	srvcs      []bluetooth.DeviceService
	chars      []bluetooth.DeviceCharacteristic
}

func (b *Bluetooth) Enable() error {
	b.adapter = bluetooth.DefaultAdapter
	return b.adapter.Enable()
}

func (b *Bluetooth) find(deviceName string) {
	b.chScan = make(chan bluetooth.ScanResult, 1)
	go b.adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if device.LocalName() == deviceName {
			adapter.StopScan()
			b.chScan <- device
		}
	})
}

func (b *Bluetooth) Connected() bool {
	return b.connected
}

func (b *Bluetooth) Disconnect() (err error) {
	err = nil

	if b.connected {
		b.connected = false
		err = b.device.Disconnect()
	}

	return
}

func (b *Bluetooth) Connect(ctx context.Context, deviceName string) (err error) {
	err = nil
	if !b.connected {
		b.find(deviceName)

		select {
		case <-ctx.Done():
			b.adapter.StopScan()
			err = ctx.Err()
		case device := <-b.chScan:
			b.device, err = b.adapter.Connect(device.Address, bluetooth.ConnectionParams{})
			b.connected = err == nil
		}
	}

	return
}

func (b *Bluetooth) StopScan() error {
	return b.adapter.StopScan()
}

func (b *Bluetooth) ScanDevices() {
	go b.adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		log.Println(result, result.LocalName())
	})
}

func (b *Bluetooth) ListUUIDs() error {
	srvcs, err := b.device.DiscoverServices(nil)
	if err != nil {
		return err
	}

	for _, srvc := range srvcs {
		fmt.Printf("Service UUID: %v\n", srvc.UUID())
		char, err := srvc.DiscoverCharacteristics(nil)
		if err != nil {
			return err
		}

		for _, c := range char {
			fmt.Printf(" - Characteristic UUID: %v\n", c.UUID())
		}

		fmt.Println("")
	}

	return nil
}

func (b *Bluetooth) Read(ServiceUUID [16]byte, CharacteristicUUID [16]byte) (err error) {
	err = nil

	if !b.connected {
		err = errors.New("bluetooth not connected")
		return
	}

	b.srvcs, err = b.device.DiscoverServices([]bluetooth.UUID{bluetooth.NewUUID(ServiceUUID)})
	if err != nil {
		return
	}

	if len(b.srvcs) == 0 {
		err = errors.New("could not find service")
		return
	}

	b.chars, err = b.srvcs[0].DiscoverCharacteristics([]bluetooth.UUID{bluetooth.NewUUID(CharacteristicUUID)})
	if err != nil {
		return
	}

	if len(b.chars) == 0 {
		err = errors.New("could not find characteristic")
		return
	}

	b.ChReceived = make(chan []byte, 1)
	b.chars[0].EnableNotifications(func(byteArray []byte) {
		b.ChReceived <- byteArray
	})

	return
}
