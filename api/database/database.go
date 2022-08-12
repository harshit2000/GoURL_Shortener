package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

var Ctx  = context.Background()

func CreateClient(dbNo int) *redis.Client {
