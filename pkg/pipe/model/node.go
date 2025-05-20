package model

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/yxxchange/pipefree/helper/serialize"
)

type Node struct {
	ApiVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       Kind     `json:"kind" yaml:"kind"`
	MetaData   MetaData `json:"metadata" yaml:"metadata"`
	Spec       Spec     `json:"spec" yaml:"spec"`
	Env        Env      `json:"env" yaml:"env"`
	Status     Status   `json:"status" yaml:"status"`
}

func (n Node) ToString() string {
	b, err := serialize.JsonSerialize(n)
	if err != nil {
		return ""
	}
	return string(b)
}

func (n Node) Validate() error {
	if n.ApiVersion == "" {
		return fmt.Errorf("apiVersion is required")
	}
	if n.Kind == "" {
		return fmt.Errorf("kind is required")
	}
	if n.MetaData.Name == "" {
		return fmt.Errorf("name is required")
	}
	if n.MetaData.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if n.MetaData.Operation == "" {
		return fmt.Errorf("operation is required")
	}
	if n.MetaData.UUID == "" {
		return fmt.Errorf("uuid is required")
	}
	return nil
}

func (n *Node) AddToNode(meta *MetaData) {
	if n.MetaData.To == nil {
		n.MetaData.To = make(map[string]*MetaData)
	}
	n.MetaData.To[meta.Name] = meta
	return
}

func (n *Node) AddFromNode(meta *MetaData) {
	if n.MetaData.From == nil {
		n.MetaData.From = make(map[string]*MetaData)
	}
	n.MetaData.From[meta.Name] = meta
	return
}

func (n Node) ToSchema() Schema {
	return Schema{
		ApiVersion: n.ApiVersion,
		Namespace:  n.MetaData.Namespace,
		Kind:       n.Kind,
		Operation:  n.MetaData.Operation,
	}
}

func (n Node) ToSnapshot() Node {
	if n.Status.Phase == "" {
		n.MetaData.RuntimeUUID = uuid.New().String()
		n.Status.Phase = PhaseReady
	}
	return n
}

type PipeFlow struct {
	Nodes  []Node `json:"nodes,omitempty" yaml:"nodes,omitempty"` // the nodes of the pipe
	Graph  Graph  `json:"graph,omitempty" yaml:"graph,omitempty"` // the graph of the node
	Global Env    `json:"global,omitempty" yaml:"global,omitempty"`
}

func (n PipeFlow) ValidateStaticCfg() error {
	names := make(map[string]bool)
	namespace := make(map[string]bool)
	for _, node := range n.Nodes {
		if names[node.MetaData.Name] {
			return fmt.Errorf("duplicate node name: %s", node.MetaData.Name)
		}
		names[node.MetaData.Name] = true
		if err := node.Validate(); err != nil {
			return err
		}
		namespace[node.MetaData.Namespace] = true
	}
	if len(namespace) > 1 {
		return fmt.Errorf("namespace is not same")
	}
	for _, vertex := range n.Graph.Vertexes {
		if vertex.UUID == "" {
			return fmt.Errorf("vertex uuid is required")
		}
	}

	return nil
}

func (n PipeFlow) ValidateDynamicCfg() error {
	for _, edge := range n.Graph.Edges {
		if edge.SrcUUID == "" || edge.DstUUID == "" {
			return fmt.Errorf("edge src or dst uuid is empty")
		}
	}
	for _, vertex := range n.Graph.Vertexes {
		if vertex.UUID == "" {
			return fmt.Errorf("vertex uuid is empty")
		}
		if vertex.RuntimeUUID == "" {
			return fmt.Errorf("vertex runtime uuid is empty")
		}
	}
	return nil
}

func (n PipeFlow) ToExec() PipeExec {
	var res PipeExec
	m := make(map[string]string) // name: runtime_uuid
	for _, node := range n.Nodes {
		snapshot := node.ToSnapshot()
		res.Nodes = append(res.Nodes, snapshot)
		m[snapshot.MetaData.Name] = snapshot.MetaData.RuntimeUUID
	}
	for _, node := range res.Nodes {
		vertex := Vertex{
			Name:        node.MetaData.Name,
			UUID:        node.MetaData.UUID,
			RuntimeUUID: node.MetaData.RuntimeUUID,
		}
		res.Graph.Vertexes = append(res.Graph.Vertexes, vertex)
	}
	for _, edge := range n.Graph.Edges {
		edge.SrcUUID = m[edge.From]
		edge.DstUUID = m[edge.To]
		res.Graph.Edges = append(res.Graph.Edges, edge)
	}

	return res
}

