package querygen

import (
	"testing"

	"github.com/andreyvit/diff"
)

func Test_QueryGen(t *testing.T) {
	clusters := TableDesc{
		Name:       "clusters",
		PrimaryKey: "id",
		Cols: []ColumnDesc{
			{"id", TInt},
			{"name", TString},
			{"cockroach_version", TString},
			{"machine_type", TString},
			{"created_at", TTimestamp},
			{"updated_at", TTimestamp},
		},
		Timestamps: Timestamps,
	}

	nodes := TableDesc{
		Name:       "nodes",
		PrimaryKey: "id",
		Cols: []ColumnDesc{
			{"id", TInt},
			{"cluster_id", TInt},
			{"name", TString},
			{"locality", TString},
			{"machine_type", TString},
			{"internal_address", TInet},
			{"external_address", TInet},
			{"created_at", TTimestamp},
			{"updated_at", TTimestamp},
		},
		Timestamps: Timestamps,
	}

	nodesForCluster := &Node{
		TableName: nodes.Name,
		Where:     "nodes.cluster_id = clusters.id",
		Columns:   nodes.Cols,
	}

	clusterRegions := TableDesc{
		Name:       "cluster_regions",
		PrimaryKey: "id",
		Cols: []ColumnDesc{
			{"id", TInt},
			{"cluster_id", TInt},
			{"locality", TString},
			{"num_nodes", TInt},
			{"created_at", TTimestamp},
			{"updated_at", TTimestamp},
		},
		Timestamps: Timestamps,
	}

	regionsForCluster := &Node{
		TableName: clusterRegions.Name,
		Where:     "cluster_regions.cluster_id = clusters.id",
		Columns:   clusterRegions.Cols,
	}

	query := clusters.All().OrderBy("created_at DESC").WithChildren(map[string]*Node{
		"regions": regionsForCluster,
		"nodes":   nodesForCluster,
	})

	expected := `SELECT json_agg(json_build_object(
  'id', id,
  'name', name,
  'cockroach_version', cockroach_version,
  'machine_type', machine_type,
  'created_at', experimental_strftime(created_at, '%Y-%m-%dT%H:%M:%S.%fZ'),
  'updated_at', experimental_strftime(updated_at, '%Y-%m-%dT%H:%M:%S.%fZ'),
  'nodes', (
    SELECT json_agg(json_build_object(
      'id', id,
      'cluster_id', cluster_id,
      'name', name,
      'locality', locality,
      'machine_type', machine_type,
      'internal_address', internal_address,
      'external_address', external_address,
      'created_at', experimental_strftime(created_at, '%Y-%m-%dT%H:%M:%S.%fZ'),
      'updated_at', experimental_strftime(updated_at, '%Y-%m-%dT%H:%M:%S.%fZ')
    ))
    FROM nodes
    WHERE nodes.cluster_id = clusters.id
  ),
  'regions', (
    SELECT json_agg(json_build_object(
      'id', id,
      'cluster_id', cluster_id,
      'locality', locality,
      'num_nodes', num_nodes,
      'created_at', experimental_strftime(created_at, '%Y-%m-%dT%H:%M:%S.%fZ'),
      'updated_at', experimental_strftime(updated_at, '%Y-%m-%dT%H:%M:%S.%fZ')
    ))
    FROM cluster_regions
    WHERE cluster_regions.cluster_id = clusters.id
  )
))
FROM clusters
ORDER BY created_at DESC`
	actual := query.ToSQL().String()
	if expected != actual {
		t.Logf("GOT:\n\n%s\n\nWANTED:\n\n%s", actual, expected)
		t.Fatalf(diff.LineDiff(actual, expected))
	}
}
