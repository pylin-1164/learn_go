### RedisClient 基于go-redis的客户端再封装

概述   
---

由于项目部署方式的多样性，在服务器配置和业务较为简单的情况下。
需要在集群和单实例的模式下快速切换，同时对代码进行复用使代码更简单

基于以上需求，在go-redis上层做了Redis相关方法的再封装

Redis配置
---
[config.ini]

    [redis_server]
    #single = 127.0.0.1:10000
    cluster = 127.0.0.1:10000,127.0.0.1:10001,127.0.0.1:10002,127.0.0.1:10003,127.0.0.1:10004,127.0.0.1:10005

[example]
```cgo
// 定义 MasterIdMap 继承 RedisCacheIntf 接口

//声明Redis存储的数据类型，初始化已存在的key，用于后续删除过期数据
cacheFun := mapper.NewRedisCacheFun(&MasterIdMap{})


for ...
    // TODO 构造要存储的MasterIdMap数据
    redisMaster := mapper.MasterIdMap{
        NameKey 		string
        MasterId 		string,
    }

    //缓存到Redis,存在则覆盖
    cacheFun.RedisSave(&redisMaster)
end 


//删除过期的无效KEY
cacheFun.RemoveInvalidKeys()

```

