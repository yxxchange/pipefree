package model

type PipeSpace struct {
	Basic
	Space string `json:"space" yaml:"space" gorm:"column:space"` // 流水线命名空间
}

func (*PipeSpace) TableName() string {
	return "pipe_space"
}

type NodeNamespace struct {
	Basic
	Namespace string `json:"namespace" yaml:"namespace" gorm:"column:namespace"` // 节点命名空间
}

func (*NodeNamespace) TableName() string {
	return "node_namespace"
}

type PermissionItem struct {
	Basic
	SpaceId            int64  `json:"space_id" yaml:"space_id" gorm:"column:space_id"`                                  // 流水线空间ID
	Space              string `json:"space" yaml:"space" gorm:"column:space"`                                           // 流水线命名空间
	NamespaceId        int64  `json:"namespace_id" yaml:"namespace_id" gorm:"column:namespace_id"`                      // 节点命名空间ID
	Namespace          string `json:"namespace" yaml:"namespace" gorm:"column:namespace"`                               // 节点命名空间
	PermissionInstance string `json:"permission_instance" yaml:"permission_instance" gorm:"column:permission_instance"` // 权限实例
}

func (*PermissionItem) TableName() string {
	return "permission_item"
}
