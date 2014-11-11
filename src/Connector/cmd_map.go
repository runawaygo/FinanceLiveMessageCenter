package main

const (
	PING = 0
	PONG = 1

	//心跳消息
	PUMP = 3
)

const (
	INTERNAL_ERROR = 500
)

const (
	//普通消息
	MESSAGE   = 20000
	VOICE     = 20001
	TELEPHONE = 20002

	//认证消息
	CUSTOMER_AUTH = 21001
	CSR_AUTH      = 21002

	AUTH_SUCCESS = 21102
	AUTH_FAILED  = 21103

	CSRINFO           = 22001
	CSR_STATUS_CHANGE = 22002
)
