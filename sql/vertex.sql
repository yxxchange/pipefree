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