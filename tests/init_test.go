//go:build integration
// +build integration

package tests

import "pvz_controller/tests/postgresql"

var (
	db *postgresql.TDB
)

func init() {
	// тут мы запрашиваем тестовые креды для бд из енв
	// cfg,err := config.FromEnv
	db = postgresql.NewFromEnv()
}
