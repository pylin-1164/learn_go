package mapper

import (
	"fmt"
	util "redis.cli/redis"
	"strings"
)

type RedisCacheFun struct {
	RedisSavedKeys		[]string 	 // Redis 保存过程中的KEY
	RedisOldExitsKeys 	[]string	// Redis 操作前已存在的所有的KEY
}

func NewRedisCacheFun(t RedisCacheIntf) *RedisCacheFun{
	return &RedisCacheFun{RedisSavedKeys:make([]string,0),RedisOldExitsKeys:t.QueryAllKeys()}
}

func (r *RedisCacheFun)RedisSave(t RedisCacheIntf){
	saveRedisKeys := t.RedisSave()
	r.RedisSavedKeys = append(r.RedisSavedKeys,saveRedisKeys...)
}

func (r *RedisCacheFun) RemoveInvalidKeys(){
	redisUtil := util.GetInstance()
	savedRedisKeys := fmt.Sprintf("/%s/", strings.Join(r.RedisSavedKeys, "/"))
	for _,oldKey := range r.RedisOldExitsKeys {
		keyMatch := fmt.Sprintf("/%s/",oldKey)
		//如果历史KEY和新保存的KEY一致，则跳过
		if strings.Contains(savedRedisKeys,keyMatch){
			continue
		}
		//如果历史KEY在新保存的KEY中不存在，则删除
		redisUtil.Del(oldKey)
	}
}

func RedisRead(t RedisCacheIntf){
	t.RedisRead()
}

/**
 * @Date 2020-09-11
 * @author pyl
 * @Description 所有的Redis缓存数据均需要继承一下接口，方便统一维护
 */
type RedisCacheIntf interface {

	RedisSave() []string

	RedisRead()

	QueryAllKeys() []string

}