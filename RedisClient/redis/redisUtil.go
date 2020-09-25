package util

import (
    "errors"
    "fmt"
    "github.com/go-redis/redis"
    "github.com/larspensjo/config"
    "reflect"
    "strconv"
    "strings"
    "sync"
    Logger "log"
)

var ru *RedisUtil
var once sync.Once

/**
 *  Redis支持批量执行，pipline和mset等，但是在集群模式下，go-redis，jredis等都未做很好的支持
 */

//	get instance util and  init redis client connect
func GetInstance() *RedisUtil {
    once.Do(func() {
        ru = &RedisUtil{}
        ru.redisUtil()
    })
    return ru
}

type RedisUtil struct {
    client      *redis.Cmdable
    cluster     *redis.ClusterClient
    index       int
    IPAddr      string
    NatAddr     string
}

//根据配置决定是redis连接方式：集群或单服务
func (ru *RedisUtil) redisUtil()  {
    r := Rediscli{}
    //return r.ConnectRedis()
    configs := ru.getRedisConfig()
    if configs["redis_single"] != "" {
        split := strings.Split(configs["redis_single"], ":")
        port, _ := strconv.Atoi(split[1])
        config := Config{split[0], port, ""}
        ru.client = r.ConnectRedis(config)
        return
    } else if configs["redis_cluster"] != "" {
        clusters := strings.Split(configs["redis_cluster"], ",")
        var addrConfigs = make([]Config, len(clusters))
        for i, cluster := range clusters {
            split := strings.Split(cluster, ":")
            port, _ := strconv.Atoi(split[1])
            addrConfigs[i] = Config{split[0], port, ""}
        }
        ru.client,ru.cluster = r.ConnectRedisCluster(addrConfigs)
        return
    }
    Logger.Println("config.ini with redis can't be find")
    panic("config.ini with redis can't be find")
}

func (ru *RedisUtil) GetClient() redis.Cmdable {
    return *ru.client
}

// get redis value with key
func (ru *RedisUtil) Get(key string) string {
    if value, err := ru.GetClient().Get(key).Result(); err == nil {
        return value
    } else if err.Error() != "nil" {
        fmt.Printf("Redis Get %s ERROR : %v \n", key, err)
    }
    return ""
}

// get redis all keys like pattern
func (ru *RedisUtil) Keys(pattern string)[]string{
    keys := make([]string, 0)
    // 集群单独处理
    if strings.Contains(reflect.TypeOf(ru.GetClient()).String(),"redis.ClusterClient"){
        ru.cluster.ForEachMaster(func(client *redis.Client) error {
            keys = append(keys,scanKeys(client, pattern) ... )
            return nil
        })
    } else {
        keys = append(keys,scanKeys(ru.GetClient(),pattern) ... )
    }

    return keys
}

func scanKeys(client redis.Cmdable, pattern string) []string{
    keys := make([]string, 0)
    iter := client.Scan(0, pattern, 0).Iterator()
    for iter.Next() {
        keys = append(keys,iter.Val())
    }
    return keys
}

// override if exists or add
func (ru *RedisUtil) Set(key string, value interface{}) bool {
    cmd := ru.GetClient().Set(key, value, 0)
    value, err := cmd.Result()
    if err != nil {
        fmt.Printf("Redis Set ERROR : %v \n", err)
        return false
    }
    return true
}

// delete all keys if exists ,ignore keys if not exits
func (ru *RedisUtil) Del(keys ...string) bool {
    if _, err := ru.GetClient().Del(keys...).Result(); err != nil {
        fmt.Printf("Redis Del keys %v ERROR :  %v \n", keys, err)
        return false
    }
    return true
}

// get redis hash value with key
func (ru *RedisUtil) HGet(key string, field string) string {
    if value, err := ru.GetClient().HGet(key, field).Result(); err == nil {
        return value
    } else {
        Logger.Printf("Redis HGet key:%s field:%s ERROR : %v", key, field, err)
    }
    return ""
}

//get redis all hash map with key
func (ru *RedisUtil) HGETAll(key string) map[string]string {
    if value, err := ru.GetClient().HGetAll(key).Result(); err == nil {
        return value
    } else {
        fmt.Printf("Redis HGetAll key:%s ERROR : %v \n", key, err)
    }
    return nil
}

// set redis hash field-value with key
func (ru *RedisUtil) HSet(key string, field string, value interface{}) bool {
    if _, err := ru.GetClient().HSet(key, field, value).Result(); err != nil {
        fmt.Printf("Redis HSet key:%s  key-value:[%s:%s] ERROR : %v \n", key, field, value, err)
        return false
    }
    return true
}

//multiple set hash filed-value with key
//example: map[string]interface{}{"EN":"78","CN":"45","SX":"66"}
func (ru *RedisUtil) HMSet(key string, maps map[string]interface{}) bool {

    if _, err := ru.GetClient().HMSet(key, maps).Result(); err != nil {
        fmt.Printf("Redis HMSet key:%s  map:%v ERROR : %v \n", key, maps, err)
        return false
    }
    return true
}

//add item to collect
func (ru *RedisUtil) SAdd(key string, values interface{}) bool {
    if _, err := ru.GetClient().SAdd(key, values).Result(); err != nil {
        fmt.Printf("Redis SAdd key:%s  value:%v ERROR : %v \n", key, values, err)
        return false
    }
    return true
}

//Get The collect values by key
func (ru *RedisUtil) SMembers(key string) []string {
    var values []string
    var err error
    if values, err = ru.GetClient().SMembers(key).Result(); err != nil {
        fmt.Printf("Redis SMembers key:%s ERROR: %v\n", key, err)
        return make([]string, 0)
    }
    return values
}

//自增索引
func (ru *RedisUtil) Incr(key string) string {
    value, _ := ru.GetClient().Incr(key).Result()
    return fmt.Sprintf("%v", value)
}

func (r *RedisUtil) getRedisConfig() map[string]string {
    cfg, err := config.ReadDefault("config.ini")
    if err != nil {
        Logger.Fatalln("读取配置文件config.ini失败")
    }

    redisSingleServer, _ := cfg.String("redis_server","single")

    redisClusterServer, _ := cfg.String("redis_server","cluster")

    if redisSingleServer == "" && redisClusterServer == ""{
        Logger.Println("读取Redis Server 失败")
        panic(errors.New(" Read Redis Service Setting ERROR "))
    }

    data := make(map[string]string)
    //Redis 独立服务
    if redisSingleServer != ""{
        data["redis_single"] = redisSingleServer
    }

    //Redis 集群服务
    if redisClusterServer != ""{
        data["redis_cluster"] = redisClusterServer
    }

    Logger.Printf("data is : %v", data)
    return data
}


func (r *RedisUtil) GetServerIp() string{
    if r.NatAddr == "" {
        return r.IPAddr
    }
    return r.NatAddr
}
