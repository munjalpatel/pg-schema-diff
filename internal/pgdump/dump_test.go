package pgdump_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/munjalpatel/pg-schema-diff/internal/pgdump"
	"github.com/munjalpatel/pg-schema-diff/internal/pgengine"
	"github.com/stretchr/testify/require"
)

func TestGetDump(t *testing.T) {
	pgEngine, err := pgengine.StartEngine()
	require.NoError(t, err)
	defer pgEngine.Close()

	db, err := pgEngine.CreateDatabase()
	require.NoError(t, err)
	defer db.DropDB()

	connPool, err := sql.Open("pgx", db.GetDSN())
	require.NoError(t, err)
	defer connPool.Close()

	_, err = connPool.ExecContext(context.Background(), `
			CREATE TABLE foobar(foobar_id text);

			INSERT INTO foobar VALUES ('some-id');

			CREATE SCHEMA test;
			CREATE TABLE test.bar(bar_id text);
		`)
	require.NoError(t, err)

	dump, err := pgdump.GetDump(db)
	require.NoError(t, err)
	require.Contains(t, dump, "public.foobar")
	require.Contains(t, dump, "test.bar")
	require.Contains(t, dump, "some-id")

	onlySchemasDump, err := pgdump.GetDump(db, pgdump.WithSchemaOnly())
	require.NoError(t, err)
	require.Contains(t, onlySchemasDump, "public.foobar")
	require.Contains(t, onlySchemasDump, "test.bar")
	require.NotContains(t, onlySchemasDump, "some-id")

	onlyPublicSchemaDump, err := pgdump.GetDump(db, pgdump.WithSchemaOnly(), pgdump.WithExcludeSchema("test"))
	require.NoError(t, err)
	require.Contains(t, onlyPublicSchemaDump, "public.foobar")
	require.NotContains(t, onlyPublicSchemaDump, "test.bar")
	require.NotContains(t, onlyPublicSchemaDump, "some-id")
}
