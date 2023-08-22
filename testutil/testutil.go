package testutil

import (
	"database/sql"
	"fmt"
	"strings"
)

// PrepareTestData 外部キー制約のチェックを無効化した状態で第二引数の処理を実行します
func PrepareTestData(db *sql.DB, closure func(db *sql.DB)) {
	MustExec(db, "SET FOREIGN_KEY_CHECKS = 0")
	closure(db)
	MustExec(db, "SET FOREIGN_KEY_CHECKS = 1")
}

// MustInsert データを挿入し、エラーが発生した場合はpanicを発生させます
func MustInsert(db *sql.DB, table string, records []map[string]interface{}) {
	if len(records) == 0 {
		return
	}

	// カラム名とプレースホルダーを生成
	columns := make([]string, 0, len(records[0]))
	placeholders := make([]string, 0, len(records[0]))
	for column := range records[0] {
		columns = append(columns, column)
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	// 各レコードを挿入
	for _, record := range records {
		values := make([]interface{}, 0, len(record))
		for _, column := range columns {
			values = append(values, record[column])
		}

		_, err := db.Exec(query, values...)
		if err != nil {
			panic(err)
		}
	}
}

// MustExec SQLを実行し、エラーが発生した場合はpanicを発生させます
func MustExec(db *sql.DB, sql string) sql.Result {
	result, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
	return result
}
