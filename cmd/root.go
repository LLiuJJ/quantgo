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

package cmd

import (
	"fmt"
	"os"

	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	quantgo_storage "github.com/LLiuJJ/quantgo/storage"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "quantgo",
	Short: "quantgo is a Quantitative Trading System",
	Long:  `quantgo is a Quantitative Trading System`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite3", "./stock_history")
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
		defer db.Close()
		quantgo_storage.InitStorage(db)
		fmt.Println("A Quantitative Trading System In Golang.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
