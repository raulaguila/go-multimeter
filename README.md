# Go - BLE Multimeter

Get BLE Multimeter data to PC. With approx 2 samples per second (in DC Volt mode).

Works with FS9721-LP3 based Bluetooth Multimeters, including:

* AOPUTTRIVER AP-90EPD
* Infurider YF-90EPD
* HoldPeak HP-90EPD

![](/img/example.png)

Tested with Linux using golang 1.20.3.

## Credits

This work derives heavily from the main source:

* [pyBleMultimeter](https://github.com/mechaot/pyBleMultimeter)

## Quickstart

* Install dependencies via: `go mod download`.
* Switch on your multimeter in close proximity (~5m) to the PC.
* Run `go run main.go`
* Watch terminal.