package data

import (
	"context"
	"convoo-accounts/internal/conf"
	"fmt"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/types/known/durationpb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	DB    *gorm.DB
	Redis *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {

	// 创建数据库连接
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// 2. 连接 Redis
	readTimeout := durationToDuration(c.Redis.ReadTimeout)
	writeTimeout := durationToDuration(c.Redis.WriteTimeout)
	redisClient := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	})

	// 测试 Redis 连接
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("failed to get db instance for cleanup: %v", err)
		}
		sqlDB.Close() // 在应用结束时关闭数据库连接
		log.NewHelper(logger).Info("closing the data resources")
	}

	return &Data{
		DB:    db,
		Redis: redisClient,
	}, cleanup, nil
}

// 转换 google.protobuf.Duration 到 time.Duration
func durationToDuration(d *durationpb.Duration) time.Duration {
	return time.Duration(d.GetSeconds())*time.Second + time.Duration(d.GetNanos())*time.Nanosecond
}
