package owon

var (
	DeviceName               string   = "BDM"
	ServiceUUID              [16]byte = [16]byte{0x00, 0x00, 0xff, 0xf0, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0x80, 0x5f, 0x9b, 0x34, 0xfb}
	CharacteristicNotifyUUID [16]byte = [16]byte{0x00, 0x00, 0xff, 0xf4, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0x80, 0x5f, 0x9b, 0x34, 0xfb}
	CharacteristicWriteUUID  [16]byte = [16]byte{0x00, 0x00, 0xff, 0xf3, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0x80, 0x5f, 0x9b, 0x34, 0xfb}
)
