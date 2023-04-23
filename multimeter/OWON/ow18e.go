package owon

import (
	"fmt"
	"log"
)

var (
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
	auto       string = "0100"
)

type OW18E struct {
	Bytearray []string
}

func (m *OW18E) ShowByteArray(bytearray []byte) {
	str := []string{}
	for _, b := range bytearray {
		str = append(str, m.toBintString(b))
	}
	log.Printf("%v <-> %v \n", str, bytearray)
}

func (m *OW18E) ProccessArray(bytearray []byte) (float64, string, []string) {
	// m.ShowByteArray(bytearray)

	value := m.extractValue(bytearray)
	unit := m.extractUnit(bytearray)
	flags := m.extractFlags(bytearray)

	return value, unit, flags
}

func (m *OW18E) toBintString(b byte) string {
	return fmt.Sprintf("%08b", b)
}

func (m *OW18E) extractValue(bytearray []byte) (ret float64) {
	str := m.toBintString(bytearray[0])
	div := 10
	switch str[5:] {
	case "100":
		div = 10000
	case "011":
		div = 1000
	case "010":
		div = 100
	case "001":
		div = 10
	case "111":
		ret = 0
		return
	}

	if bytearray[5] < 128 {
		ret = ((float64(bytearray[5]) * 256) + float64(bytearray[4])) / float64(div)
	} else {
		ret = ((((float64(bytearray[5]) - 128.0) * 256.0) + float64(bytearray[4])) / float64(div)) * -1
	}

	return
}

func (m *OW18E) extractUnit(bytearray []byte) (unit string) {
	firstByte := m.toBintString(bytearray[0])
	secondByte := m.toBintString(bytearray[1])

	unitRepresentation := firstByte[2:5]
	finalFunction := firstByte[0:2] + secondByte[4:]

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
	}

	for _, item := range arrUnits {
		if item[0].(bool) {
			unit += item[1].(string)
		}
	}

	return
}

func (m *OW18E) extractFlags(bytearray []byte) (flags []string) {
	firstByte := m.toBintString(bytearray[0])
	secondByte := m.toBintString(bytearray[1])
	thirdByte := m.toBintString(bytearray[2])

	funcRepresentation := firstByte[0:2]
	finalFunction := funcRepresentation + secondByte[4:]

	arrFlags := [][2]interface{}{
		{funcRepresentation == dc, "DC"},                      // DC Voltage Measure
		{funcRepresentation == ac, "AC"},                      // AC Voltage Measure
		{thirdByte[4:] == auto, "Auto"},                       // Auto range
		{finalFunction == dc+continuity, "Temp celsius"},      // Temp celsius
		{finalFunction == ac+continuity, "Temp fahrenheit"},   // Temp fahrenheit
		{finalFunction == diod+continuity, "Diode test"},      // Diode test
		{finalFunction == cont+continuity, "Continuity test"}, // Continuity test
	}

	for _, flag := range arrFlags {
		if flag[0].(bool) {
			flags = append(flags, flag[1].(string))
		}
	}

	return
}
