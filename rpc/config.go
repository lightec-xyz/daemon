package rpc

import "time"

const BatchRequestLimit = 1000

const BatchResponseMaxSize = 25 * 1000 * 1000

const HttpReadTimeOut = 15 * time.Second
const HttpWriteTimeOut = 15 * time.Second
const MaxHeaderBytes = 10 << 20

type Permission int

const (
	NonePermission Permission = iota + 1
	JwtPermission
)

const (
	AuthorizationHeader = "Authorization"
)
