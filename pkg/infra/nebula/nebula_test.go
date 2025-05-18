package nebula

import (
	"github.com/yxxchange/pipefree/config"
	"github.com/yxxchange/pipefree/helper/log"
	"testing"
)

type TestVertex struct {
	Vid  string `nebula:"vid"`
	Name string `nebula:"name"`
}

func (t TestVertex) VID() string {
	return t.Vid
}

func (t TestVertex) TagName() string {
	return "test"
}

func (t TestVertex) Props() map[string]interface{} {
	return map[string]interface{}{
		"name": "fghjk",
	}
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

}

func initTestGraph() (err error) {
	t11 := TestVertex{
		Vid:  "111",
		Name: "t11",
	}
	res := HandleSQL("test", BuildInsertVertexSQL(t11), t11.Props())
	if res.Err != nil {
		panic(res.Err)
	}
	//res := HandleSQL("test", BuildGoNStepsSQL(1, t11.VID(), "edge_test", Yield(t11)), t11.Props())
	//if res.Err != nil {
	//	log.Errorf("err: %v", res.Err)
	//	return res.Err
	//}
	//var tRes []TestVertex
	//err = res.Res.Scan(&tRes)
	//if err != nil {
	//	log.Error(err.Error())
	//	return
	//}
	//log.Infof("%+v", tRes)
	//return
	return

}
