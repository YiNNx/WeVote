package models

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/YiNNx/WeVote/internal/common/log"
)

var (
	db  *gorm.DB
	rdb *redis.Ring
)

type rtx struct {
	redis.Pipeliner
	ctx context.Context
}

type tx struct {
	*gorm.DB
}

func BeginRedisTx(ctx context.Context) rtx {
	return rtx{
		Pipeliner: rdb.TxPipeline(),
		ctx:       ctx,
	}
}

func BeginPostgresTx() tx {
	return tx{db.Begin()}
}

func initPostgresConn(dsn string) {
	var err error

	db, err = gorm.Open(postgres.New(postgres.Config{DSN: dsn}))
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

func initRedisClusterConns(addrs []string) {
	var err error

	rdbAddrs := make(map[string]string, len(addrs))
	for i, addr := range addrs {
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

func InitIOWrapper(postgresDSN string, redisAddrs []string) {
	initPostgresConn(postgresDSN)
	initRedisClusterConns(redisAddrs)
	initVoteDataWrapper(rdb, db)
}
