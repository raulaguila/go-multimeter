package bluetooth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/raulaguila/go-multimeter/bluetooth/enum"
	"tinygo.org/x/bluetooth"
)

type Bluetooth struct {
	ChReceived      chan []byte
	ChWrite         chan []byte
	chScan          chan bluetooth.ScanResult
	connected       bool
	adapter         *bluetooth.Adapter
	device          *bluetooth.Device
	characteristics [2]bluetooth.DeviceCharacteristic
}

func (b *Bluetooth) enable() error {
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

func (b *Bluetooth) StopScan() error {
	return b.adapter.StopScan()
}

func (b *Bluetooth) Connected() bool {
	return b.connected
}

func (b *Bluetooth) Disconnect() error {
	if b.connected {
		b.connected = false
		time.Sleep(100 * time.Millisecond)
		return b.device.Disconnect()
	}

	return nil
}

func (b *Bluetooth) Connect(ctx context.Context, deviceName string) (err error) {
	err = nil
	if !b.connected {
		if err = b.enable(); err != nil {
			return
		}

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

func (b *Bluetooth) ScanDevices() {
	go b.adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		log.Println(result, result.LocalName())
	})
}

func (b *Bluetooth) ListUUIDs() error {
	services, err := b.device.DiscoverServices(nil)
	if err != nil {
		return err
	}

	for _, service := range services {
		fmt.Printf("Service UUID: %v\n", service.UUID())
		characteristics, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			return err
		}

		for _, characteristic := range characteristics {
			fmt.Printf(" - Characteristic UUID: %v\n", characteristic.UUID())
			fmt.Printf("   - Properties: %v\n\n", characteristic.Properties())
		}

		fmt.Println("")
	}

	return nil
}

func (b *Bluetooth) getCharacteristic(ServiceUUID [16]byte, CharacteristicUUID [16]byte) (*bluetooth.DeviceCharacteristic, error) {
	if !b.connected {
		return nil, errors.New("bluetooth not connected")
	}

	services, err := b.device.DiscoverServices([]bluetooth.UUID{bluetooth.NewUUID(ServiceUUID)})
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, errors.New("could not find service")
	}

	characteristics, err := services[0].DiscoverCharacteristics([]bluetooth.UUID{bluetooth.NewUUID(CharacteristicUUID)})
	if err != nil {
		return nil, err
	}

	if len(characteristics) == 0 {
		return nil, errors.New("could not find characteristic")
	}

	return &characteristics[0], nil
}

func (b *Bluetooth) StartNotifier(ServiceUUID [16]byte, CharacteristicUUID [16]byte) error {
	characteristic, err := b.getCharacteristic(ServiceUUID, CharacteristicUUID)
	if err != nil {
		return err
	}

	b.characteristics[enum.Reader] = *characteristic
	b.ChReceived = make(chan []byte, 1)
	b.characteristics[enum.Reader].EnableNotifications(func(byteArray []byte) {
		b.ChReceived <- byteArray
	})

	return nil
}

func (b *Bluetooth) StartWriter(ServiceUUID [16]byte, CharacteristicUUID [16]byte) error {
	characteristic, err := b.getCharacteristic(ServiceUUID, CharacteristicUUID)
	if err != nil {
		return err
	}

	b.characteristics[enum.Writer] = *characteristic
	b.ChWrite = make(chan []byte, 1)
	go func() {
		for b.connected {
			b.characteristics[enum.Writer].WriteWithoutResponse(<-b.ChWrite)
		}
	}()

	return nil
}
