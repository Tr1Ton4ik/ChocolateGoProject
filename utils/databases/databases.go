package databases

import (
	"database/sql"
	"fmt"
)

func CreateTables(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Ошибка при создании транзакции: ", err)
	}
	createChocolateTable(tx)
	createCategoryOfChocolateTable(tx)
	createAdminsTable(tx)

	err = tx.Commit()
	if err!= nil {
        fmt.Println("Ошибка при сохранении изменений в базу CreateTables:", err)
    }
}
func createChocolateTable(tx *sql.Tx) {
	_, err := tx.Exec("CREATE TABLE IF NOT EXISTS chocolate (id INTEGER PRIMARY KEY, name TEXT, price INTEGER, description TEXT, category_id INTEGER, FOREIGN KEY (category_id) REFERENCES categories (id))")
	if err != nil {
		fmt.Println("Ошибка при создании таблицы chocolate: ", err)
	}
}
func createCategoryOfChocolateTable(tx *sql.Tx) {
	_, err := tx.Exec("CREATE TABLE IF NOT EXISTS categories (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		fmt.Println("Ошибка при создании таблицы categories: ", err)
	}
}
func createAdminsTable(tx *sql.Tx) {
	_, err := tx.Exec("CREATE TABLE IF NOT EXISTS admins (id INTEGER PRIMARY KEY, name TEXT, password TEXT)")
	if err != nil {
		fmt.Println("Ошибка при создании таблицы admins: ", err)
	}
}
