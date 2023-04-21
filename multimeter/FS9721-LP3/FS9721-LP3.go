package fs9721lp3

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	DeviceName         = "FS9721-LP3"
	ServiceUUID        = [16]byte{0x00, 0x00, 0xff, 0xb0, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0x80, 0x5f, 0x9b, 0x34, 0xfb}
	CharacteristicUUID = [16]byte{0x00, 0x00, 0xff, 0xb2, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0x80, 0x5f, 0x9b, 0x34, 0xfb}
)

type Fs9721lp3 struct {
	Bytearray []string
}

func (m *Fs9721lp3) AddToByteArray(bytearray []byte) (float64, string, []string) {
	if len(bytearray) == 8 {
		m.Bytearray = m.Bytearray[:0]
	}

	if len(bytearray) == 8 || len(bytearray) == 6 {
		for _, b := range bytearray {
			aux := fmt.Sprintf("%08b", b)
			m.Bytearray = append(m.Bytearray, aux[len(aux)-4:])
		}

		if len(m.Bytearray) == 14 {
			return m.proccessArray(m.Bytearray)
		}
	}

	return 0, "", []string{}
}

func (m *Fs9721lp3) proccessArray(bytearray []string) (float64, string, []string) {
	str := strings.Join(bytearray, "")
	if len(str) != 56 {
		return 0, "", []string{}
	}

	value := m.extractValue(str)
	unit := m.extractUnit(str)
	flags := m.extractFlags(str)

	// log.Printf("%v %v %v\n", value, unit, flags)

	return value, unit, flags
}

func (m *Fs9721lp3) extractValue(str string) (ret float64) {
	digits := map[string]string{
		"1111101": "0",
		"0000101": "1",
		"1011011": "2",
		"0011111": "3",
		"0100111": "4",
		"0111110": "5",
		"1111110": "6",
		"0010101": "7",
		"1111111": "8",
		"0111111": "9",
		"0000000": "",
		"1101000": "L",
	}

	arrValues := []string{
		str[5:12],  // Digito 01
		str[12:13], // Ponto 01
		str[13:20], // Digito 02
		str[20:21], // Ponto 02
		str[21:28], // Digito 03
		str[28:29], // Ponto 03
		str[29:36], // Digito 04
	}

	var value string
	for i, item := range arrValues {
		switch i % 2 {
		case 0:
			if val, exist := digits[item]; exist {
				value += val
			}
		case 1:
			if item == "1" {
				value += "."
			}
		}
	}

	ret, _ = strconv.ParseFloat(value, 64)
	if str[4:5] == "1" {
		ret = ret * -1
	}

	return
}

func (m *Fs9721lp3) extractUnit(str string) (unit string) {
	arrUnits := [][2]interface{}{
		{str[37:38] == "1", "n"},  // nano
		{str[36:37] == "1", "u"},  // micro
		{str[38:39] == "1", "k"},  // kilo
		{str[40:41] == "1", "m"},  // mili
		{str[42:43] == "1", "M"},  // mega
		{str[41:42] == "1", "%"},  // percent
		{str[45:46] == "1", "Ω"},  // ohm
		{str[48:49] == "1", "A"},  // amp
		{str[49:50] == "1", "V"},  // volts
		{str[44:45] == "1", "F"},  // cap
		{str[50:51] == "1", "Hz"}, // hertz
		{str[53:54] == "1", "°C"}, // temp
	}

	for _, item := range arrUnits {
		if item[0].(bool) {
			unit += item[1].(string)
		}
	}

	return
}

func (m *Fs9721lp3) extractFlags(str string) (flags []string) {
	arrFlags := [][2]interface{}{
		{str[0:1] == "1", "AC"},
		{str[1:2] == "1" && str[53:54] == "0", "DC"},
		{str[2:3] == "1", "Auto"},
		{str[39:40] == "1", "Diode test"},
		{str[43:44] == "1", "Conti test"},
		{str[44:45] == "1", "Capacity"},
		{str[46:47] == "1", "Rel"},
		{str[47:48] == "1", "Hold"},
		{str[52:53] == "1", "Min"},
		{str[55:56] == "1", "Max"},
		{str[51:52] == "1", "LowBat"},
	}

	for _, item := range arrFlags {
		if item[0].(bool) {
			flags = append(flags, item[1].(string))
		}
	}

	return
}
