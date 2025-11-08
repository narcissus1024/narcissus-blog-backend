package mysql

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Client *gorm.DB
	once   sync.Once
)

const (
	DB_TRANSACTION_CONTEXT_KEY = "tx_db"
)

type MysqlConfig struct {
	User     string `json:"user,omitempty" yaml:"user,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	Host     string `json:"host,omitempty" yaml:"host,omitempty"`
	Port     int    `json:"port,omitempty" yaml:"port,omitempty"`
	DBName   string `json:"dbname,omitempty" yaml:"dbname,omitempty"`
}

func NewDefaultMysqlCfg() MysqlConfig {
	return MysqlConfig{
		User:     "root",
		Password: "123456",
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   "blog",
	}
}

func MustInit(mysqlConf MysqlConfig) {
	once.Do(func() {
		Client = initMysqlClient(mysqlConf)
	})
}

func initMysqlClient(mysqlConf MysqlConfig) *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConf.User,
		mysqlConf.Password,
		mysqlConf.Host,
		mysqlConf.Port,
		mysqlConf.DBName)
	// todo check mysql config
	mysqlClient, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err.Error())
	}
	return mysqlClient
}

func RunDBTransaction(c *gin.Context, fn func() error) error {
	tx := Client.Begin()
	if tx.Error != nil {
		return fmt.Errorf("begin transaction failed: %w", tx.Error)
	}

	c.Set(DB_TRANSACTION_CONTEXT_KEY, tx)
	defer c.Set(DB_TRANSACTION_CONTEXT_KEY, nil)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func GetDBFromContext(c *gin.Context) *gorm.DB {
	tx, dbExist := c.Get(DB_TRANSACTION_CONTEXT_KEY)
	if !dbExist || tx == nil {
		return Client
	}
	txDB, ok := tx.(*gorm.DB)
	if !ok {
		return Client
	}
	return txDB
}

func GetDBFromContext2(ctx context.Context) *gorm.DB {
	tx := ctx.Value(DB_TRANSACTION_CONTEXT_KEY)
	if tx == nil {
		return Client
	}
	txDB, ok := tx.(*gorm.DB)
	if !ok {
		return Client
	}
	return txDB
}
