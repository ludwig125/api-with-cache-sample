// go:build integration

package integration

/* VS Codeでの実行はGOFLAGの設定が必要そう
"go.toolsEnvVars": {
	"GO111MODULE": "on",
	"GOBIN": "/home/ludwig125/go/bin",
	"GOFLAGS": "-tags=integration"
},

$ go test -v ./integration -tags=integration

VSCodeのReload Windowなども必要かもしれない

*/

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"os"

	_ "github.com/mattn/go-sqlite3"
)

func TestHoge(t *testing.T) {

	dbName := "../item_db"
	if err := makeDB(dbName); err != nil {
		log.Fatalf("failed to makeDB: %v", err)
	}

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatalf("failed to open db: %s", dbName)
	}

	if err := makeTables(db, dbName); err != nil {
		log.Fatalf("failed to makeTables: %v", err)
	}

	if err := insertTestData(db); err != nil {
		log.Fatalf("failed to insertTestData: %v", err)
	}
	log.Println("insert testData successfully")

}

func makeDB(dbName string) error {
	if err := os.Remove(dbName); err != nil {
		return fmt.Errorf("failed to remove db: %s", dbName)
	}

	if _, err := os.Create(dbName); err != nil {
		return fmt.Errorf("failed to create db: %s", dbName)
	}
	return nil
}

func makeTables(db *sql.DB, dbName string) error {
	// テーブル作成
	if _, err := db.Exec(
		`CREATE TABLE item(id INTEGER PRIMARY KEY ASC, name TEXT, price INTEGER)`,
	); err != nil {
		return fmt.Errorf("failed to create table item")
	}
	if _, err := db.Exec(
		`CREATE TABLE score(id INTEGER PRIMARY KEY ASC, name TEXT, score INTEGER)`,
	); err != nil {
		return fmt.Errorf("failed to create table score")
	}
	return nil
}

func insertTestData(db *sql.DB) error {
	n := 9
	order := 100

	// 一度に大量に入れるとsqliteのバイト数制限に引っかかるので分けて入れる
	for i := 0; i < n; i++ {
		start := i*order + 1
		end := (i + 1) * order

		if err := insertItem(db, start, end); err != nil {
			return fmt.Errorf("failed to insertItem: %v", err)
		}
		if err := insertScore(db, start, end); err != nil {
			return fmt.Errorf("failed to insertScore: %v", err)
		}
	}
	return nil
}

func insertItem(db *sql.DB, start, end int) error {
	var values string
	for i := start; i <= end; i++ {
		name := fmt.Sprintf("Item%d", i)
		price := i
		if i == end {
			values += fmt.Sprintf(`("%s", %d)`, name, price)
			break
		}
		values += fmt.Sprintf(`("%s", %d),`, name, price)
	}

	// fmt.Println("values", values)
	if _, err := db.Exec(
		fmt.Sprintf(`INSERT INTO item (name, price) VALUES %s`, values),
	); err != nil {
		return fmt.Errorf("failed to insert: %v", err)
	}
	return nil
}

func insertScore(db *sql.DB, start, end int) error {
	var values string
	for i := start; i <= end; i++ {
		name := fmt.Sprintf("Item:%d", i)
		score := i % 10 // 10の余りをスコアにする
		if i == end {
			values += fmt.Sprintf(`("%s", %d)`, name, score)
			break
		}
		values += fmt.Sprintf(`("%s", %d),`, name, score)
	}

	if _, err := db.Exec(
		fmt.Sprintf(`INSERT INTO score (name, score) VALUES %s`, values),
	); err != nil {
		return fmt.Errorf("failed to insert: %v", err)
	}
	return nil
}
