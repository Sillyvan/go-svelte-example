//go:build js && wasm

package store

import (
	"database/sql"

	_ "github.com/syumai/workers/cloudflare/d1"
)

const d1BindingName = "DB"

func openDB() (*sql.DB, error) {
	return sql.Open("d1", d1BindingName)
}

func SourceDescription() string {
	return `Cloudflare D1 binding "DB"`
}
