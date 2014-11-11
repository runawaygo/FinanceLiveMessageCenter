package main

func csrAuth(csrInfo *map[interface{}]interface{}) (string, error) {

	return "superowlf", nil
}

func customerAuth(customerInfo *map[interface{}]interface{}) (string, error) {

	return "superowlf", nil
}

func pipe(message *Message) {
	println(message)
}

func authHandler(message *Message) (string, error) {
	switch message.Cmd {
	case CUSTOMER_AUTH:
		return customerAuth(message.Content)
	case CSR_AUTH:
		return csrAuth(message.Content)
	default:
		panic("传入消息类型不合法!")
	}
}

func messageHandler(message *Message) {
	pipe(message)
}