func (n PipeFlow) ToPipeCfg() PipeConfig {
	var res PipeConfig
	m := make(map[string]string)
	for _, node := range n.Nodes {
		_copy := node
		_copy.MetaData.UUID = uuid.New().String()
		_copy.MetaData.RuntimeUUID = ""
		res.Nodes = append(res.Nodes, _copy)
		m[node.MetaData.Name] = _copy.MetaData.UUID
	}
	for _, vertex := range n.Graph.Vertexes {
		vertex.UUID = m[vertex.Name]
		res.Graph.Vertexes = append(res.Graph.Vertexes, vertex)
	}
	res.Graph.Edges = n.Graph.Edges
	res.Graph.Reference = n.Graph.Reference
	return res
}

// Reference is used to reference the node
// must point to the node of compound kind
type Reference struct {
	Vertex
}

type MetaData struct {
	// Static config
	UUID      string `json:"uuid" yaml:"uuid"` // static uuid
	Name      string `json:"name" yaml:"name"`
	Operation string `json:"operation" yaml:"operation"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Desc      string `json:"desc" yaml:"desc"`

	// Dynamic config
	RuntimeUUID     string `json:"runtime_uuid" yaml:"runtime_uuid"`
	ResourceVersion uint64 `json:"resource_version" yaml:"resource_version"`

	// Graph config
	// the graph config is used to build the graph of the node
	Ancestor *MetaData            `json:"-" yaml:"-"` // the ancestor of the node
	From     map[string]*MetaData `json:"-" yaml:"-"`
	To       map[string]*MetaData `json:"-" yaml:"-"`
}

// Schema means that node can be processed by which engine
// TODO: engine need to subscribe the pipe
type Schema struct {
	ApiVersion string `json:"apiVersion"`
	Namespace  string `json:"namespace"`
	Kind       Kind   `json:"kind"`
	Operation  string `json:"operation"`
}

func (n Schema) Identifier() string {
	// apiVersion/space/tag/kind/operation
	return fmt.Sprintf("/%s/%s/%s/%s", n.ApiVersion, n.Namespace, n.Kind, n.Operation)
}

func (m *MetaData) GetPredecessor() map[string]*MetaData {
	return m.From
}

func (m *MetaData) GetSuccessor() map[string]*MetaData {
	return m.To
}

func (m *MetaData) GetAncestor() *MetaData {
	return m.Ancestor
}

func (m *MetaData) AddFrom(from *MetaData) {
	if m.From == nil {
		m.From = make(map[string]*MetaData)
	}
	m.From[from.Name] = from
}

func (m *MetaData) AddTo(to *MetaData) {
	if m.To == nil {
		m.To = make(map[string]*MetaData)
	}
	m.To[to.Name] = to
}

func (m *MetaData) AddAncestor(ancestor *MetaData) {
	m.Ancestor = ancestor
}

type Spec struct {
	json.RawMessage
}

type Status struct {
	Phase     Phase                    `json:"phase" yaml:"phase"`
	Chains    []Record                 `json:"chains" yaml:"chains"`
	Customize []map[string]interface{} `json:"customize" yaml:"customize"`
}

type Record struct {
	Err  string `json:"err" yaml:"err"`
	Info string `json:"info" yaml:"info"`
	Warn string `json:"warn" yaml:"warn"`
}

type Graph struct {
	Edges     []Edge    `json:"edges,omitempty" yaml:"edges,omitempty"`
	Vertexes  []Vertex  `json:"vertexes,omitempty" yaml:"vertexes,omitempty"`
	Reference Reference `json:"reference,omitempty" yaml:"reference,omitempty"`
}

type Env struct {
	KeyValues map[string]interface{}
}
