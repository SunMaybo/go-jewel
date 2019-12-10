package context

import (
	"github.com/jinzhu/gorm"
	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2"
	"github.com/SunMaybo/jewel-template/rest"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/SunMaybo/jewel-inject/inject"
	"go.uber.org/zap"
)

type BasePlugin struct {
	MysqlDb      map[string]*gorm.DB
	PostDb       map[string]*gorm.DB
	RedisDb      map[string]*redis.Client
	MgoDb        map[string]*mgo.Database
	RestTemplate map[string]*rest.RestTemplate
}

func NewBasePlugin() *BasePlugin {
	return &BasePlugin{
		MysqlDb:      make(map[string]*gorm.DB),
		PostDb:       make(map[string]*gorm.DB),
		RedisDb:      make(map[string]*redis.Client),
		MgoDb:        make(map[string]*mgo.Database),
		RestTemplate: make(map[string]*rest.RestTemplate),
	}
}

func (d *BasePlugin) DefaultMysql() *gorm.DB {
	return d.MysqlDb["default"]
}
func (d *BasePlugin) Mysql(name string) *gorm.DB {
	return d.MysqlDb[name]
}

func (d *BasePlugin) DefaultPost() *gorm.DB {
	return d.PostDb["default"]
}
func (d *BasePlugin) Post(name string) *gorm.DB {
	return d.PostDb[name]
}
func (d *BasePlugin) DefaultRedis() *redis.Client {
	return d.RedisDb["default"]
}
func (d *BasePlugin) Redis(name string) *redis.Client {
	return d.RedisDb[name]
}
func (d *BasePlugin) DefaultMgo() *mgo.Database {
	return d.MgoDb["default"]
}
func (d *BasePlugin) Mgo(name string) *mgo.Database {
	return d.MgoDb[name]
}

func (d BasePlugin) DefaultRest() *rest.RestTemplate {
	return d.RestTemplate["default"]
}
func (d *BasePlugin) Rest(name string) *rest.RestTemplate {
	return d.RestTemplate[name]
}

func (d *BasePlugin) Open(injector *inject.Injector) error {
	var jewel JewelProperties
	jewel = injector.Service(&jewel).(JewelProperties)
	mysql := jewel.Jewel.MySql
	if mysql != nil {
		for name, sqlDataSource := range mysql {
			if sqlDataSource.Enabled != nil && !*sqlDataSource.Enabled {
				continue
			}
			db, err := sqlDataSource.Create("mysql")
			if err != nil {
				return err
			}
			d.MysqlDb[name] = db
			zap.S().Infof("mysql connection success,dataSource:%s", name)
		}
	}

	redis := jewel.Jewel.Redis
	if redis != nil {
		for name, redisDataSource := range redis {
			if redisDataSource.Enabled != nil && !*redisDataSource.Enabled {
				continue
			}
			client, err := redisDataSource.Create()
			if err != nil {
				return err
			}
			d.RedisDb[name] = client
			zap.S().Infof("redis connection success,dataSource:%s", name)
		}
	}

	mgo := jewel.Jewel.Mgo
	if mgo != nil {
		for name, mgoDataSource := range mgo {
			if mgoDataSource.Enabled != nil && !*mgoDataSource.Enabled {
				continue
			}
			db, err := mgoDataSource.Create()
			if err != nil {
				return err
			}
			d.MgoDb[name] = db
			zap.S().Infof("mgo connection success,dataSource:%s", name)
		}
	}

	postgres := jewel.Jewel.Postgres
	if postgres != nil {
		for name, postgresDataSource := range postgres {
			if postgresDataSource.Enabled != nil && !*postgresDataSource.Enabled {
				continue
			}
			db, err := postgresDataSource.Create("postgres")
			if err != nil {
				return err
			}
			d.PostDb[name] = db
			zap.S().Infof("postgres connection success,dataSource:%s", name)
		}
	}

	rests := jewel.Jewel.Rest
	if rests != nil {
		for name, restDataSource := range rests {
			if restDataSource.Enabled != nil && !*restDataSource.Enabled {
				continue
			}
			restTemplate, err := restDataSource.Create()
			if err != nil {
				return err
			}
			d.RestTemplate[name] = restTemplate
			zap.S().Infof("rest create success,templateName:%s", name)
		}
	}

	return nil
}

func (d *BasePlugin) Health() error {
	if d.MysqlDb != nil {
		for name, db := range d.MysqlDb {
			err := db.Exec("select 1").Error
			if err != nil {
				zap.S().Error("mysql db health error:" + name)
				return err
			}
		}
	}
	if d.RedisDb != nil {
		for name, db := range d.RedisDb {
			_, err := db.Ping().Result()
			if err != nil {
				zap.S().Error("redis client health error:" + name)
				return err
			}
		}
	}
	if d.MgoDb != nil {
		for name, db := range d.MgoDb {
			_, err := db.CollectionNames()
			if err != nil {
				zap.S().Error("mgo session health error:" + name)
				return err
			}
		}
	}
	if d.PostDb != nil {
		for name, db := range d.PostDb {
			err := db.Exec("select 1").Error
			if err != nil {
				zap.S().Error("postgres db health error:" + name)
				return err
			}
		}
	}
	return nil
}

func (d *BasePlugin) Close() {
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
func (d *BasePlugin) InterfaceName() string {
	return "base_plugin"
}
