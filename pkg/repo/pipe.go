package repo

import (
	"context"
	"github.com/yxxchange/pipefree/pkg/infra/mongoDB"
	"github.com/yxxchange/pipefree/pkg/infra/nebula"
	"github.com/yxxchange/pipefree/pkg/pipe/model"
	"go.mongodb.org/mongo-driver/bson"
)

var PipeRepo Pipe

const (
	PipeCfgCollection      = "pipe_cfg"
	PipeExecCollection     = "pipe_exec"
	NodeSnapshotCollection = "node_snapshot"
	PipeDBName             = "pipe"

	NebulaPipeExecSpace     = "pipe_exec"
	NebulaPipeExecBasicTag  = "pipe_exec_basic_tag"
	NebulaPipeExecBasicEdge = "pipe_exec_basic_edge"
)

type Pipe struct{}

func (p Pipe) CreatePipeCfg(ctx context.Context, pipe *model.PipeConfig) (id interface{}, err error) {
	db := mongoDB.AssignDB(PipeDBName, PipeCfgCollection)
	res, err := db.InsertOne(ctx, pipe)
	id = res.InsertedID
	return id, err
}

func (p Pipe) GetPipeCfg(ctx context.Context, id string) (pipe model.PipeConfig, err error) {
	db := mongoDB.AssignDB(PipeDBName, PipeCfgCollection)
	err = db.FindOne(ctx, bson.M{"_id": id}).Decode(&pipe)
	return pipe, err
}

func (p Pipe) UpdatePipeCfg(ctx context.Context, filter, fields bson.M) (err error) {
	db := mongoDB.AssignDB(PipeDBName, PipeCfgCollection)
	_, err = db.UpdateOne(ctx, filter, fields)
	return
}

func (p Pipe) CreatePipeExec(ctx context.Context, exec *model.PipeExec) (id interface{}, err error) {
	db := mongoDB.AssignDB(PipeDBName, PipeExecCollection)
	res, err := db.InsertOne(ctx, exec)
	id = res.InsertedID
	return id, err
}

func (p Pipe) BatchCreateNodeSnapshot(ctx context.Context, nodes []interface{}) (ids []interface{}, err error) {
	db := mongoDB.AssignDB(PipeDBName, NodeSnapshotCollection)
	res, err := db.InsertMany(ctx, nodes)
	return res.InsertedIDs, err
}

func (p Pipe) CreatePipeExecVertex(vertex interface{}, ifNotExist ...bool) error {
	return nebula.Use(NebulaPipeExecSpace).InsertVertex(vertex, ifNotExist...).Exec()
}

func (p Pipe) CreatePipeExecEdge(edge interface{}, ifNotExist ...bool) error {
	return nebula.Use(NebulaPipeExecSpace).InsertEdge(edge, ifNotExist...).Exec()
}
