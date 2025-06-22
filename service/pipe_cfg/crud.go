package pipe_cfg

import (
	"context"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/infra/dal/dao"
	"github.com/yxxchange/pipefree/infra/dal/model"
)

const ErrorCode = 10001

type Service struct {
	db dao.IPipeCfgDo
}

func NewService(ctx context.Context) *Service {
	return &Service{
		db: dao.Q.PipeCfg.WithContext(ctx),
	}
}

func (s *Service) GetById(pipeId int64) (*model.PipeCfg, error) {
	pipeCfg, err := s.db.Where(dao.PipeCfg.Id.Eq(pipeId)).First()
	if err != nil {
		return nil, err
	}
	return pipeCfg, nil
}

func (s *Service) Create(pipeCfg *model.PipeCfg) error {
	if err := s.Create(pipeCfg); err != nil {
		return err
	}
	return nil
}

func (s *Service) Update(pipeCfg *model.PipeCfg) error {
	if err := s.db.Save(pipeCfg); err != nil {
		return err
	}
	return nil
}

func (s *Service) Delete(pipeId int64) error {
	_, err := s.db.Where(dao.PipeCfg.Id.Eq(pipeId)).Delete()
	if err != nil {
		return err
	}
	log.Infof("Deleted pipe configuration with ID: %d", pipeId)
	return nil
}
