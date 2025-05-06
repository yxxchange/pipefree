package nebula

import (
	"github.com/yxxchange/pipefree/config"
	"github.com/yxxchange/pipefree/helper/log"
	"testing"
)

type TestVertex struct {
	VID  string `norm:"vertex_id"`
	Name string `norm:"prop:name"`
}

func (t TestVertex) VertexID() string {
	return t.VID
}

func (t TestVertex) VertexTagName() string {
	return "test"
}

type TestEdge struct {
	SrcID string `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
	Rank  int    `norm:"edge_rank"`
	Name  string `norm:"prop:name"`
}

func (s TestEdge) EdgeTypeName() string {
	return "edge_test"
}

func TestNebula(t *testing.T) {
	config.Init("../../../config.yaml")
	Init()

	log.Info("nebula test ok")
	err := initTestGraph()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	Use("test").Go(1).From("1").Yield("edge_test")
}

func initTestGraph() (err error) {
	t1 := TestVertex{
		VID:  "1",
		Name: "t1",
	}
	t2 := TestVertex{
		VID:  "2",
		Name: "t2",
	}
	t3 := TestVertex{
		VID:  "3",
		Name: "t3",
	}
	t4 := TestVertex{
		VID:  "4",
		Name: "t4",
	}
	t5 := TestVertex{
		VID:  "5",
		Name: "t5",
	}
	e1 := TestEdge{
		SrcID: "1",
		DstID: "2",
	}
	e2 := TestEdge{
		SrcID: "1",
		DstID: "3",
	}
	e3 := TestEdge{
		SrcID: "2",
		DstID: "4",
	}
	e4 := TestEdge{
		SrcID: "2",
		DstID: "5",
	}
	e5 := TestEdge{
		SrcID: "3",
		DstID: "4",
	}
	err = Use("test").InsertVertex(t1).Exec()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	err = Use("test").InsertVertex(t2).Exec()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	err = Use("test").InsertVertex(t3).Exec()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	err = Use("test").InsertVertex(t4).Exec()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	err = Use("test").InsertVertex(t5).Exec()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	err = Use("test").InsertEdge(e1).Exec()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	err = Use("test").InsertEdge(e2).Exec()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	err = Use("test").InsertEdge(e3).Exec()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	err = Use("test").InsertEdge(e4).Exec()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	err = Use("test").InsertEdge(e5).Exec()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	log.Info("init test graph ok")
	return
}
