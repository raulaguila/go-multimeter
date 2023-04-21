package bluetooth

import (
	"errors"

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
	go b.adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.LocalName() == deviceName {
			// log.Println(result, result.LocalName())
			adapter.StopScan()
			b.chScan <- result
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

func (b *Bluetooth) Connect(deviceName string) (err error) {
	err = nil
	if !b.connected {
		b.find(deviceName)
		result := <-b.chScan
		b.device, err = b.adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		b.connected = err == nil
	}

	return
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
		err = errors.New("could not find heart rate service")
		return
	}

	b.chars, err = b.srvcs[0].DiscoverCharacteristics([]bluetooth.UUID{bluetooth.NewUUID(CharacteristicUUID)})
	if err != nil {
		return
	}

	if len(b.chars) == 0 {
		err = errors.New("could not find heart rate characteristic")
		return
	}

	b.ChReceived = make(chan []byte, 1)
	b.chars[0].EnableNotifications(func(byteArray []byte) {
		b.ChReceived <- byteArray
	})

	return
}
