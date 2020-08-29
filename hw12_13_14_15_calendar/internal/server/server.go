package server

type key string

const (
	UserIDHeaderKey string = "User_id"
	UserIDGrpcKey   key    = "Grpc-Metadata-User_id"
)
