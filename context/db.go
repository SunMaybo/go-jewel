package context

import (
	"github.com/jinzhu/gorm"
	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2"
	"github.com/SunMaybo/jewel-template/template/rest"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/cihub/seelog"
)

type Db struct {
	MysqlDb      map[string]*gorm.DB
	PostDb       map[string]*gorm.DB
	RedisDb      map[string]*redis.Client
	MgoDb        map[string]*mgo.Database
	RestTemplate map[string]*rest.RestTemplate
}

func NewDb() *Db {
	return &Db{
		MysqlDb:      make(map[string]*gorm.DB),
		PostDb:       make(map[string]*gorm.DB),
		RedisDb:      make(map[string]*redis.Client),
		MgoDb:        make(map[string]*mgo.Database),
		RestTemplate: make(map[string]*rest.RestTemplate),
	}
}

func (d *Db) DefaultMysql() *gorm.DB {
	return d.MysqlDb["default"]
}
func (d *Db) Mysql(name string) *gorm.DB {
	return d.MysqlDb[name]
}

func (d *Db) DefaultPost() *gorm.DB {
	return d.PostDb["default"]
}
func (d *Db) Post(name string) *gorm.DB {
	return d.PostDb[name]
}
func (d *Db) DefaultRedis() *redis.Client {
	return d.RedisDb["default"]
}
func (d *Db) Redis(name string) *redis.Client {
	return d.RedisDb[name]
}
func (d *Db) DefaultMgo() *mgo.Database {
	return d.MgoDb["default"]
}
func (d *Db) Mgo(name string) *mgo.Database {
	return d.MgoDb[name]
}

func (d *Db) DefaultRest() *rest.RestTemplate {
	return d.RestTemplate["default"]
}
func (d *Db) Rest(name string) *rest.RestTemplate {
	return d.RestTemplate[name]
}

func (d *Db) Open(jewel JewelProperties) error {
	mysql := jewel.Jewel.MySql
	if mysql != nil {
		for name, sqlDataSource := range mysql {
			db, err := sqlDataSource.Create("mysql")
			if err != nil {
				return err
			}
			d.MysqlDb[name] = db
		}
	}

	redis := jewel.Jewel.Redis
	if redis != nil {
		for name, redisDataSource := range redis {
			client, err := redisDataSource.Create()
			if err != nil {
				return err
			}
			d.RedisDb[name] = client
		}
	}

	mgo := jewel.Jewel.Mgo
	if mgo != nil {
		for name, mgoDataSource := range mgo {
			db, err := mgoDataSource.Create()
			if err != nil {
				return err
			}
			d.MgoDb[name] = db
		}
	}

	postgres := jewel.Jewel.Postgres
	if postgres != nil {
		for name, postgresDataSource := range postgres {
			db, err := postgresDataSource.Create("postgres")
			if err != nil {
				return err
			}
			d.PostDb[name] = db
		}
	}

	rests := jewel.Jewel.Rest
	if rests != nil {
		for name, restDataSource := range rests {
			restTemplate, err := restDataSource.Create()
			if err != nil {
				return err
			}
			d.RestTemplate[name] = restTemplate
		}
	}

	return nil
}

func (d *Db) Health() error {
	if d.MysqlDb != nil {
		for name, db := range d.MysqlDb {
			err := db.Exec("select 1").Error
			if err != nil {
				seelog.Error("mysql db health error:" + name)
				return err
			}
		}
	}
	if d.RedisDb != nil {
		for name, db := range d.RedisDb {
			_, err := db.Ping().Result()
			if err != nil {
				seelog.Error("redis client health error:" + name)
				return err
			}
		}
	}
	if d.MgoDb != nil {
		for name, db := range d.MgoDb {
			_, err := db.CollectionNames()
			if err != nil {
				seelog.Error("mgo session health error:" + name)
				return err
			}
		}
	}
	if d.PostDb != nil {
		for name, db := range d.PostDb {
			err := db.Exec("select 1").Error
			if err != nil {
				seelog.Error("postgres db health error:" + name)
				return err
			}
		}
	}
	return nil
}

func (d *Db) Close() {
	if d.MysqlDb != nil {
		for _, db := range d.MysqlDb {
			db.Close()
		}
	}
	if d.PostDb != nil {
		for _, db := range d.PostDb {
			db.Close()
		}
	}
	if d.RedisDb != nil {
		for _, client := range d.RedisDb {
			client.Close()
		}
	}
	if d.MgoDb != nil {
		for _, db := range d.MgoDb {
			db.Session.Close()
		}
	}
}
