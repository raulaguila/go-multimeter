package owon

import (
	"fmt"
	"log"
)

const (
	dc         string = "00"
	ac         string = "01"
	diod       string = "10"
	cont       string = "11"
	nano       string = "001"
	micro      string = "010"
	mili       string = "011"
	one        string = "100"
	kilo       string = "101"
	mega       string = "110"
	voltage    string = "0000"
	resistance string = "0001"
	continuity string = "0010"
	ncv        string = "0011"
	auto       string = "10"
	relative   string = "01"
)

type OW18E struct {
	bytearray []byte
	binarray  []string
}

func (m *OW18E) getBinArray() []string {
	str := []string{}
	for _, b := range m.bytearray {
		str = append(str, fmt.Sprintf("%08b", b))
	}
	return str
}

func (m *OW18E) ProccessArray(bytearray []byte, printArray bool) (float64, string, []string) {
	m.bytearray = bytearray
	m.binarray = m.getBinArray()

	value := m.extractValue()
	unit := m.extractUnit()
	flags := m.extractFlags()

	if printArray {
		log.Printf("%v <-> %v <-> %v %v %v\n", m.binarray, m.bytearray, value, unit, flags)
	}

	return value, unit, flags
}

func (m *OW18E) extractValue() (ret float64) {
	div := 10
	mrange := m.binarray[0][5:]
	switch mrange {
	case "100": // range 2
		div = 10000
	case "011": // range 20
		div = 1000
	case "010": // range 200
		div = 100
	case "001": // range 2000
		div = 10
	case "111": // L
		return
	default:
		log.Printf("\tRange not tracked: %v\n", mrange)
	}

	if m.bytearray[5] < 128 {
		ret = ((float64(m.bytearray[5]) * 256) + float64(m.bytearray[4])) / float64(div)
	} else {
		ret = ((((float64(m.bytearray[5]) - 128.0) * 256.0) + float64(m.bytearray[4])) / float64(div)) * -1
	}

	return
}

func (m *OW18E) extractUnit() (unit string) {
	firstByte := m.binarray[0]
	secondByte := m.binarray[1]

	unitRepresentation := firstByte[2:5]
	finalFunction := firstByte[:2] + secondByte[4:]

	arrUnits := [][2]interface{}{
		{unitRepresentation == nano, "n"},        // nano
		{unitRepresentation == micro, "µ"},       // micro
		{unitRepresentation == mili, "m"},        // mili
		{unitRepresentation == one, ""},          // 1
		{unitRepresentation == kilo, "k"},        // kilo
		{unitRepresentation == mega, "M"},        // Mega
		{finalFunction == dc+continuity, "ºC"},   // Temp celsius
		{finalFunction == dc+voltage, "V"},       // DC Voltage Measure
		{finalFunction == dc+resistance, "Ω"},    // Resistance Measure
		{finalFunction == ac+continuity, "ºF"},   // Temp fahrenheit
		{finalFunction == ac+resistance, "F"},    // Capacitance Measure
		{finalFunction == ac+voltage, "V"},       // AC Voltage Measure
		{finalFunction == diod+continuity, "V"},  // Diode test
		{finalFunction == diod+resistance, "Hz"}, // Frequence
		{finalFunction == diod+voltage, "A"},     // Current Measure
		{finalFunction == cont+continuity, "Ω"},  // Continuity test
		{finalFunction == cont+resistance, "%"},  // Percentage
	}

	for _, item := range arrUnits {
		if item[0].(bool) {
			unit += item[1].(string)
		}
	}

	return
}

func (m *OW18E) extractFlags() (flags []string) {
	firstByte := m.binarray[0]
	secondByte := m.binarray[1]
	thirdByte := m.binarray[2]

	funcRepresentation := firstByte[0:2]
	finalFunction := funcRepresentation + secondByte[4:]

	arrFlags := [][2]interface{}{
		{funcRepresentation == dc, "DC"},                      // DC Voltage Measure
		{funcRepresentation == ac, "AC"},                      // AC Voltage Measure
		{firstByte[5:] == "111", "L"},                         // L
		{thirdByte[5:7] == auto, "Auto Mode"},                 // Auto Mode
		{thirdByte[5:7] == relative, "Relative Mode"},         // Relative Mode
		{thirdByte[7:] == "1", "Hold"},                        // Hold
		{finalFunction == dc+continuity, "Temp celsius"},      // Temp celsius
		{finalFunction == ac+continuity, "Temp fahrenheit"},   // Temp fahrenheit
		{finalFunction == ac+resistance, "Capacity"},          // Capacitance Measure
		{finalFunction == ac+ncv, "NCV Measure"},              // NCV Measure
		{finalFunction == diod+continuity, "Diode test"},      // Diode test
		{finalFunction == cont+continuity, "Continuity test"}, // Continuity test
		{finalFunction == cont+resistance, "Percentage"},      // Percentage
	}

	for _, flag := range arrFlags {
		if flag[0].(bool) {
			flags = append(flags, flag[1].(string))
		}
	}

	return
}
