package nebula

import (
	"github.com/yxxchange/pipefree/helper/str"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const (
	PlaceHolderTag   = "[tag]"
	PlaceHolderProp  = "[prop]"
	PlaceHolderValue = "[value]"
	PlaceHolderEdge  = "[edge]"
	PlaceHolderVid   = "[vid]"
	PlaceHolderN     = "[N]"
	PlaceHolderYield = "[Yield]"

	TagPropsTmpl            = PlaceHolderTag + " (" + PlaceHolderProp + ")"
	EdgePropsTmpl           = PlaceHolderEdge + " (" + PlaceHolderProp + ")"
	VertexPropValueListTmpl = PlaceHolderVid + ":" + "(" + PlaceHolderValue + ")"
	EdgePropValueListTmpl   = PlaceHolderVid + "->" + PlaceHolderVid + ":" + "(" + PlaceHolderValue + ")"

	InsertVertexTmpl = "INSERT VERTEX " + TagPropsTmpl + " VALUES " + VertexPropValueListTmpl
	InsertEdgeTmpl   = "INSERT EDGE " + EdgePropsTmpl + " VALUES " + EdgePropValueListTmpl
	GONStepsTmpl     = "GO " + PlaceHolderN + " STEPS" + " FROM " + PlaceHolderVid + " OVER " + PlaceHolderEdge + " YIELD " + PlaceHolderYield
)

type Vertex interface {
	VID() string
	TagName() string
	Props() map[string]interface{}
}

type Edge interface {
	Dst() string
	Src() string
	EdgeType() string
	Props() map[string]interface{}
}

func BuildInsertVertexSQL(vertex Vertex) string {
	tmpl := InsertVertexTmpl
	props := make([]string, 0)
	holders := make([]string, 0)
	keys := make([]string, 0, len(vertex.Props()))
	for k := range vertex.Props() {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		props = append(props, k)
		holders = append(holders, "$"+k)
	}
	tmpl = str.Replace(tmpl, PlaceHolderTag, vertex.TagName())
	tmpl = str.Replace(tmpl, PlaceHolderProp, strings.Join(props, ", "))
	tmpl = str.Replace(tmpl, PlaceHolderVid, vertex.VID())
	tmpl = str.Replace(tmpl, PlaceHolderValue, strings.Join(holders, ", "))
	return tmpl
}

func BuildInsertEdgeSQL(edge Edge) string {
	tmpl := InsertEdgeTmpl
	edgeType := edge.EdgeType()
	props := make([]string, 0)
	holders := make([]string, 0)
	keys := make([]string, 0, len(edge.Props()))
	for k := range edge.Props() {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		props = append(props, k)
		holders = append(holders, "$"+k)
	}
	tmpl = str.Replace(tmpl, PlaceHolderEdge, edgeType)
	tmpl = str.Replace(tmpl, PlaceHolderProp, strings.Join(props, ", "))
	tmpl = str.Replace(tmpl, PlaceHolderVid, edge.Src())
	tmpl = str.Replace(tmpl, PlaceHolderVid, edge.Dst())
	tmpl = str.Replace(tmpl, PlaceHolderValue, strings.Join(holders, ", "))
	return tmpl
}

func BuildGoNStepsSQL(N int, vid, edgeType, yield string) string {
	tmpl := GONStepsTmpl
	tmpl = str.Replace(tmpl, PlaceHolderN, strconv.Itoa(N))
	tmpl = str.Replace(tmpl, PlaceHolderVid, vid)
	tmpl = str.Replace(tmpl, PlaceHolderEdge, edgeType)
	tmpl = str.Replace(tmpl, PlaceHolderYield, yield)
	return tmpl
}

func Yield(vertex Vertex) string {
	yieldStr := ""
	flag := false
	tt := reflect.TypeOf(vertex)
	for i := 0; i < tt.NumField(); i++ {
		field := tt.Field(i)
		tagValue := field.Tag.Get("nebula")
		if tagValue == "" || tagValue == "vid" {
			continue
		}
		if flag {
			yieldStr += ", "
		}
		yieldStr += "$$." + vertex.TagName() + "." + tagValue + " AS " + tagValue
		flag = true
	}
	if flag {
		yieldStr += ", "
	}
	yieldStr += "id($$) AS vid"
	return yieldStr
}
