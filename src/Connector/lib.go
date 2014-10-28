package main

import (
	"bufio"
	"fmt"
)
import (
	"github.com/vmihailenco/msgpack"
)

func convertToByte(number uint16) []byte {
	return []byte{byte(number >> 8), byte(number & 255)}
}
func convertToUint16(number []byte) uint16 {
	return uint16(number[0])<<8 + uint16(number[1])
}

func readUint16(reader *bufio.Reader) uint16 {
	byte1, _ := reader.ReadByte()
	byte2, _ := reader.ReadByte()
	return convertToUint16([]byte{byte1, byte2})
}

func packageObject(cmd uint16, obj interface{}) []byte {
	data := []byte{'$'}
	data = append(data, convertToByte(cmd)...)

	objBytes, err := msgpack.Marshal(obj)
	if err != nil {
		panic(err)
	}
	data = append(data, convertToByte(uint16(len(objBytes)))...)
	data = append(data, objBytes...)
	return data
}

func readCmd(reader *bufio.Reader) uint16 {
	beginFlag, _ := reader.Peek(1)
	fmt.Println(string(beginFlag))
	cmdBytes, err := reader.Peek(2)
	if err != nil {
		panic(err)
	}

	return convertToUint16(cmdBytes)
}

func readMessage(reader *bufio.Reader) interface{} {
	var data = map[string]interface{}{}

	begin, _ := reader.ReadString('$')
	cmd := readUint16(reader)
	length := readUint16(reader)

	fmt.Println(begin)
	fmt.Println(cmd)
	fmt.Println(length)

	messageBytes := make([]byte, length)
	reader.Read(messageBytes)

	fmt.Println(len(messageBytes))
	fmt.Println(messageBytes)
	msgpack.Unmarshal(messageBytes, &data)
	return data
}
