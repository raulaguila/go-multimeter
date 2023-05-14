package owon

import (
	"fmt"
	"log"
	"strings"
)

const (
	// First 2 bits of 1st byte
	dc   string = "00"
	ac   string = "01"
	diod string = "10"
	cont string = "11"

	// Last 2 bits of 2nd byte
	voltage    string = "00"
	resistance string = "01"
	continuity string = "10"
	ncv        string = "11"
)

type OW18E struct {
	bytearray     []byte
	binarray      []string
	function      string
	finalfunction string
	unity         string
	mrange        string
}

func (m *OW18E) getBinArray() []string {
	str := []string{}
	for _, b := range m.bytearray {
		str = append(str, fmt.Sprintf("%08b", b))
	}
	return str
}

func (m *OW18E) ProccessArray(bytearray []byte, printArray bool) (value float64, unit string, flags []string) {
	m.bytearray = bytearray                          // Save byte array
	m.binarray = m.getBinArray()                     // Convert byte array to bits string
	m.unity = m.binarray[0][2:5]                     // Get unit
	m.function = m.binarray[0][:2]                   // Get function
	m.finalfunction = m.function + m.binarray[1][6:] // Get final function
	m.mrange = m.binarray[0][5:]                     // Get range

	value = m.extractValue()
	unit = m.extractUnit()
	flags = m.extractFlags()

	if printArray {
		log.Printf("%v <-> %v <-> %v %v [%v]\n", m.binarray, m.bytearray, value, unit, strings.Join(flags, ", "))
	}

	return value, unit, flags
}

func (m *OW18E) calcValue(byte5 byte, byte4 byte, div float64, negative bool) (ret float64) {
	ret = float64(byte5)
	if negative {
		ret -= 128
	}

	ret = ((ret * 256) + float64(byte4)) / div

	if negative {
		ret *= -1
	}

	return
}

func (m *OW18E) extractValue() (ret float64) {
	switch m.mrange {
	case "100": // range 2
		ret = m.calcValue(m.bytearray[5], m.bytearray[4], 10000, m.bytearray[5] >= 128)
	case "011": // range 20
		ret = m.calcValue(m.bytearray[5], m.bytearray[4], 1000, m.bytearray[5] >= 128)
	case "010": // range 200
		ret = m.calcValue(m.bytearray[5], m.bytearray[4], 100, m.bytearray[5] >= 128)
	case "001": // range 2000
		ret = m.calcValue(m.bytearray[5], m.bytearray[4], 10, m.bytearray[5] >= 128)
	case "000": // NCV
		ret = float64(m.bytearray[4])
		return
	case "111": // L
		return
	default:
		log.Printf("\tRange not tracked: %v\n", m.mrange)
		return
	}

	return
}

func (m *OW18E) extractUnit() (unit string) {
	arrUnits := [][2]interface{}{
		{m.unity == "001", "n"},                    // nano
		{m.unity == "010", "µ"},                    // micro
		{m.unity == "011", "m"},                    // mili
		{m.unity == "100", ""},                     // 1
		{m.unity == "101", "k"},                    // kilo
		{m.unity == "110", "M"},                    // Mega
		{m.finalfunction == dc+continuity, "ºC"},   // Temp celsius
		{m.finalfunction == dc+voltage, "V"},       // DC Voltage Measure
		{m.finalfunction == dc+resistance, "Ω"},    // Resistance Measure
		{m.finalfunction == ac+continuity, "ºF"},   // Temp fahrenheit
		{m.finalfunction == ac+resistance, "F"},    // Capacitance Measure
		{m.finalfunction == ac+voltage, "V"},       // AC Voltage Measure
		{m.finalfunction == ac+ncv, "NVC"},         // NVC Measure
		{m.finalfunction == diod+continuity, "V"},  // Diode test
		{m.finalfunction == diod+resistance, "Hz"}, // Frequence
		{m.finalfunction == diod+voltage, "A"},     // Current Measure
		{m.finalfunction == cont+continuity, "Ω"},  // Continuity test
		{m.finalfunction == cont+resistance, "%"},  // Percentage
	}

	for _, item := range arrUnits {
		if item[0].(bool) {
			unit += item[1].(string)
		}
	}

	return
}

func (m *OW18E) extractFlags() (flags []string) {
	arrFlags := [][2]interface{}{
		{m.function == dc, "DC"},                                // DC Voltage Measure
		{m.function == ac, "AC"},                                // AC Voltage Measure
		{m.mrange == "111", "L"},                                // L
		{m.finalfunction == dc+continuity, "Temp celsius"},      // Temp celsius
		{m.finalfunction == ac+continuity, "Temp fahrenheit"},   // Temp fahrenheit
		{m.finalfunction == ac+resistance, "Capacity"},          // Capacitance Measure
		{m.finalfunction == ac+ncv, "NCV Measure"},              // NCV Measure
		{m.finalfunction == diod+continuity, "Diode test"},      // Diode test
		{m.finalfunction == cont+continuity, "Continuity test"}, // Continuity test
		{m.finalfunction == cont+resistance, "Percentage"},      // Percentage
		{m.binarray[2][4] == '1', "Low Battery"},                // Low Battery
		{m.binarray[2][5] == '1', "Auto Range"},                 // Auto Range
		{m.binarray[2][6] == '1', "Relative Mode"},              // Relative Mode
		{m.binarray[2][7] == '1', "Hold"},                       // Hold
	}

	for _, flag := range arrFlags {
		if flag[0].(bool) {
			flags = append(flags, flag[1].(string))
		}
	}

	return
}

// Buttons command
func (m *OW18E) Select() []byte {
	return []byte{1, 0}
}

func (m *OW18E) Auto() []byte {
	return []byte{2, 0}
}

func (m *OW18E) Range() []byte {
	return []byte{2, 1}
}

func (m *OW18E) Light() []byte {
	return []byte{3, 0}
}

func (m *OW18E) Hold() []byte {
	return []byte{3, 1}
}

func (m *OW18E) Relative() []byte {
	return []byte{4, 0}
}
