package pipe_exec

import (
	"context"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/infra/dal/dao"
	"github.com/yxxchange/pipefree/infra/dal/model"
	"github.com/yxxchange/pipefree/infra/etcd"
)

const ErrorCode = 10002

type Service struct {
	pipeExec dao.IPipeExecDo
	nodeExec dao.INodeExecDo
	pipeCfg  dao.IPipeCfgDo
	nodeCfg  dao.INodeCfgDo
	*dao.Query
	ctx context.Context
}

func NewService(ctx context.Context) *Service {
	return &Service{
		pipeExec: dao.Q.PipeExec.WithContext(ctx),
		nodeExec: dao.Q.NodeExec.WithContext(ctx),
		pipeCfg:  dao.Q.PipeCfg.WithContext(ctx),
		nodeCfg:  dao.Q.NodeCfg.WithContext(ctx),
		Query:    dao.Q,
		ctx:      ctx,
	}
}

func (s *Service) Run(pipeId int64) error {
	pipeCfg, err := s.pipeCfg.Where(dao.PipeCfg.Id.Eq(pipeId)).First()
	if err != nil {
		log.Errorf("run pipe execution failed, get pipe configuration failed: %v", err)
		return err
	}
	nodeCfgList, err := s.nodeCfg.Where(dao.NodeCfg.PipeCfgId.Eq(pipeId)).Find()
	if err != nil {
		log.Errorf("run pipe execution failed, get node configurations failed: %v", err)
	}
	err = s.Query.Transaction(func(tx *dao.Query) error {
		return run(tx, s.ctx, pipeCfg, nodeCfgList)
	})
	if err != nil {
		log.Errorf("run pipe execution failed, transaction failed: %v", err)
		return err
	}
	return nil
}
func run(tx *dao.Query, ctx context.Context, pipeCfg *model.PipeCfg, nodeCfgList []*model.NodeCfg) error {
	pipeExec := model.NewPipeExec(pipeCfg)
	err := tx.PipeExec.WithContext(ctx).Create(pipeExec)
	if err != nil {
		log.Errorf("run pipe execution failed, create pipe execution failed: %v", err)
		return err
	}
	nodeExecList := make([]*model.NodeExec, 0, len(nodeCfgList))
	for _, nodeCfg := range nodeCfgList {
		nodeExec := model.NewNodeExec(nodeCfg, pipeExec)
		nodeExecList = append(nodeExecList, nodeExec)
	}
	err = tx.NodeExec.WithContext(ctx).Create(nodeExecList...)
	if err != nil {
		log.Errorf("run pipe execution failed, create node executions failed: %v", err)
		return err
	}
	initialNodes := make([]*model.NodeExec, 0)
	for _, nodeExec := range nodeExecList {
		if nodeExec.InDegree == 0 {
			initialNodes = append(initialNodes, nodeExec)
		}
	}
	kv := make(map[string]string, len(initialNodes))
	for _, nodeExec := range initialNodes {
		if val, e := ValueGen(nodeExec); e != nil {
			log.Errorf("run pipe execution failed, generate value for node %s failed: %v", nodeExec.Name, e)
			return e
		} else {
			kv[KeyGen(nodeExec)] = val
		}
	}
	err = etcd.TransactionPut(ctx, kv)
	if err != nil {
		log.Errorf("run pipe execution failed, transaction put to etcd failed: %v", err)
		return err
	}
	return nil
}
