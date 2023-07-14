package login

import (
	"app-backend/service"
)

var (
	dbCtx = service.DbCtx
	cacheClient   = service.CacheClient
	mongoClient = service.MongoDbClient
)

type Result struct {
	Error     interface{} `json:"error,omitempty"`
	Message   interface{} `json:"message,omitempty"`
	MessageTh interface{} `json:"message_th,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Total     int         `json:"total,omitempty"`
	Count     int         `json:"count,omitempty"`
	Status    int         `json:"status,omitempty"`
}
