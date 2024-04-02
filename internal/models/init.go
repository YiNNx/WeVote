package models

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/YiNNx/WeVote/pkg/log"
)

var (
	db  *gorm.DB
	rdb *redis.Ring
)

func InitDataBaseConnections(pgHost, pgPort, pgUser, pgPassword, pgDbname string, redisAddrs []string) {
	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		pgHost, pgUser, pgPassword, pgDbname, pgPort,
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Logger.Error("postgres connection failed: ", err)
		return
	}
	err = db.AutoMigrate(&Vote{})
	if err != nil {
		log.Logger.Error("postgres migrate failed: ", err)
		return
	}

	log.Logger.Info("PostgreSQL server connected!")

	rdbAddrs := make(map[string]string, len(redisAddrs))
	for i, addr := range redisAddrs {
		rdbAddrs[fmt.Sprintf("shard%d", i)] = addr
	}

	rdb = redis.NewRing(&redis.RingOptions{
		Addrs: rdbAddrs,
	})

	err = rdb.ForEachShard(context.Background(), func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})

	if err != nil {
		panic(err)
	}

	log.Logger.Info("Redis server connected!")
}
