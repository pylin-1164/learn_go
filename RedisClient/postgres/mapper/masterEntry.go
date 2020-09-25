package mapper

import (
	"encoding/json"
	"fmt"
	Logger "log"
	util "redis.cli/redis"
)

/**
 * REDIS KEY键映射常量
 */

// 主账号，支持根据ID查询
type MasterIdMap struct {
	IdKey  			string
	Value 			string
	Data			*MasterData
}

//主账号属性
type MasterData struct {
	Id 				string 				`json:"id"`
	Name 			string 				`json:"name"`
	Password 		string 				`json:"password"`
	Mobile 			string 				`json:"mobile"`
}

//主账号ID和Name映射关系
type MasterNameIdMap struct{
	NameKey 		string
	MasterId 		string
}

//根据主账号ID获取到Redis中的KEY
func (m *MasterIdMap)RedisMasterIdKey() string{
	idKey := "master:master_id:%s"
	masterIdKey := fmt.Sprintf(idKey, m.IdKey)
	return masterIdKey
}

//根据主账号名称查询主账号的ID
func (m *MasterNameIdMap)RedisMasterNameKey() string{
	nameKey := "master_name:%s:master_id"
	masterNameKey := fmt.Sprintf(nameKey, m.NameKey)
	return masterNameKey
}


type RedisMaster struct {
	MasterIdMap 		*MasterIdMap
	MasterNameIdMap 	*MasterNameIdMap
	QueryMasterId 		string
	QueryMasterName		string
}

func (r *RedisMaster) RedisSave() []string{
	redisUtil := util.GetInstance()
	//保存到Redis
	redisUtil.Set(r.MasterIdMap.RedisMasterIdKey(), r.MasterIdMap.Value)
	redisUtil.Set(r.MasterNameIdMap.RedisMasterNameKey(), r.MasterNameIdMap.MasterId)
	return []string{r.MasterIdMap.RedisMasterIdKey(),r.MasterNameIdMap.RedisMasterNameKey()}
}

func (r *RedisMaster) QueryAllKeys() (keys []string){
	redisUtil := util.GetInstance()
	masterIdKeys := redisUtil.Keys("master:master_id:*")
	masterNameIdKeys := redisUtil.Keys("master_name:*:master_id")

	if masterIdKeys != nil && masterNameIdKeys != nil{
		return append(masterIdKeys,masterNameIdKeys...)
	}
	Logger.Println("Redis Master Query All Keys ERROR")
	return
}

func (r *RedisMaster) RedisRead() {
	redisUtil := util.GetInstance()
	masterId := r.QueryMasterId
	masterName := r.QueryMasterName

	//根据用户名查询用户
	if masterId == "" && masterName != ""{
		masterNameIdMap := MasterNameIdMap{NameKey: masterName}
		masterId = redisUtil.Get(masterNameIdMap.RedisMasterNameKey())
	}

	//根据用户ID查询
	if masterId != ""{
		masterIdMap := MasterIdMap{IdKey: masterId,Data:&MasterData{}}
		masterInf := redisUtil.Get(masterIdMap.RedisMasterIdKey())
		json.Unmarshal([]byte(masterInf),masterIdMap.Data)
		r.MasterIdMap = &masterIdMap
		return
	}

}
