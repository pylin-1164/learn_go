package postgres

import (
    "database/sql"
    "fmt"
    "github.com/larspensjo/config"
    _ "github.com/lib/pq"
    Logger "log"
)



type PgDatas struct {
    Host     string `json:"host,omitempty"`
    DbName   string `json:"dbName,omitempty"`
    UserName string `json:"userName,omitempty"`
    Password string `json:"password,omitempty"`
}

var Pg = &PgDatas{}

func init() {
    Pg.ReadConfig()
}

// read dbconfig to init
//if testFile is not empty,just use testFile
//if testFile is empty ,use the build path dbconfig.json
func (pgDatas *PgDatas) ReadConfig()  {
    cfg, err := config.ReadDefault("config.ini")
    if err != nil{
        Logger. Println("file config.ini not exits")
        panic("file config.ini not exits")
    }
    dbHost, _ := cfg.String("postgres", "host")
    dbName, _ := cfg.String("postgres", "db_name")
    dbUser, _ := cfg.String("postgres", "user_name")
    dbPassword, _ := cfg.String("postgres", "password")
    if dbHost == "" || dbName == "" || dbUser == "" || dbPassword == ""{
        Logger. Println("config.ini with postgres can't be null ")
        panic("config.ini with postgres can't be null ")
    }
    pgDatas = &PgDatas{
        Host:     dbHost,
        DbName:   dbName,
        UserName: dbUser,
        Password: dbPassword,
    }

    Pg = pgDatas
}

func (pg *PgDatas) ConnectPG() (*sql.DB, error) {
    connStr := "host=%s port=5432 dbname=%s  user=%s password=%s sslmode=disable"
    connStr = fmt.Sprintf(connStr, pg.Host, pg.DbName, pg.UserName, pg.Password)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        Logger.Printf("can not connect postgresql url: %s", err)
    }
    return db, err
}