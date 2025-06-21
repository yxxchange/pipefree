package dal

import (
	"github.com/spf13/viper"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/infra/dal/dao"
	"github.com/yxxchange/pipefree/infra/dal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"time"
)

const DsnTemplate = "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=UTC"

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(renderDsn()), renderConfig())
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	setConnPool(db, 30, 100, 10*time.Minute)

	log.Info("MySQL 数据库连接成功！")

	afterInitDB(db) // Perform post-initialization tasks
	return db
}

func setConnPool(db *gorm.DB, maxIdles, maxOpens int, connMaxLifetime time.Duration) {
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(maxIdles)           // 设置空闲连接池最大连接数
	sqlDB.SetMaxOpenConns(maxOpens)           // 设置数据库最大连接数
	sqlDB.SetConnMaxLifetime(connMaxLifetime) // 设置连接最大复用时间
	if err := sqlDB.Ping(); err != nil {
		panic("failed to ping database: " + err.Error())
	}
	return
}

func renderConfig() *gorm.Config {
	return &gorm.Config{
		Logger: customLogger(),
	}
}

func customLogger() logger.Interface {
	return logger.New(
		log.AsGormLoggerPlugin(), // 使用自定义的日志插件
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // 日志级别（Debug 会打印所有 SQL）
			Colorful:      true,        // 启用颜色
		},
	)
}

func renderDsn() string {
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	dbName := viper.GetString("mysql.database")
	dsn := strings.Replace(DsnTemplate, "user", user, 1)
	dsn = strings.Replace(dsn, "password", password, 1)
	dsn = strings.Replace(dsn, "host", host, 1)
	dsn = strings.Replace(dsn, "port", port, 1)
	dsn = strings.Replace(dsn, "dbname", dbName, 1)
	return dsn
}

func afterInitDB(db *gorm.DB) {
	migrateIfNeeded(db) // Perform migration if needed
	dao.SetDefault(db)  // Set the default DAO with the initialized DB
}

func migrateIfNeeded(db *gorm.DB) {
	if !viper.GetBool("mysql.migrate") {
		log.Info("Skipping database migration as per configuration.")
		return
	}
	log.Info("Starting database migration...")
	// Here you would typically call your migration logic, e.g.:
	migrateList := []interface{}{
		model.PipeCfg{},
		model.PipeExec{},
		model.PipeVersion{},
		model.NodeCfg{},
		model.NodeExec{},
	}
	err := db.AutoMigrate(migrateList...) // Perform migrations for the specified models
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}
	log.Info("Database migration completed successfully.")
}
