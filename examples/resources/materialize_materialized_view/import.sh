# Materialized views can be imported using the materialized view id:
terraform import materialize_materialized_view.example_materialize_view <view_id>

# Materialized view id and information be found in the `mz_catalog.mz_materialized_views` table
