package resources

var readConnection string = `
SELECT
	mz_connections.id,
	mz_connections.name AS connection_name,
	mz_schemas.name AS schema_name,
	mz_databases.name AS database_name,
	mz_connections.type AS connection_type
FROM mz_connections
JOIN mz_schemas
	ON mz_connections.schema_id = mz_schemas.id
JOIN mz_databases
	ON mz_schemas.database_id = mz_databases.id
WHERE mz_connections.id = 'u1';`

var readConnectionAwsPrivatelink string = `
SELECT
	mz_connections.id,
	mz_connections.name AS connection_name,
	mz_schemas.name AS schema_name,
	mz_databases.name AS database_name,
	mz_aws_privatelink_connections.principal
FROM mz_connections
JOIN mz_schemas
	ON mz_connections.schema_id = mz_schemas.id
JOIN mz_databases
	ON mz_schemas.database_id = mz_databases.id
LEFT JOIN mz_aws_privatelink_connections
	ON mz_connections.id = mz_aws_privatelink_connections.id
WHERE mz_connections.id = 'u1';`

var readConnectionSshTunnellink string = `
SELECT
	mz_connections.id,
	mz_connections.name AS connection_name,
	mz_schemas.name AS schema_name,
	mz_databases.name AS database_name,
	mz_ssh_tunnel_connections.public_key_1,
	mz_ssh_tunnel_connections.public_key_2
FROM mz_connections
JOIN mz_schemas
	ON mz_connections.schema_id = mz_schemas.id
JOIN mz_databases
	ON mz_schemas.database_id = mz_databases.id
LEFT JOIN mz_ssh_tunnel_connections
	ON mz_connections.id = mz_ssh_tunnel_connections.id
WHERE mz_connections.id = 'u1';`
