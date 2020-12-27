package dsnbuilder

import "testing"

func TestDsnBuilder_Build(t *testing.T) {
	expectedPostgresDSN := "host=localhost port=5432 user=bogaduit_be dbname=bogaduit_core password=debbiegibson1985"

	dsn, err := New("localhost", 5432, "bogaduit_be",
		"debbiegibson1985", "bogaduit_core").Build("postgres")
	if err != nil {
		t.Error("unexpected " + err.Error())
	}

	if dsn != expectedPostgresDSN {
		t.Error("mismatch dsn, expected " + expectedPostgresDSN + " got " + dsn)
	}
}
