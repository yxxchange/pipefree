# CREATE TAG [IF NOT EXISTS] <tag_name>
#     (
#       <prop_name> <data_type> [NULL | NOT NULL] [DEFAULT <default_value>] [COMMENT '<comment>']
#       [{, <prop_name> <data_type> [NULL | NOT NULL] [DEFAULT <default_value>] [COMMENT '<comment>']} ...]
#     )
#     [TTL_DURATION = <ttl_duration>]
#     [TTL_COL = <prop_name>]
#     [COMMENT = '<comment>'];

CREATE TAG IF NOT EXISTS node_cfg (
    node_name STRING NOT NULL COMMENT 'node name',
    node_vid STRING NOT NULL COMMENT 'node vid',
    node_kind STRING NOT NULL COMMENT 'node kind',
    operation string NOT NULL COMMENT 'operation of the node',
    node_desc string NOT NULL COMMENT 'description of the node',
    spec string NOT NULL COMMENT 'spec param of the node',
    `from` string,
    `to` string,
)

# CREATE EDGE [IF NOT EXISTS] <edge_type_name>
#     (
#       <prop_name> <data_type> [NULL | NOT NULL] [DEFAULT <default_value>] [COMMENT '<comment>']
#       [{, <prop_name> <data_type> [NULL | NOT NULL] [DEFAULT <default_value>] [COMMENT '<comment>']} ...]
#     )
#     [TTL_DURATION = <ttl_duration>]
#     [TTL_COL = <prop_name>]
#     [COMMENT = '<comment>'];

CREATE EDGE IF NOT EXISTS node_flow (
    src_vid STRING NOT NULL COMMENT 'source node vid',
    dst_vid STRING NOT NULL COMMENT 'destination node vid',
)

# INSERT VERTEX [IF NOT EXISTS] <tag_name> (<prop_name_list>) [, <tag_name> (<prop_name_list>), ...]
#     {VALUES | VALUE} VID: (<prop_value_list>[, <prop_value_list>])
#
#     prop_name_list:
#     [prop_name [, prop_name] ...]
#
#     prop_value_list:
#     [prop_value [, prop_value] ...]

INSERT VERTEX node_cfg (node_name, node_vid, node_kind, operation, node_desc, spec) VALUES
    "node1": (node_name="node1", node_vid="node1", node_kind="kind1", operation="op1", node_desc="desc1", spec="spec1"),
    "node2": (node_name="node2", node_vid="node2", node_kind="kind2", operation="op2", node_desc="desc2", spec="spec2");