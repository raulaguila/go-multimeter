# Bluetooth communication protocol 

This work was carried using bluetooth to connect to the multimeter and tracking the message bytes in each multimeter function.

Tested with Owon - OW18E multimeter

![](/screenshot/OW18E.png)

* ### Example message:

    * Receive a slice of byte with five positions
    * Example received message: [27 240 4 0 [50 128](#5th--6th-bytes)]
    * Convert to 8 digits binary: [[00011011](#1st-byte) [11110000](#2nd-byte) [00000100](#3rd-byte) [00000000](#4th-byte) 00110010 10000000]

* ### 1st byte

    * 8 digits, ex: [00011011](#example-message)
    * Digits (0, 1): Represents the function
    * Digits (2, 3, 4): Represents the unit of measure
    * Digits (5, 6, 7): Represents the range of the measured value

        | 0-1 | func | -   | 2-4 | unit | -   | 5-7 | range |
        | --- | ---  | --- | --- | ---  | --- | --- | ---   |
        | 00  | DC   |     | 001 | n    |     | 001 | 2000  |
        | 01  | AC   |     | 010 | µ    |     | 010 | 200   |
        | 10  | Diod |     | 011 | m    |     | 011 | 20    |
        | 11  | Cont |     | 100 | 1    |     | 100 | 2     |
        |     |      |     | 101 | k    |     | 111 | L     |
        |     |      |     | 110 | M    |     |     |       |

* ### 2nd byte

    * 8 digits, ex: [11110000](#example-message)
    * Digits (0, 1, 2, 3): Apparently they are not used [*](#ps)
    * Digits (4, 5, 6, 7): Represents the function

        | 0-3  | func | 4-7  | func        |
        | ---  | ---  | ---  | ---         |
        | 1111 | -    | 0000 | Voltage     |
        |      |      | 0001 | Resistance  |
        |      |      | 0010 | Continuity  |
        |      |      | 0100 | Capacitance |

* ### 3rd byte

    * 8 digits, ex: [00000100](#example-message)
    * Digits (0, 1, 2, 3): Apparently they are not used [*](#ps)
    * Digits (4, 5, 6, 7): Represents if the multimeter is in automatic range

        | 0-3  | func | 4-7  | func     |
        | ---  | ---  | ---  | ---      |
        | 0000 | -    | 0100 | Auto ON  |
        |      |      | 0000 | Auto OFF |

* ### 4th byte

    * 8 digits, ex: [00000000](#example-message)
    * Apparently this byte is not used [*](#ps)

* ### 5th & 6th bytes

    * Ex: [\[50 128\]](#example-message)
    * Represents the measurement value
    * Use them without converting to binary
    * 6th byte counts the overflow of 5th byte
    * If the 5th byte >= 128, it is a negative value

* ### Combinations

    * Combining the first byte and second byte items, we have the final function
    * First two characters of first byte with last four characters of second byte

        | Conbination                                 | Final function      | unity |
        | ---                                         | ---                 | ---   |
        | [DC](#1st-byte) + [Continuity](#2nd-byte)   | Temperature         | ºC    |
        | [AC](#1st-byte) + [Continuity](#2nd-byte)   | Temperature         | ºF    |
        | [Cont](#1st-byte) + [Continuity](#2nd-byte) | Continuity test     | Ω     |
        | [Diod](#1st-byte) + [Continuity](#2nd-byte) | Diode test          | V     |
        | [Diod](#1st-byte) + [Resistance](#2nd-byte) | Frequence           | Hz    |
        | [DC](#1st-byte) + [Voltage](#2nd-byte)      | DC Voltage Measure  | V     |
        | [AC](#1st-byte) + [Voltage](#2nd-byte)      | AC Voltage Measure  | V     |
        | [DC](#1st-byte) + [Resistance](#2nd-byte)   | Resistance Measure  | Ω     |
        | [AC](#1st-byte) + [Resistance](#2nd-byte)   | Capacitance Measure | F     |
        | [Diod](#1st-byte) + [Voltage](#2nd-byte)    | Current Measure     | A     |

* #### PS

    * \* Maybe the unused data are used to represent other information that was not tracked, such as low battery, MIN, MAX, and others.