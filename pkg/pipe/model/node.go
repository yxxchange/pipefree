package model

import (
	"encoding/json"
	"fmt"
)

type NodeInfo struct {
	ApiVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       Kind     `json:"kind" yaml:"kind"`
	MetaData   MetaData `json:"metadata" yaml:"metadata"`
	Spec       Spec     `json:"spec" yaml:"spec"`
	Env        Env      `json:"env" yaml:"env"`
	Status     Status   `json:"status" yaml:"status"`
}

type Node struct {
	NodeInfo
	Graph *Graph `json:"graph,omitempty" yaml:"graph,omitempty"` // the graph of the node
}

type MetaData struct {
	// Static config
	Name      string `json:"name" yaml:"name"`
	VID       string `json:"vid" yaml:"vid"`     // the unique id of the node
	Space     string `json:"space" yaml:"space"` // like schema
	Tag       string `json:"tag" yaml:"tag"`     // like table
	Operation string `json:"operation" yaml:"operation"`
	Desc      string `json:"desc" yaml:"desc"`

	// Dynamic config
	RuntimeUUID     string `json:"runtime_uuid" yaml:"runtime_uuid"` // the instance id of the pipe exec snapshot
	ResourceVersion uint64 `json:"resource_version" yaml:"resource_version"`

	// Graph config
	// the graph config is used to build the graph of the node
	Ancestor *MetaData            `json:"ancestor" yaml:"ancestor"` // the ancestor of the node
	From     map[string]*MetaData `json:"predecessor" yaml:"predecessor"`
	To       map[string]*MetaData `json:"successor" yaml:"successor"`
}

// NodeIdentifier means that node can be processed by which engine
// TODO: engine need to subscribe the pipe
type NodeIdentifier struct {
	ApiVersion string `json:"apiVersion"`
	Space      string `json:"space"`
	Tag        string `json:"tag"`
	Kind       Kind   `json:"kind"`
	Operation  string `json:"operation"`
}

func (n NodeIdentifier) Identifier() string {
	// apiVersion/space/tag/kind/operation
	return fmt.Sprintf("/%s/%s/%s/%s/%s", n.ApiVersion, n.Space, n.Tag, n.Kind, n.Operation)
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
	Edges     []Edge   `json:"edges,omitempty" yaml:"edges,omitempty"`
	Vertexes  []Node   `json:"vertexes,omitempty" yaml:"vertexes,omitempty"`
	Reference MetaData `json:"reference,omitempty" yaml:"reference,omitempty"`
}

type Edge struct {
	From string `json:"from" yaml:"from"`
	To   string `json:"to" yaml:"to"`
}

type Env struct {
	KeyValues map[string]interface{}
}
