package context

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"encoding/xml"
	"strings"
	"sync"
	"encoding/json"
	"log"
	"github.com/jinzhu/gorm"
	"time"
	"github.com/go-redis/redis"
	"crypto/tls"
	"gopkg.in/mgo.v2"
	"github.com/SunMaybo/jewel-template/template/rest"
	"net/http"
	"strconv"
)

type SqlDataSource struct {
	ConnMaxLifeTime *int    `json:"conn_max_life_time" yaml:"conn_max_life_time" xml:"conn_max_life_time"`
	MaxIdleConns    *int    `json:"max_idle_conns" yaml:"max_idle_conns" xml:"max_idle_conns"`
	MaxOpenConns    *int    `json:"max_open_conns" yaml:"max_open_conns" xml:"max_open_conns"`
	SqlShow         *bool   `json:"sql_show" yaml:"show_sql" xml:"show_sql"`
	URL             *string `json:"url" yaml:"url" xml:"url"`
	Enabled         *bool   `json:"enabled" yaml:"enabled" xml:"enabled"`
}

func (ds *SqlDataSource) Create(name string) (*gorm.DB, error) {
	db, err := gorm.Open(name, *ds.URL)
	if err != nil {
		return nil, err
	}
	if ds.SqlShow != nil {
		db.LogMode(*ds.SqlShow)
	}
	if ds.ConnMaxLifeTime != nil {
		db.DB().SetConnMaxLifetime(time.Duration(*ds.ConnMaxLifeTime * 1000000))
	}
	if ds.MaxIdleConns != nil {
		db.DB().SetMaxIdleConns(*ds.MaxIdleConns)
	}
	if ds.MaxOpenConns != nil {
		db.DB().SetMaxOpenConns(*ds.MaxOpenConns)
	}
	return db, nil
}

type RedisDataSource struct {
	Addr               *string `json:"addr" yaml:"addr" xml:"addr"`
	Db                 *int    `json:"db" yaml:"db" xml:"db"`
	DialTimeout        *int    `json:"dial_timeout" yaml:"dial_timeout" xml:"dial_timeout"`
	IdleCheckFrequency *int    `json:"idle_check_frequency" yaml:"idle_check_frequency" xml:"idle_check_frequency"`
	IdleTimeout        *int    `json:"idle_timeout" yaml:"idle_timeout" xml:"idle_timeout"`
	MaxRetries         *int    `json:"max_retries" yaml:"max_retries" xml:"max_retries"`
	MaxRetryBackOff    *int    `json:"max_retry_back_off" yaml:"max_retry_back_off" xml:"max_retry_back_off"`
	MinRetryBackOff    *int    `json:"min_retry_back_off" yaml:"min_retry_back_off" xml:"min_retry_back_off"`
	Network            *string `json:"network" yaml:"network" xml:"network"`
	Password           *string `json:"password" yaml:"password" xml:"password"`
	PoolSize           *int    `json:"pool_size" yaml:"pool_size" xml:"pool_size"`
	PoolTimeout        *int    `json:"pool_timeout" yaml:"pool_timeout" xml:"pool_timeout"`
	ReadTimeout        *int    `json:"read_timeout" yaml:"read_timeout" xml:"read_timeout"`
	TLS                *bool   `json:"tls" yaml:"tls" xml:"tls"`
	WriteTimeout       *int    `json:"write_timeout" yaml:"write_timeout" xml:"write_timeout"`
	Enabled            *bool   `json:"enabled" yaml:"enabled" xml:"enabled"`
}

func (ds *RedisDataSource) Create() (*redis.Client, error) {
	option := &redis.Options{}
	if ds.Addr != nil {
		option.Addr = *ds.Addr
	}
	if ds.Db != nil {
		option.DB = *ds.Db
	}
	if ds.WriteTimeout != nil {
		option.WriteTimeout = time.Duration(*ds.WriteTimeout * 1000000)
	}
	if ds.DialTimeout != nil {
		option.DialTimeout = time.Duration(*ds.DialTimeout * 1000000)
	}
	if ds.ReadTimeout != nil {
		option.ReadTimeout = time.Duration(*ds.ReadTimeout * 1000000)
	}
	if ds.IdleCheckFrequency != nil {
		option.IdleCheckFrequency = time.Duration(*ds.IdleCheckFrequency)
	}
	if ds.IdleTimeout != nil {
		option.IdleTimeout = time.Duration(*ds.IdleTimeout * 1000000)
	}
	if ds.MaxRetries != nil {
		option.MaxRetries = *ds.MaxRetries
	}
	if ds.MaxRetryBackOff != nil {
		option.MaxRetryBackoff = time.Duration(*ds.MaxRetryBackOff * 1000000)
	}
	if ds.MinRetryBackOff != nil {
		option.MinRetryBackoff = time.Duration(*ds.MinRetryBackOff * 1000000)
	}
	if ds.Network != nil {
		option.Network = *ds.Network
	}
	if ds.Password != nil {
		option.Password = *ds.Password
	}
	if ds.PoolSize != nil {
		option.PoolSize = *ds.PoolSize
	}
	if ds.PoolTimeout != nil {
		option.PoolTimeout = time.Duration(*ds.PoolTimeout * 1000000)
	}
	if ds.TLS != nil {
		option.TLSConfig = &tls.Config{
			InsecureSkipVerify: *ds.TLS,
		}
	}
	db := redis.NewClient(option)
	_, err := db.Ping().Result()
	if err != nil {
		return nil, err
	}
	return db, nil
}

