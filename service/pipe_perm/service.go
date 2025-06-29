package pipe_perm

import (
	"context"
	"errors"
	"fmt"
	"github.com/yxxchange/pipefree/infra/dal/dao"
	"gorm.io/gen"
	"gorm.io/gorm"
)

const (
	PipePermissionBind = "pipe.permission.bind"

	ErrorCode = 10003
)

type Service struct {
	permissionItem dao.IPermissionItemDo
	pipeSpace      dao.IPipeSpaceDo
	nodeNamespace  dao.INodeNamespaceDo
	query          *dao.Query
	ctx            context.Context
}

func NewService(ctx context.Context) *Service {
	return &Service{
		permissionItem: dao.Q.PermissionItem.WithContext(ctx),
		pipeSpace:      dao.Q.PipeSpace.WithContext(ctx),
		query:          dao.Q,
		ctx:            ctx,
	}
}

func (s *Service) CreatePipeSpace(space string) (err error) {
	_, err = s.pipeSpace.Where(dao.PipeSpace.Space.Eq(space)).FirstOrCreate()
	return
}

func (s *Service) CreateNodeNamespace(namespace string) (err error) {
	_, err = s.nodeNamespace.Where(dao.NodeNamespace.Namespace.Eq(namespace)).FirstOrCreate()
	return
}

func (s *Service) CreatePermissionItem(space, namespace string) (err error) {
	pipeSpace, err := s.pipeSpace.Where(dao.PipeSpace.Space.Eq(space)).First()
	if err != nil {
		return err
	}
	nodeNamespace, err := s.nodeNamespace.Where(dao.NodeNamespace.Namespace.Eq(namespace)).First()
	if err != nil {
		return err
	}
	conditions := []gen.Condition{
		dao.PermissionItem.SpaceId.Eq(pipeSpace.Id),
		dao.PermissionItem.NamespaceId.Eq(nodeNamespace.Id),
		dao.PermissionItem.Space.Eq(space),
		dao.PermissionItem.Namespace.Eq(namespace),
		dao.PermissionItem.PermissionInstance.Eq(PipePermissionBind),
	}
	_, err = s.permissionItem.Where(conditions...).FirstOrCreate()
	return
}

func (s *Service) PermissionCheck(space, namespace, expectedPerm string) error {
	conditions := []gen.Condition{
		dao.PermissionItem.Space.Eq(space),
		dao.PermissionItem.Namespace.Eq(namespace),
		dao.PermissionItem.PermissionInstance.Eq(expectedPerm),
	}
	_, err := s.permissionItem.Where(conditions...).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("no permission")
		}
		return fmt.Errorf("permssion check failed: %v", err)
	}
	return nil
}

func (s *Service) PermissionBatchCheck(space string, namespaces []string, expectedPerm string) error {
	conditions := []gen.Condition{
		dao.PermissionItem.Space.Eq(space),
		dao.PermissionItem.PermissionInstance.Eq(expectedPerm),
	}
	items, err := s.permissionItem.Where(conditions...).Find()
	if err != nil {
		return fmt.Errorf("permssion check failed: %v", err)
	}
	namespaceMap := make(map[string]struct{})
	for _, item := range items {
		namespaceMap[item.Namespace] = struct{}{}
	}
	for _, ns := range namespaces {
		if _, ok := namespaceMap[ns]; !ok {
			return fmt.Errorf("permssion check failed: lack the permission of %s for space %s", ns, space)
		}
	}
	return nil
}
