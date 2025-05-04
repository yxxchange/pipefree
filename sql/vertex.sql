# CREATE TAG [IF NOT EXISTS] <tag_name>
#     (
#       <prop_name> <data_type> [NULL | NOT NULL] [DEFAULT <default_value>] [COMMENT '<comment>']
#       [{, <prop_name> <data_type> [NULL | NOT NULL] [DEFAULT <default_value>] [COMMENT '<comment>']} ...]
#     )
#     [TTL_DURATION = <ttl_duration>]
#     [TTL_COL = <prop_name>]
#     [COMMENT = '<comment>'];

CREATE TAG IF NOT EXISTS node (
    name STRING NOT NULL COMMENT 'node name',



)