type MgoDataSource struct {
	Address        *string `json:"address" yaml:"address" xml:"address"`
	Database       *string `json:"database" yaml:"database" xml:"database"`
	Direct         *bool   `json:"direct" yaml:"direct" xml:"direct"`
	FailFast       *bool   `json:"fail_fast" yaml:"fail_fast" xml:"fail_fast"`
	Password       *string `json:"password" yaml:"password" xml:"password"`
	PoolLimit      *int    `json:"pool_limit" yaml:"pool_limit" xml:"pool_limit"`
	ReplicaSetName *string `json:"replica_set_name" yaml:"replica_set_name" xml:"replica_set_name"`
	Service        *string `json:"service" yaml:"service" xml:"service"`
	ServiceHost    *string `json:"service_host" yaml:"service_host" xml:"service_host"`
	Source         *string `json:"source" yaml:"source" xml:"source"`
	Timeout        *int    `json:"timeout" yaml:"timeout" xml:"timeout"`
	UserName       *string `json:"user_name" yaml:"user_name" xml:"user_name"`
	Enabled        *bool   `json:"enabled" yaml:"enabled" xml:"enabled"`
}

func (ds *MgoDataSource) Create() (*mgo.Database, error) {
	info := &mgo.DialInfo{}
	if ds.Address != nil {
		info.Addrs = strings.Split(*ds.Address, ",")
	}
	if ds.Password != nil {
		info.Password = *ds.Password
	}
	if ds.Service != nil {
		info.Service = *ds.Service
	}
	if ds.Database != nil {
		info.Database = *ds.Database
	}
	if ds.Direct != nil {
		info.Direct = *ds.Direct
	}
	if ds.FailFast != nil {
		info.FailFast = *ds.FailFast
	}
	if ds.PoolLimit != nil {
		info.PoolLimit = *ds.PoolLimit
	}
	if ds.ReplicaSetName != nil {
		info.ReplicaSetName = *ds.ReplicaSetName
	}
	if ds.ServiceHost != nil {
		info.ServiceHost = *ds.ServiceHost
	}
	if ds.Source != nil {
		info.Source = *ds.Source
	}
	if ds.Timeout != nil {
		info.Timeout = time.Duration(*ds.Timeout * 1000000)
	}
	if ds.UserName != nil {
		info.Username = *ds.UserName
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		return nil, err
	}
	db := session.DB(info.Database)
	return db, nil
}

type RestProperties struct {
	Authorization      *string `json:"authorization" yaml:"authorization" xml:"authorization"`
	DisableCompression *bool   `json:"disable_compression" yaml:"disable_compression" xml:"disable_compression"`
	IdleConnTimeout    *int    `json:"idle_conn_timeout" yaml:"idle_conn_timeout" xml:"idle_conn_timeout"`
	MaxIdleConns       *int    `json:"max_idle_conns" yaml:"max_idle_conns" xml:"max_idle_conns"`
	ReplyCount         *int    `json:"reply_count" yaml:"reply_count" xml:"reply_count"`
	SocketTimeout      *int    `json:"socket_timeout" yaml:"socket_timeout" xml:"socket_timeout"`
	Enabled            *bool   `json:"enabled" yaml:"enabled" xml:"enabled"`
}

type Manager struct {
	User     *string `json:"user" yaml:"user" xml:"user"`
	Password *string `json:"password" yaml:"password" xml:"password"`
	Enabled  *bool   `json:"enabled" yaml:"enabled" xml:"enabled"`
}

type ServerProperties struct {
	GinMode *string `json:"gin_mode" yaml:"gin_mode" xml:"gin_mode"`

	ContextPath *string `json:"context_path" yaml:"context_path" xml:"context_path"`

	Port *int64 `json:"port" yaml:"port" xml:"port"`

	Templates *string `json:"templates" yaml:"templates" xml:"templates"`

	EnableTls *bool `json:"enable_tls" yaml:"enable_tls" xml:"enable_tls"`

	X509KeyPairServerPem *string `json:"x509_key_pair_server_pem" yaml:"x509_key_pair_server_pem" xml:"x509_key_pair_server_pem"`

	X509KeyPairServerKey *string `json:"x509_key_pair_server_key" yaml:"x509_key_pair_server_key" xml:"x509_key_pair_server_key"`

	ReadTimeout *int64 `json:"read_timeout" yaml:"read_timeout" xml:"read_timeout"`

	ReadHeaderTimeout *int64 `json:"read_header_timeout" yaml:"read_header_timeout" xml:"read_header_timeout"`

	WriteTimeout *int64 `json:"write_timeout" yaml:"write_timeout" xml:"write_timeout"`

	IdleTimeout *int64 `json:"idle_timeout" yaml:"idle_timeout" xml:"idle_timeout"`

	MaxHeaderBytes *int `json:"max_header_bytes" yaml:"max_header_bytes" xml:"max_header_bytes"`

	Manager Manager `json:"manager" yaml:"manager" xml:"max_header_bytes"`
}

