package models

import (
	"context"
	"fmt"

	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/YiNNx/WeVote/internal/common/log"
	"github.com/YiNNx/WeVote/internal/config"
)

var (
	db *gorm.DB

	rdb     *redis.Ring
	redlock *redsync.Redsync
)

func initPostgresConn() {
	var err error

	db, err = gorm.Open(
		postgres.New(postgres.Config{DSN: config.C.Postgres.DSN}))
	if err != nil {
		log.Logger.Error("postgres connection failed: ", err)
		return
	}
	err = db.AutoMigrate(&Vote{})
	if err != nil {
		log.Logger.Error("postgres migrate failed: ", err)
		return
	}

	// set the connections pool
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)

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

	syncPools := []redsyncredis.Pool{}
	err := rdb.ForEachShard(context.Background(), func(ctx context.Context, shard *redis.Client) error {
		if err := shard.Ping(ctx).Err(); err != nil {
			return err
		}
		syncPools = append(syncPools, goredis.NewPool(shard))
		return nil
	})
	if err != nil {
		log.Logger.Fatal(err)
	}

	redlock = redsync.New(syncPools...)

	log.Logger.Info("Redis server connected!")
}

func InitIOWrapper() {
	initPostgresConn()
	initRedisClusterConns()
	initVoteDataWrapper()
}
