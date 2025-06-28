package pipe_cfg

import (
	"context"
	"fmt"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/infra/dal/dao"
	"github.com/yxxchange/pipefree/infra/dal/model"
	"github.com/yxxchange/pipefree/service/graph"
)

const ErrorCode = 10001

type Service struct {
	pipeCfg dao.IPipeCfgDo
	nodeCfg dao.INodeCfgDo
	*dao.Query
	ctx context.Context
}

func NewService(ctx context.Context) *Service {
	return &Service{
		pipeCfg: dao.Q.PipeCfg.WithContext(ctx),
		nodeCfg: dao.Q.NodeCfg.WithContext(ctx),
		Query:   dao.Q,
		ctx:     ctx,
	}
}

func (s *Service) GetById(pipeId int64) (*model.PipeCfg, error) {
	pipeCfg, err := s.pipeCfg.Where(dao.PipeCfg.Id.Eq(pipeId)).First()
	if err != nil {
		return nil, err
	}
	return pipeCfg, nil
}

func (s *Service) Create(pipe model.PipeCfg, nodes []model.NodeCfg) error {
	dag, err := graph.Extract(pipe.Graph)
	if err != nil {
		log.Errorf("create pipe configuration failed, extract graph failed: %v", err)
		return err
	}
	err = graph.IsDAG(dag)
	if err != nil {
		log.Errorf("create pipe configuration failed, graph is not a DAG: %v", err)
		return err
	}
	// Set the pipe configuration ID for each node
	nodeCfgMap := make(map[string]*model.NodeCfg)
	nodeCfgList := make([]*model.NodeCfg, 0)
	for _, node := range nodes {
		_, exists := nodeCfgMap[node.Name]
		if exists {
			return fmt.Errorf("create pipe configuration failed, duplicate node name: %s", node.Name)
		}
		_, exists = dag.VertexMap[node.Name]
		if !exists {
			return fmt.Errorf("create pipe configuration failed, node %s not found in graph", node.Name)
		}
		nodeCfgMap[node.Name] = &node
		nodeCfgList = append(nodeCfgList, &node)
	}

	for name, vertex := range dag.VertexMap {
		_, exists := nodeCfgMap[name]
		if !exists {
			return fmt.Errorf("create pipe configuration failed, node %s not found in node configurations", name)
		}
		nodeCfgMap[name].InDegree = vertex.InDegree
	}

	err = s.Transaction(func(tx *dao.Query) error {
		err = tx.PipeCfg.WithContext(s.ctx).Create(&pipe)
		if err != nil {
			return fmt.Errorf("create pipe configuration failed: %w", err)
		}
		for i := range nodeCfgList {
			nodeCfgList[i].PipeCfgId = pipe.Id
		}
		err = tx.NodeCfg.WithContext(s.ctx).CreateInBatches(nodeCfgList, 100)
		if err != nil {
			return fmt.Errorf("create pipe configuration failed, create node configurations failed: %w", err)
		}
		return nil
	})

	return err
}

func (s *Service) Update(pipeCfg *model.PipeCfg) error {
	if err := s.pipeCfg.Save(pipeCfg); err != nil {
		return err
	}
	return nil
}

func (s *Service) Delete(pipeId int64) error {
	_, err := s.pipeCfg.Where(dao.PipeCfg.Id.Eq(pipeId)).Delete()
	if err != nil {
		return err
	}
	log.Infof("Deleted pipe configuration with ID: %d", pipeId)
	return nil
}
