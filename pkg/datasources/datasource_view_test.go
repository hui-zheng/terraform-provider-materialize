package datasources

import (
	"context"
	"testing"

	"github.com/MaterializeInc/terraform-provider-materialize/pkg/testhelpers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestViewDatasource(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"schema_name":   "schema",
		"database_name": "database",
	}
	d := schema.TestResourceDataRaw(t, View().Schema, in)
	r.NotNil(d)

	testhelpers.WithMockDb(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		p := `WHERE mz_databases.name = 'database' AND mz_schemas.name = 'schema'`
		testhelpers.MockViewScan(mock, p)

		if err := viewRead(context.TODO(), d, db); err != nil {
			t.Fatal(err)
		}
	})

}
