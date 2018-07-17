package context

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	//_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/go-redis/redis"
	"github.com/cihub/seelog"
	"github.com/streadway/amqp"
	"log"
)

type Db struct {
	MysqlDb     *gorm.DB
	PostDb      *gorm.DB
	SqlServerDb *gorm.DB
	Sqlite3Db   *gorm.DB
	RedisDb     *redis.Client
	AmqpConnect *amqp.Connection
}

func (d *Db) Open(c Config) error {
	//mysql
	mysql := c.Jewel.Mysql
	maxIdleConns := c.Jewel.Max_Idle_Conns
	maxOpenConns := c.Jewel.Max_Open_Conns
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
		db.LogMode(c.Jewel.SqlShow)
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
			Addr:     redisConfig.Host,
			Password: redisConfig.Password, // no password set
			DB:       redisConfig.Db,       // use default DB
		})
		pong, err := d.RedisDb.Ping().Result()
		if err != nil {
			return err
		}
		seelog.Info("redis ping result:" + pong)
		seelog.Info("db connection success")
	}
	amqpConfig := c.Jewel.Amqp
	if amqpConfig != "" {
		conn, err := amqp.Dial(amqpConfig)
		if err != nil {
			log.Fatal(err)
		}
		d.AmqpConnect = conn
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
	if d.AmqpConnect != nil {
		d.AmqpConnect.Close()
	}
}
