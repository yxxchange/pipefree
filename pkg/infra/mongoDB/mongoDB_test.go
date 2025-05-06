package mongoDB

import (
	"context"
	"github.com/yxxchange/pipefree/config"
	"testing"
)

func TestInitMongoDB(t *testing.T) {
	config.InitConfig("../../../config.yaml")
	Init()
	res, err := AssignDB("test", "users").InsertOne(context.Background(), map[string]interface{}{
		"name": "Donald Trump",
		"age":  18,
	})
	if err != nil {
		t.Errorf("Failed to insert document: %v", err)
	} else {
		t.Logf("Inserted document with ID: %v", res.InsertedID)
	}
}
