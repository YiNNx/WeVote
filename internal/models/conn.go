package models

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/YiNNx/WeVote/internal/common/log"
	"github.com/YiNNx/WeVote/internal/config"
)

var (
	db  *gorm.DB
	rdb *redis.Ring
)

func initPostgresConn() {
	var err error

	db, err = gorm.Open(postgres.New(postgres.Config{DSN: config.C.Postgres.DSN}))
	if err != nil {
		log.Logger.Error("postgres connection failed: ", err)
		return
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Logger.Error("postgres migrate failed: ", err)
		return
	}

	log.Logger.Info("PostgreSQL server connected!")
}

func initRedisClusterConns() {
	addrs := config.C.Redis.Addrs
	rdbAddrs := make(map[string]string, len(addrs))
	for i, addr := range addrs {
		rdbAddrs[fmt.Sprintf("shard%d", i)] = addr
	}

	rdb = redis.NewRing(&redis.RingOptions{
		Addrs: rdbAddrs,
	})

	err := rdb.ForEachShard(context.Background(), func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})

	if err != nil {
		panic(err)
	}

	log.Logger.Info("Redis server connected!")
}

func InitIOWrapper() {
	initPostgresConn()
	initRedisClusterConns()
	initVoteDataWrapper()
}
