package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	testData := map[string]interface{}{"name": "superwolf", "sex": 1}

	if err != nil {
		// handle error
		fmt.Println(err)
		return
	}

	for {
		data := packageObject(2, testData)
		fmt.Println(data)
		conn.Write(data)

		status, err := bufio.NewReader(conn).ReadString('3')
		if err != nil {
			fmt.Println("Error")
			fmt.Println(err)
			continue
		}
		fmt.Println(status)
		fmt.Println(time.Now())

	}

}
