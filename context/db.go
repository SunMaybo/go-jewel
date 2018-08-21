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
	"os"
	"github.com/mgo"
)

type Db struct {
	MysqlDb     *gorm.DB
	PostDb      *gorm.DB
	SqlServerDb *gorm.DB
	Sqlite3Db   *gorm.DB
	RedisDb     *redis.Client
	AmqpConnect *amqp.Connection
	MgoDb       *mgo.Database
}

func (d *Db) Open(c Config) error {
	//mysql
	mysql := c.Jewel.Mysql
	mgoUrl := c.Jewel.Mgo
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
			seelog.Error("mysql connection failed......")
			seelog.Flush()
			os.Exit(-1)
			return err
		}
		db.LogMode(c.Jewel.SqlShow)
		db.DB().SetMaxIdleConns(maxIdleConns)
		db.DB().SetMaxOpenConns(maxOpenConns)
		d.MysqlDb = db
		seelog.Info("mysql connection success......")
	}
	if mgoUrl != "" {
		db, err := mgo.Dial(mgoUrl)
		if err != nil {
			seelog.Error("mgo connection failed......")
			seelog.Flush()
			os.Exit(-1)
			return err
		}

		d.MgoDb = db
		d.MgoDb.SetMode(mgo.Monotonic, true)
		seelog.Info("mgo connection success......")
	}
	//postgres

	postgres := c.Jewel.Postgres
	if postgres != "" {
		db, err := gorm.Open("postgres", postgres)
		if err != nil {
			seelog.Error("postgres connection failed......")
			seelog.Flush()
			os.Exit(-1)
			return err
		}
		db.DB().SetMaxIdleConns(maxIdleConns)
		db.DB().SetMaxOpenConns(maxOpenConns)
		d.PostDb = db
		seelog.Info("postgres connection success......")
	}
	//sqlite3
	sqlite3 := c.Jewel.Sqlite3
	if sqlite3 != "" {
		db, err := gorm.Open("sqlite3", sqlite3)
		if err != nil {
			seelog.Error("sqlite3 connection failed......")
			seelog.Flush()
			os.Exit(-1)
			return err
		}
		db.DB().SetMaxIdleConns(maxIdleConns)
		db.DB().SetMaxOpenConns(maxOpenConns)
		d.Sqlite3Db = db
		seelog.Info("sqlite3 connection success......")
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
			seelog.Error("redis connection failed......")
			seelog.Flush()
			os.Exit(-1)
			return err
		}
		seelog.Info("redis ping result:" + pong + "......")
		seelog.Info("redis connection success......")
	}
	amqpConfig := c.Jewel.Amqp
	if amqpConfig != "" {
		if c.Jewel.Amqp_Vhost == "" {
			c.Jewel.Amqp_Vhost = "/"
		}
		if c.Jewel.Amqp_Max_Channel <= 0 {
			c.Jewel.Amqp_Max_Channel = 10
		}
		amqpCfg := amqp.Config{
			Vhost:      c.Jewel.Amqp_Vhost,
			ChannelMax: c.Jewel.Amqp_Max_Channel,
		}
		conn, err := amqp.DialConfig(amqpConfig, amqpCfg)
		if err != nil {
			seelog.Error("amqp connection failed......")
			seelog.Flush()
			os.Exit(-1)
			return err
		}
		d.AmqpConnect = conn
		seelog.Info("amqp connection success......")
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
