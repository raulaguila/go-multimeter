# Multimeter bluetooth communication protocol

Get BLE Multimeter data to PC.

Tested with Owon [OW18E](/multimeter/owon/README.md) multimeter and with [FS9721-LP3](/multimeter/fs9721/README.md) based bluetooth multimeters.

Tested with Linux using golang 1.20.

## Requirements

[![Golang](https://img.shields.io/badge/Golang-v1.20-%2300ADD8.svg?style=flat&logo=go&logoColor=2300ADD8&labelColor=0D1117)](https://go.dev/)

## Quickstart

* Download dependencies: `go mod download`.
* Switch on your multimeter in close proximity (~4m) to the PC.
* Run: `go run main.go <ow18e or fs9721>`
* You can also pass the **true** flag after "ow18e" or "fs9721" to print the received array and the respective value.
* Watch terminal.
* Press **ENTER** to disconnect bluetooth and close the program.
