# Multimeter bluetooth communication protocol

Get BLE Multimeter data to PC.

Tested with Owon [OW18E](/multimeter/owon/README.md) multimeter and with [FS9721-LP3](/multimeter/fs9721/README.md) based bluetooth multimeters.

Tested with Linux using golang 1.20.3.

## Quickstart

* Download dependencies: `go mod download`.
* Switch on your multimeter in close proximity (~4m) to the PC.
* Run `go run main.go <ow18e or fs9721>`
* Watch terminal.
* Press **ENTER** to disconnect bluetooth and close the program.
