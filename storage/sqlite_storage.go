// Copyright 2025 The quantgo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"database/sql"
	"log"
)

func InitStorage(db *sql.DB) {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS stock_daily_his (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		symbol VARCHAR NOT NULL,
		date DATE,
		open NUMERIC,
		high NUMERIC,
		low  NUMERIC,
		close NUMERIC,
		volume NUMERIC
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}
