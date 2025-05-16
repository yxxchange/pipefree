package nebula

import "testing"

// 测试用结构体实现 Vertex 接口
type testVertex struct {
	vid      string
	tag      string
	propsMap map[string]interface{}
}

func (t testVertex) VID() string {
	return t.vid
}

func (t testVertex) TagName() string {
	return t.tag
}

func (t testVertex) Props() map[string]interface{} {
	return t.propsMap
}

type testEdge struct {
	edgeType string
	src      string
	dst      string
	props    map[string]interface{}
}

func (t testEdge) EdgeType() string {
	return t.edgeType
}

func (t testEdge) Src() string {
	return t.src
}

func (t testEdge) Dst() string {
	return t.dst
}

func (t testEdge) Props() map[string]interface{} {
	return t.props
}

func TestBuildInsertVertexSQL(t *testing.T) {
	tests := []struct {
		name        string
		vertex      Vertex
		expectedSQL string
	}{
		{
			name: "单属性测试",
			vertex: testVertex{
				tag:      "person",
				propsMap: map[string]interface{}{"name": "Alice"},
			},
			expectedSQL: "INSERT VERTEX person (name) VALUES ($name)",
		},
		{
			name: "多属性测试（排序验证）",
			vertex: testVertex{
				tag: "product",
				propsMap: map[string]interface{}{
					"price": 10.5,
					"stock": 100,
					"name":  "ItemA",
				},
			},
			expectedSQL: "INSERT VERTEX product (name, price, stock) VALUES ($name, $price, $stock)",
		},
		{
			name: "空属性测试",
			vertex: testVertex{
				tag:      "empty",
				propsMap: map[string]interface{}{},
			},
			expectedSQL: "INSERT VERTEX empty () VALUES ()",
		},
		{
			name: "特殊字符属性名测试",
			vertex: testVertex{
				tag: "special",
				propsMap: map[string]interface{}{
					"_id":       1,
					"age2":      30,
					"full-name": "Test User",
				},
			},
			expectedSQL: "INSERT VERTEX special (_id, age2, full-name) VALUES ($_id, $age2, $full-name)",
		},
		{
			name: "空标签名测试",
			vertex: testVertex{
				tag:      "",
				propsMap: map[string]interface{}{"foo": "bar"},
			},
			expectedSQL: "INSERT VERTEX  (foo) VALUES ($foo)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := BuildInsertVertexSQL(tt.vertex)
			if sql != tt.expectedSQL {
				t.Errorf("\nGot:  %s\nWant: %s", sql, tt.expectedSQL)
			}
		})
	}
}

func TestBuildInsertEdgeSQL(t *testing.T) {
	tests := []struct {
		name     string
		edge     Edge
		expected string
	}{
		{
			name: "Simple edge with one prop",
			edge: testEdge{
				edgeType: "like",
				src:      "user1",
				dst:      "user2",
				props: map[string]interface{}{
					"degree": 95,
				},
			},
			expected: `INSERT EDGE like (degree) VALUES user1->user2:($degree)`,
		},
		{
			name: "Edge with multiple props",
			edge: testEdge{
				edgeType: "follow",
				src:      "player100",
				dst:      "player101",
				props: map[string]interface{}{
					"since":  "2025-01-01",
					"score":  90,
					"active": true,
				},
			},
			expected: `INSERT EDGE follow (active, score, since) VALUES player100->player101:($active, $score, $since)`,
		},
		{
			name: "Edge with empty props",
			edge: testEdge{
				edgeType: "connect",
				src:      "a",
				dst:      "b",
				props:    map[string]interface{}{},
			},
			expected: `INSERT EDGE connect () VALUES a->b:()`,
		},
		{
			name: "Edge with numeric keys",
			edge: testEdge{
				edgeType: "visit",
				src:      "user1",
				dst:      "page2",
				props: map[string]interface{}{
					"timestamp": 167890,
					"duration":  "30s",
				},
			},
			expected: `INSERT EDGE visit (duration, timestamp) VALUES user1->page2:($duration, $timestamp)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildInsertEdgeSQL(tt.edge)
			if got != tt.expected {
				t.Errorf("BuildInsertEdgeSQL() = \n%s\nwant\n%s\n", got, tt.expected)
			}
		})
	}
}
