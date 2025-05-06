package mongoDB

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/yxxchange/pipefree/helper/log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var once sync.Once

func Init() {
	once.Do(func() {
		initMongoDB()
	})
}

// 初始化 MongoDB 客户端
func initMongoDB() {
	// 设置连接超时和心跳间隔
	clientOptions := options.Client().
		ApplyURI(viper.GetString("mongoDB.uri")).
		SetTimeout(viper.GetDuration("mongoDB.options.timeout") * time.Second).
		SetMaxPoolSize(viper.GetUint64("mongoDB.options.maxPoolSize")).
		SetMinPoolSize(viper.GetUint64("mongoDB.options.minPoolSize")).
		SetMaxConnIdleTime(viper.GetDuration("mongoDB.options.maxConnIdleTime") * time.Second).
		SetHeartbeatInterval(10 * time.Second)

	if viper.GetBool("mongoDB.options.needAuth") {
		clientOptions.SetAuth(options.Credential{
			Username: viper.GetString("mongoDB.username"),
			Password: viper.GetString("mongoDB.password"),
		})
	}

	// 创建连接
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Errorf("Failed to connect to MongoDB: %v", err)
		panic(err)
	}

	// 检查连接是否成功
	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Errorf("Failed to ping MongoDB: %v", err)
		panic(err)
	}

	fmt.Println("Connected to MongoDB!")

	Client = client
}

func AssignDB(dbName, collName string) *mongo.Collection {
	return Client.Database(dbName).Collection(collName)
}

func Close() {
	if err := Client.Disconnect(context.TODO()); err != nil {
		log.Errorf("Failed to disconnect from MongoDB: %v", err)
	}
	log.Info("Disconnected from MongoDB!")
}
