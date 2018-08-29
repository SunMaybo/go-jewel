package context

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	//_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/go-redis/redis"
	"github.com/cihub/seelog"
	"time"
)

type Db struct {
	MysqlDb     *gorm.DB
	PostDb      *gorm.DB
	SqlServerDb *gorm.DB
	Sqlite3Db   *gorm.DB
	RedisDb     *redis.Client
}

func (d *Db) Open(c Config) error {
	//mysql
	mysql := c.Jewel.Mysql
	maxIdleConns := c.Jewel.Max_idle_conns
	maxOpenConns := c.Jewel.Max_Open_conns
	if maxIdleConns == 0 {
		maxIdleConns = 10
	}
	if maxOpenConns == 0 {
		maxIdleConns = 100
	}
	if mysql != "" {
		db, err := gorm.Open("mysql", mysql)
		if err != nil {
			return err
		}
		db.Debug()
		//db.LogMode(true)
		db.DB().SetMaxIdleConns(maxIdleConns)
		db.DB().SetMaxOpenConns(maxOpenConns)
		d.MysqlDb = db
		seelog.Info("db connection success")
	}
	//postgres

	postgres := c.Jewel.Postgres
	if postgres != "" {
		db, err := gorm.Open("postgres", postgres)
		if err != nil {
			return err
		}
		db.DB().SetMaxIdleConns(maxIdleConns)
		db.DB().SetMaxOpenConns(maxOpenConns)
		d.PostDb = db
		seelog.Info("db connection success")
	}
	//sqlite3
	sqlite3 := c.Jewel.Sqlite3
	if sqlite3 != "" {
		db, err := gorm.Open("sqlite3", sqlite3)
		if err != nil {
			return err
		}
		db.DB().SetMaxIdleConns(maxIdleConns)
		db.DB().SetMaxOpenConns(maxOpenConns)
		d.Sqlite3Db = db
		seelog.Info("db connection success")
	}
	//redis
	redisConfig := c.Jewel.Redis
	if redisConfig.Host != "" {
		d.RedisDb = redis.NewClient(&redis.Options{
			Addr:         redisConfig.Host,
			Password:     redisConfig.Password, // no password set
			DB:           redisConfig.Db,       // use default DB
			ReadTimeout:  30 * time.Second,
			DialTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		})
		pong, err := d.RedisDb.Ping().Result()
		if err != nil {
			return err
		}
		seelog.Info("redis ping result:" + pong)
		seelog.Info("db connection success")
	}
	return nil
}

func (d *Db) Close() {
	if d.MysqlDb != nil {
		d.MysqlDb.Close()
	}
	if d.PostDb != nil {
		d.PostDb.Close()
	}
	if d.Sqlite3Db != nil {
		d.Sqlite3Db.Close()
	}
	if d.RedisDb != nil {
		d.RedisDb.Close()
	}
}
