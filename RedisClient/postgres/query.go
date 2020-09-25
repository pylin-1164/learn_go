package postgres

import (
	Logger "log"
	"database/sql"
	"encoding/json"
	"redis.cli/postgres/mapper"
)

/**
 * @author pyl
 * @Date 2020-09-11
 */

//主账号信息缓存
func (pg *PgDatas) MasterDBCache(db *sql.DB) {
	masterSQL := `SELECT id,name,password,mobile FROM master `
	masterRows, err := db.Query(masterSQL)
	if err != nil{
		Logger.Printf("MasterDBCache SQL Query Error : %s \n",err)
	}
	masterIdMap := mapper.MasterIdMap{}
	masterNameIdMap := mapper.MasterNameIdMap{}
	defer masterRows.Close()

	cacheFun := mapper.NewRedisCacheFun(&mapper.RedisMaster{})

	//遍历结果集
	for masterRows.Next() {
		masterData := mapper.MasterData{}
		err = masterRows.Scan(&masterData.Id, &masterData.Name, &masterData.Password, &masterData.Mobile)
		if err != nil{
			Logger.Printf("master row scan Error : %s \n",err)
			continue
		}
		masterIdMap.IdKey = masterData.Id
		if masterRedisValue, err := json.Marshal(masterData); err != nil {
			Logger.Printf("json marshal master error with master id[%s] and name[%s] \n", masterData.Id, masterData.Name)
			continue
		} else {
			masterIdMap.Value = string(masterRedisValue)
		}

		masterNameIdMap.NameKey = masterData.Name
		masterNameIdMap.MasterId = masterData.Id

		//Redis缓存用户信息和用户ID与用户名称关联信息
		redisMaster := mapper.RedisMaster{
			MasterIdMap:     &masterIdMap,
			MasterNameIdMap: &masterNameIdMap,
		}

		//缓存到Redis
		cacheFun.RedisSave(&redisMaster)
	}

	//删除过期的无效KEY
	cacheFun.RemoveInvalidKeys()

}

