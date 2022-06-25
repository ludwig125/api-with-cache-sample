package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteRepository struct {
	db *sql.DB
}

// interfaceを実装しているか保証する
// See: http://golang.org/doc/faq#guarantee_satisfies_interface
var _ ItemRepository = (*sqliteRepository)(nil)

func NewSQLiteItemRepository(dbName string) (ItemRepository, error) {
	// 対象のDBがなくても新規に作ってしまうようなので、DBファイルの存在確認する
	if !exists(dbName) {
		return nil, fmt.Errorf("no such db file: %s", dbName)
	}

	db, err := connSQLite(dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to connection db: %v", err)
	}
	log.Printf("connected %s successfully", dbName)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %v", err)
	}
	log.Printf("ping %s successfully", dbName)
	return &sqliteRepository{db: db}, nil
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func connSQLite(dbName string) (*sql.DB, error) {
	// DNS: root:password@tcp(ipaddress:port)/dbname
	// https://github.com/go-sql-driver/mysql#examples
	// パスワードなしで、localhostに対して、デフォルトの3306 portに接続する場合は以下でいい
	return sql.Open("sqlite3", dbName)
}

func (r *sqliteRepository) GetAll() ([]Item, error) {
	rows, err := r.db.Query("SELECT * FROM item")
	if err != nil {
		return nil, fmt.Errorf("failed to select all items, err: %v", err)
	}
	return scanItems(rows)
}

func (r *sqliteRepository) GetItems(ids []ID) ([]Item, error) {
	// PrepareのIN句の中のplaceholder"?"をidsの数だけカンマでつなぐ
	inStmt := ""
	for i := 1; i <= len(ids); i++ {
		if i == len(ids) {
			inStmt += "?"
			break
		}
		inStmt += "?,"
	}

	stmt, err := r.db.Prepare(fmt.Sprintf("SELECT * FROM item WHERE id IN(%s)", inStmt))
	if err != nil {
		return nil, fmt.Errorf("failed to Prepare: %v", err)
	}

	// idsをinterfaceのsliceにしてQueryに渡す
	// ref: https://stackoverflow.com/questions/53983170/sql-converting-argument-1-type-unsupported-type-int-a-slice-of-in
	args := idsToArgs(ids)
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to Query: %v", err)
	}
	return scanItems(rows)
}

func idsToArgs(ids []ID) []interface{} {
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	return args
}

func scanItems(rows *sql.Rows) ([]Item, error) {
	var Items []Item
	defer rows.Close()
	for rows.Next() {
		var a Item
		err := rows.Scan(&a.ID, &a.Name, &a.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %v", err)
		}
		Items = append(Items, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %v", err)
	}
	return Items, nil
}
