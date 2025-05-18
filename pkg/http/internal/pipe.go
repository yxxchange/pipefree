package internal

import (
	"context"
	"fmt"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/pkg/infra/etcd"
	"github.com/yxxchange/pipefree/pkg/pipe/model"
	"github.com/yxxchange/pipefree/pkg/pipe/orca"
	"github.com/yxxchange/pipefree/pkg/repo"
)

func CreatePipe(ctx context.Context, pipe model.PipeConfig) error {
	err := pipe.PipeFlow.ValidateStaticCfg()
	if err != nil {
		log.Errorf("validate pipe error: %v", err)
		return err
	}
	parser := orca.NewGraphParser()
	err = parser.Parse(pipe.PipeFlow).IsValid()
	if err != nil {
		log.Errorf("validate graph error: %v", err)
		return err
	}

	id, err := repo.PipeRepo.CreatePipeCfg(ctx, &pipe)
	if err != nil {
		return err
	}
	log.Infof("create pipe cfg: %v", id)
	return nil
}

func RunPipe(ctx context.Context, id string) error {
	pipe, err := repo.PipeRepo.GetPipeCfg(ctx, id)
	if err != nil {
		return err
	}
	err = pipe.PipeFlow.ValidateStaticCfg()
	if err != nil {
		log.Errorf("validate pipe error: %v", err)
		return err
	}
	pipeExec := pipe.ToPipeExec()
	err = pipeExec.PipeFlow.ValidateDynamicCfg()
	if err != nil {
		log.Errorf("validate pipe exec error: %v", err)
		return err
	}
	parser := orca.NewGraphParser()
	err = parser.Parse(pipeExec.PipeFlow).IsValid()
	if err != nil {
		log.Errorf("validate graph error: %v", err)
		return err
	}

	err = savePipeExec(ctx, &pipeExec)
	if err != nil {
		return fmt.Errorf("save pipe exec doc error: %v", err)
	}
	err = saveGraph(ctx, &pipeExec)
	if err != nil {
		log.Errorf("save graph error: %v", err)
		return err
	}
	origin, err := parser.FindTheOrigin()
	if err != nil {
		log.Errorf("find the origin error: %v", err)
		return err
	}
	return etcd.Put(ctx, origin.ToIdentifier().Identifier(), origin.ToString())
}

func savePipeExec(ctx context.Context, pipeExec *model.PipeExec) error {
	execId, err := repo.PipeRepo.CreatePipeExec(ctx, pipeExec)
	if err != nil {
		log.Errorf("create pipe exec error: %v", err)
		return err
	}
	log.Infof("create pipe exec: %v", execId)
	return nil
}

func saveGraph(_ context.Context, graph *model.PipeExec) error {
	for _, vertex := range graph.Graph.Vertexes {
		err := repo.PipeRepo.CreatePipeExecVertex(vertex)
		if err != nil {
			log.Errorf("create pipe exec vertex error: %v", err)
			return err
		}
	}
	for _, edge := range graph.Graph.Edges {
		err := repo.PipeRepo.CreatePipeExecEdge(edge)
		if err != nil {
			log.Errorf("create pipe exec edge error: %v", err)
			return err
		}
	}
	return nil
}
