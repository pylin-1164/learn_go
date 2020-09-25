package listen

import (
	Logger "log"
	"redis.cli/postgres"
	"redis.cli/redis"
	"time"
)

// 动态监控Redis版本号是否刷新，同步更新缓存
var VERSION = ""
var VERSION_KEY = "redis_cache_version"


type RedisListenUpdate struct {

}

//轮循获取版本号
func (r *RedisListenUpdate) StartListen() {

	for {
		var redisUtil = util.GetInstance()
		version := redisUtil.Get(VERSION_KEY)
		if VERSION == "" && version == ""{
			//初始化
			redisUtil.Incr(VERSION_KEY)
			VERSION = "1"
			freshRedisDatas()
			Logger.Print("Update Redis Datas By Init")

		} else if version != VERSION {
			Logger.Printf("version update from %s to %s", VERSION, version)
			VERSION = version
			freshRedisDatas()
			Logger.Printf("Update Redis Datas ... ")
		}
		time.Sleep(time.Second * 10)
	}
}

func freshRedisDatas() {
	//创建数据库连接
	pgDatas := postgres.Pg
	db, _ := pgDatas.ConnectPG()
	defer db.Close()

	//缓存所有的用户信息
	pgDatas.MasterDBCache(db)


}

