package nebula

import (
	"github.com/yxxchange/pipefree/helper/str"
	"sort"
	"strings"
)

// INSERT VERTEX [IF NOT EXISTS] [tag_props, [tag_props] ...]
// VALUES VID: ([prop_value_list])
//
//	tag_props:
//  	tag_name ([prop_name_list])
//
//	prop_name_list:
//   	[prop_name [, prop_name] ...]
//
//	prop_value_list:
//   	[prop_value [, prop_value] ...]

// INSERT EDGE [IF NOT EXISTS] <edge_type> ( <prop_name_list> ) VALUES
// <src_vid> -> <dst_vid>[@<rank>] : ( <prop_value_list> )
// [, <src_vid> -> <dst_vid>[@<rank>] : ( <prop_value_list> ), ...];
//
// <prop_name_list> ::=
//   [ <prop_name> [, <prop_name> ] ...]
//
// <prop_value_list> ::=
//   [ <prop_value> [, <prop_value> ] ...]

const (
	PlaceHolderTag   = "[tag]"
	PlaceHolderProp  = "[prop]"
	PlaceHolderValue = "[value]"
	PlaceHolderEdge  = "[edge]"
	PlaceHolderVid   = "[vid]"

	TagPropsTmpl          = PlaceHolderTag + " (" + PlaceHolderProp + ")"
	EdgePropsTmpl         = PlaceHolderEdge + " (" + PlaceHolderProp + ")"
	PropValueListTmpl     = "(" + PlaceHolderValue + ")"
	EdgePropValueListTmpl = PlaceHolderVid + "->" + PlaceHolderVid + ":" + "(" + PlaceHolderValue + ")"

	InsertVertexTmpl = "INSERT VERTEX " + TagPropsTmpl + " VALUES " + PropValueListTmpl
	InsertEdgeTmpl   = "INSERT EDGE " + EdgePropsTmpl + " VALUES " + EdgePropValueListTmpl
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
	tag := vertex.TagName()
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
	tmpl = str.Replace(tmpl, PlaceHolderTag, tag)
	tmpl = str.Replace(tmpl, PlaceHolderProp, strings.Join(props, ", "))
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