func (serverProperties *ServerProperties) Create() (*http.Server, error) {
	server := &http.Server{}
	if serverProperties.Port != nil {
		server.Addr = ":" + strconv.FormatInt(*serverProperties.Port, 10)
	}
	if serverProperties.IdleTimeout != nil {
		server.IdleTimeout = time.Duration(*serverProperties.IdleTimeout * 1000000)
	}
	if serverProperties.ReadTimeout != nil {
		server.ReadTimeout = time.Duration(*serverProperties.ReadTimeout * 1000000)
	}
	if serverProperties.WriteTimeout != nil {
		server.WriteTimeout = time.Duration(*serverProperties.WriteTimeout * 1000000)
	}
	if serverProperties.ReadHeaderTimeout != nil {
		server.ReadHeaderTimeout = time.Duration(*serverProperties.ReadHeaderTimeout * 1000000)
	}
	if serverProperties.MaxHeaderBytes != nil {
		server.MaxHeaderBytes = *serverProperties.MaxHeaderBytes
	}
	if serverProperties.EnableTls != nil && *serverProperties.EnableTls {
		crt, err := tls.LoadX509KeyPair(*serverProperties.X509KeyPairServerPem, *serverProperties.X509KeyPairServerKey)
		if err != nil {
			return server, err
		}
		tlsConfig := &tls.Config{}
		tlsConfig.Certificates = []tls.Certificate{crt}
		// Time returns the current time as the number of seconds since the epoch.
		// If Time is nil, TLS uses time.Now.
		tlsConfig.Time = time.Now
		server.TLSConfig = tlsConfig
	}
	return server, nil
}

func (restOptions *RestProperties) Create() (*rest.RestTemplate, error) {
	config := rest.ClientConfig{}
	if restOptions.MaxIdleConns != nil {
		config.MaxIdleConns = *restOptions.MaxIdleConns
	} else {
		config.MaxIdleConns = 2
	}
	if restOptions.ReplyCount != nil {
		config.ReplyCount = *restOptions.ReplyCount
	} else {
		config.ReplyCount = 3
	}
	if restOptions.SocketTimeout != nil {
		config.SocketTimeout = time.Duration(*restOptions.SocketTimeout) * time.Millisecond
	} else {
		config.SocketTimeout = 3 * time.Second
	}
	if restOptions.IdleConnTimeout != nil {
		config.IdleConnTimeout = time.Duration(*restOptions.IdleConnTimeout) * time.Millisecond
	} else {
		config.IdleConnTimeout = 10 * time.Millisecond
	}
	if restOptions.Authorization != nil {
		config.Authorization = *restOptions.Authorization
	}
	if restOptions.DisableCompression != nil {
		config.DisableCompression = *restOptions.DisableCompression
	}
	return rest.Config(config), nil
}

type JewelProperties struct {
	Jewel struct {
		Name string `json:"name" yaml:"name" xml:"name"`
		Log struct {
			Level string `json:"level" yaml:"level" xml:"level"`
		} `json:"log" yaml:"log" xml:"log"`
		Server ServerProperties `json:"server" yaml:"server" xml:"server"`
		Profiles struct {
			Active string `yaml:"active" xml:"active" json:"active"`
		} `json:"profiles" yaml:"profiles" xml:"profiles"`
		MySql    map[string]SqlDataSource   `json:"mysql" yaml:"mysql" xml:"mysql"`
		Postgres map[string]SqlDataSource   `json:"postgres" yaml:"postgres" xml:"postgres"`
		Redis    map[string]RedisDataSource `json:"redis" yaml:"redis" xml:"redis"`
		Mgo      map[string]MgoDataSource   `json:"mgo" yaml:"mgo" xml:"mgo"`
		Rest     map[string]RestProperties  `json:"rest" yaml:"rest" xml:"rest"`
	} `json:"jewel" yaml:"jewel" xml:"jewel"`
}
type Properties struct {
	locker sync.Mutex
}

func NewProperties() *Properties {
	return &Properties{
		locker: sync.Mutex{},
	}
}
func (prop *Properties) Load(fileName string, inter interface{}) {
	prop.locker.Lock()
	defer prop.locker.Unlock()
	if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
		prop.loadYml(fileName, inter)
	} else if strings.HasSuffix(fileName, ".xml") {
		prop.loadXml(fileName, inter)
	} else {
		prop.loadJson(fileName, inter)
	}
}
func (prop *Properties) loadYml(fileName string, inter interface{}) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(buff, inter)
	if err != nil {
		log.Fatalln(err)
	}

}
func (prop *Properties) loadXml(fileName string, inter interface{}) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	err = xml.Unmarshal(buff, &inter)
	if err != nil {
		log.Fatalln(err)
	}
}
func (prop *Properties) loadJson(fileName string, inter interface{}) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(buff, &inter)
	if err != nil {
		log.Fatalln(err)
	}
}
