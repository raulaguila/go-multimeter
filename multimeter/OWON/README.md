# Bluetooth communication protocol 

This work was carried using bluetooth to connect to the multimeter and tracking the message bytes in each multimeter function.

Tested with Owon - OW18E multimeter

![](/screenshot/OW18E.png)

* ### Example message:

    * Receive a slice of byte with five positions
    * Example received message: [27 240 4 0 50 128]
    * Convert to 8 digits binary: [00011011 11110000 00000100 00000000 00110010 10000000]

* ### 1st byte (8 digits, ex: [00011011])

    * Digits (0, 1): Represents the function
    * Digits (2, 3, 4): Represents the unit of measure
    * Digits (5, 6, 7): Represents the range of the measured value

        | 1-2 | func | 3-5 | unit | 6-8 | range |
        | --- | ---  | --- | ---  | --- | ---   |
        | 00  | DC   | -   | -    | -   | -     |
        | 01  | AC   | -   | -    | -   | -     |
        | 10  | Diod | -   | -    | -   | -     |
        | 11  | Cont | -   | -    | -   | -     |
        | -   | -    | 001 | n    | -   | -     |
        | -   | -    | 010 | µ    | -   | -     |
        | -   | -    | 011 | m    | -   | -     |
        | -   | -    | 100 | 1    | -   | -     |
        | -   | -    | 101 | k    | -   | -     |
        | -   | -    | 110 | M    | -   | -     |
        | -   | -    | -   | -    | 001 | 2000  |
        | -   | -    | -   | -    | 010 | 200   |
        | -   | -    | -   | -    | 011 | 20    |
        | -   | -    | -   | -    | 100 | 2     |
        | -   | -    | -   | -    | 111 | L     |

* ### 2nd byte (8 digits, ex: [11110000])

    * Digits (0, 1, 2, 3): Apparently they are not used [*](#ps)
    * Digits (4, 5, 6, 7): Represents the function

        | 1-4  | func | 5-8  | func        |
        | ---  | ---  | ---  | ---         |
        | 1111 | -    | 0000 | Voltage     |
        | -    | -    | 0001 | Resistance  |
        | -    | -    | 0010 | Continuity  |
        | -    | -    | 0100 | Capacitance |

* ### 3rd byte (8 digits, ex: [00000100])

    * Digits (0, 1, 2, 3): Apparently they are not used [*](#ps)
    * Digits (4, 5, 6, 7): Represents if the multimeter is in automatic range

        | 1-4  | func | 5-8  | func     |
        | ---  | ---  | ---  | ---      |
        | 0000 | -    | 0100 | Auto ON  |
        | -    | -    | 0000 | Auto OFF |

* ### 4th byte (8 digits, ex: [00000000])

    * Apparently this byte is not used [*](#ps)

* ### 5th and 6th bytes (8 digits each, ex: [50 128])

    * Represents the measurement value
    * Use them without converting to binary
    * 6th byte counts the overflow of 5th byte
    * If the 5th byte >= 128, it is a negative value

* ### Combinations

    * Combining the first byte and second byte items, we have the final function
    * First two characters of first byte with last four characters of second byte

        | Conbination          | Final function      | unity |
        | ---                  | ---                 | ---   |
        | DC + Continuity      | Temperature         | ºC    |
        | AC + Continuity      | Temperature         | ºF    |
        | Cont + Continuity    | Continuity test     | Ω     |
        | Diod + Continuity    | Diode test          | V     |
        | Diod + Resistance    | Frequence           | Hz    |
        | DC + Voltage         | DC Voltage Measure  | V     |
        | AC + Voltage         | AC Voltage Measure  | V     |
        | DC + Resistance      | Resistance Measure  | Ω     |
        | AC + Resistance      | Capacitance Measure | F     |
        | Diod + Voltage       | Current Measure     | A     |

* #### PS

    * \* Maybe the unused data are used to represent other information that was not tracked, such as low battery, MIN, MAX, and others.