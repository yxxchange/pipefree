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
	err := validatePipe(pipe)
	if err != nil {
		log.Errorf("validate pipe error: %v", err)
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
	err = validatePipe(pipe)
	if err != nil {
		log.Errorf("validate pipe error: %v", err)
		return err
	}
	pipeExec := pipe.ToPipeExec()
	// 1. store exec doc to mongoDB for persistence
	// 2. store exec node snapshot to mongoDB for quick query
	// 3. store exec node and edge to nebula for quick query
	// 4. find the root node and put to etcd for orchestration
	// step 1
	execId, err := repo.PipeRepo.CreatePipeExec(ctx, &pipeExec)
	if err != nil {
		log.Errorf("create pipe exec error: %v", err)
		return err
	}
	log.Infof("create pipe exec: %v", execId)
	// step 2
	pipeFragment := pipeExec.Decompose()
	_, err = repo.PipeRepo.BatchCreateNodeSnapshot(ctx, pipeFragment.NodeSnapshots)
	if err != nil {
		log.Errorf("create pipe exec node snapshot error: %v", err)
		return err
	}
	// step 3
	for _, vertex := range pipeFragment.Vertexes {
		err = repo.PipeRepo.CreatePipeExecVertex(vertex, true)
		if err != nil {
			log.Errorf("create pipe exec vertex error: %v", err)
			return err
		}
	}
	for _, edge := range pipeFragment.Edges {
		err = repo.PipeRepo.CreatePipeExecEdge(edge, true)
		if err != nil {
			log.Errorf("create pipe exec edge error: %v", err)
			return err
		}
	}
	// step 4
	origin, err := findTheOrigin(pipeExec, pipeFragment)
	if err != nil {
		return err
	}
	err = etcd.Put(ctx, origin.ToIdentifier().Identifier(), origin.ToString())
	if err != nil {
		log.Errorf("put origin to etcd error: %v", err)
		return err
	}
	// TODO： update runtime uuid and resource version
	return nil
}

func validatePipe(pipe model.PipeConfig) error {
	_, err := orca.NewGraphBuilder().ProcessPipeCfg(pipe).ProcessGraph().Build()
	return err
}

func findTheOrigin(pipe model.PipeExec, graph model.PipeFragment) (model.NodeBasicTag, error) {
	rootVid := pipe.VID
	// 99.999999% of the time, the root node is the first node in the graph
	for _, vertex := range graph.Vertexes {
		if vertex.VID == rootVid {
			return vertex, nil
		}
	}
	return model.NodeBasicTag{}, fmt.Errorf("root node not found in graph")
}
