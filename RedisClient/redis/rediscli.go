package util

import (
	"fmt"
	"github.com/go-redis/redis"
)

const PWD = "123456"

type Rediscli struct {
}

type Config struct {
	ServerIP  string
	RedisPort int
	RedisPwd  string
}

//Cluster Mode
//return the Client
func (r *Rediscli) ConnectRedisCluster(redisArr []Config) (*redis.Cmdable,*redis.ClusterClient) {
	var clusterCli redis.Cmdable
	addrs := make([]string, len(redisArr))
	for i, redisConfig := range redisArr {
		addrs[i] = fmt.Sprintf("%s:%d", redisConfig.ServerIP, redisConfig.RedisPort)
	}

	cluster := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
		//初始化默认密码
		Password:     PWD,
		PoolSize:     1000,
		MinIdleConns: 100,
	})
	clusterCli = cluster
	return &clusterCli,cluster
}

//Single Mode
//return the Client
func (r *Rediscli) ConnectRedis(config Config) *redis.Cmdable {
	var client redis.Cmdable
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.ServerIP, config.RedisPort), // use default Addr
		Password: PWD,                                                     // no password set
		DB:       0,                                                       // use default DB
	})
	fmt.Printf("connect reids client memery address : %v ; \n", &client)
	return &client
}

//Master-Slave Mode
//return Master client
func (r *Rediscli) ConnectFailoverRedis() *redis.Cmdable {
	var redisdb redis.Cmdable
	redisdb = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: []string{"172.16.134.20:8000", "172.16.134.20:8001", "172.16.134.20:8002"},
		Password:      PWD,
	})
	return &redisdb
}
