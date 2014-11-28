package main

func convertToByte(number uint16) []byte {
	return []byte{byte(number >> 8), byte(number & 255)}
}
func convertToUint16(bytes []byte) uint16 {
	return uint16(bytes[0])<<8 + uint16(bytes[1])
}
