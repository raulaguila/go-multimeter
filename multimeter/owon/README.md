# Bluetooth communication protocol 

This work was carried using bluetooth to connect to the multimeter and tracking the message bytes in each multimeter function.

Tested with [Owon - OW18E Digital Multimeter](https://owon.com.hk/products_owon_ow18d%7Ce_4_1%7C2_digits__handheld_digital_multimeter)

![](/screenshot/ow18e.png)

### Example message:

* Receive a 5-byte array
* Example received: [ 98 240 4 0 [147 49](#5th--6th-bytes) ]
* Convert to 8-digit binary: [ [01100010](#1st-byte) [11110000](#2nd-byte) [00000100](#3rd-byte) [00000000](#4th-byte) 10010011 00110001 ]
* Final output struct: `<value> <unit> <flags>`
* Final output of example: `126.91 V [AC, Auto Range]`

### 1st Byte

* 8 bits, ex: [01100010](#example-message)

    * Bits (0, 1): Represents the function (On example: *01*)
    * Bits (2, 3, 4): Represents the unit of measure (On example: *100*)
    * Bits (5, 6, 7): Represents the range of the measured value (On example: *010*)

    | 0-1 | func | -   | 2-4 | unit | -   | 5-7 | range |
    | --- | ---  | --- | --- | ---  | --- | --- | ---   |
    | 00  | DC   |     | 001 | n    |     | 000 | NCV   |
    | 01  | AC   |     | 010 | µ    |     | 001 | 2000  |
    | 10  | Diod |     | 011 | m    |     | 010 | 200   |
    | 11  | Cont |     | 100 | 1    |     | 011 | 20    |
    |     |      |     | 101 | k    |     | 100 | 2     |
    |     |      |     | 110 | M    |     | 111 | L     |

### 2nd Byte

* 8 bits, ex: [11110000](#example-message)

    | Bits | Value | Function                              |
    | ---: | :---: | :---                                  |
    | 0-3  | 1111  | Apparently they are not used [*](#ps) |
    | 4-5  | 00    | Apparently they are not used [*](#ps) |
    | 6-7  | 00    | Voltage                               |
    | 6-7  | 01    | Resistance                            |
    | 6-7  | 10    | Continuity                            |
    | 6-7  | 11    | NCV                                   |

### 3rd Byte

* 8 bits, ex: [00000100](#example-message)

    | Bits | Value | Function                              |
    | ---: | :---: | :---                                  |
    | 0-3  | 0000  | Apparently they are not used [*](#ps) |
    | 4    | 1     | Low Battery                           |
    | 5    | 1     | Auto range enabled                    |
    | 6    | 1     | Relative mode enabled                 |
    | 7    | 1     | Hold enabled                          |

### 4th Byte

* 8 bits, ex: [00000000](#example-message)
* Apparently this byte is unused [*](#ps)

### 5th & 6th Bytes

* Ex: [\[... 147 49\]](#example-message)
* Represents the measurement value
* Use them without converting to binary
* 6th byte counts the overflow of 5th byte
* If the 5th byte >= 128, it is a negative value

### Final function

* Combining the 1st and 2nd byte items, it has the final function.
* [First two bits of 1st byte](#1st-byte) + [last two of 2nd](#2nd-byte).

    | 1st Byte | 2nd Byte   | Final function      | Symbol |
    | ---:     | :---       | :---                | :---:  |
    | DC       | Continuity | Temperature         | ºC     |
    | DC       | Resistance | Resistance Measure  | Ω      |
    | DC       | Voltage    | DC Voltage Measure  | V      |
    | AC       | Continuity | Temperature         | ºF     |
    | AC       | NCV        | NCV Measure         | -      |
    | AC       | Resistance | Capacitance Measure | F      |
    | AC       | Voltage    | AC Voltage Measure  | V      |
    | Diod     | Continuity | Diode test          | V      |
    | Diod     | Resistance | Frequence Measure   | Hz     |
    | Diod     | Voltage    | Current Measure     | A      |
    | Cont     | Continuity | Continuity test     | Ω      |
    | Cont     | Resistance | Frequence Measure   | %      |

#### PS

* [\*](#ps) Maybe the unused data are used to represent other information that was not tracked, such as MIN, MAX, and others.