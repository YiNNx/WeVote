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

func InitDataBaseConnections(postgresDSN string, redisAddrs []string) {
	var err error

	db, err = gorm.Open(postgres.New(postgres.Config{DSN: postgresDSN}))
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

type rtx struct {
	redis.Pipeliner
	ctx context.Context
}

func BeginRedisTx(ctx context.Context) rtx {
	return rtx{
		Pipeliner: rdb.Pipeline(),
		ctx:       ctx,
	}
}

type tx struct {
	*gorm.DB
}

func BeginPostgresTx() tx {
	return tx{db.Begin()}
}
