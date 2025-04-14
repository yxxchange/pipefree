package data

type Node struct {
	ApiVersion string `json:"apiVersion" yaml:"apiVersion"` // the api version of the node
	Kind       string `json:"kind" yaml:"kind"`             // the kind of the node
	MetaData   `json:"metadata" yaml:"metadata"`
	Spec       `json:"spec" yaml:"spec"`
	Status     `json:"status" yaml:"status"`

	// attributes only for compound node
	Graph `json:"graph,omitempty" yaml:"graph,omitempty"` // the graph of the node
}

type MetaData struct {
	// Static config
	Name      string   `json:"name" yaml:"name"`
	Namespace string   `json:"namespace" yaml:"namespace"`
	Type      NodeType `json:"type" yaml:"type"`
	Desc      string   `json:"desc" yaml:"desc"`

	// Dynamic config
	RuntimeUUID     string `json:"runtime_uuid" yaml:"runtime_uuid"` // the instance id of the pipe exec snapshot
	ResourceVersion uint64 `json:"resource_version" yaml:"resource_version"`

	// Graph config
	// the graph config is used to build the graph of the node
	ancestor    *MetaData
	predecessor map[string]*MetaData
	successor   map[string]*MetaData
}

func (m *MetaData) GetPredecessor() map[string]*MetaData {
	return m.predecessor
}

func (m *MetaData) GetSuccessor() map[string]*MetaData {
	return m.successor
}

func (m *MetaData) GetAncestor() *MetaData {
	return m.ancestor
}

func (m *MetaData) SetPredecessor(predecessor map[string]*MetaData) {
	m.predecessor = predecessor
}

func (m *MetaData) SetSuccessor(successor map[string]*MetaData) {
	m.successor = successor
}

func (m *MetaData) SetAncestor(ancestor *MetaData) {
	m.ancestor = ancestor
}

func (m *MetaData) AddPredecessor(predecessor *MetaData) {
	if m.predecessor == nil {
		m.predecessor = make(map[string]*MetaData)
	}
	m.predecessor[predecessor.Name] = predecessor
}

func (m *MetaData) AddSuccessor(successor *MetaData) {
	if m.successor == nil {
		m.successor = make(map[string]*MetaData)
	}
	m.successor[successor.Name] = successor
}

func (m *MetaData) AddAncestor(ancestor *MetaData) {
	m.ancestor = ancestor
}

type Spec map[string]interface{}

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
	Nodes     []Node   `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	Reference MetaData `json:"reference,omitempty" yaml:"reference,omitempty"`
}

type Edge struct {
	From string `json:"from" yaml:"from"`
	To   string `json:"to" yaml:"to"`
}
