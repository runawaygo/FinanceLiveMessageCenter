package main

const (
	PING = 0
	PONG = 1

	//心跳消息
	CLOSE = 2
	PUMP  = 3
)

const (
	NOT_ACCEPTABLE = 406
	INTERNAL_ERROR = 500
)

const (
	AUTH_FAILED  = 1000
	AUTH_SUCCESS = 1001

	//认证消息
	CSR_AUTH      = 1100
	CUSTOMER_AUTH = 1101
)

const (
	//普通消息
	MESSAGE   = 20000
	VOICE     = 20001
	TELEPHONE = 20002

	CSRINFO           = 21000
	CSR_STATUS_CHANGE = 21001
)
