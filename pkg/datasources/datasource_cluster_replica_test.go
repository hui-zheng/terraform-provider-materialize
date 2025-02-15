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

func TestClusterReplicaDatasource(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{}
	d := schema.TestResourceDataRaw(t, ClusterReplica().Schema, in)
	r.NotNil(d)

	testhelpers.WithMockDb(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		testhelpers.MockClusterReplicaScan(mock, "")

		if err := clusterReplicaRead(context.TODO(), d, db); err != nil {
			t.Fatal(err)
		}

		if d.Get("cluster_replicas") == nil {
			t.Fatal("Data source not set")
		}
	})
}
