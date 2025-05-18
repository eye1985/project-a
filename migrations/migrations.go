package migrations

import _ "embed"

//go:embed 000001_create_users_table.up.sql
var InitMigrationCreate string
