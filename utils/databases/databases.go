package databases

import ("database/sql")

func CreateTables() {
	db, err := sql.Open("sqlite3", "./db/chocolate.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	createChocolateTable(db)
	createCategoryOfChocolateTable(db)
	createAdminsTable(db)
}
func createChocolateTable(db *sql.DB) {
	stmnt, err := db.Prepare("CREATE TABLE IF NOT EXISTS chocolate (id INTEGER PRIMARY KEY, name TEXT, price INTEGER, description TEXT, category_id INTEGER, FOREIGN KEY (category_id) REFERENCES categories (id))")
	if err != nil {
		panic(err)
	}
	defer stmnt.Close()
	_, err = stmnt.Exec()
	if err != nil {
		panic(err)
	}
}
func createCategoryOfChocolateTable(db *sql.DB) {
	stmnt, err := db.Prepare("CREATE TABLE IF NOT EXISTS categories (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		panic(err)
	}
	defer stmnt.Close()
	_, err = stmnt.Exec()
	if err != nil {
		panic(err)
	}
}
func createAdminsTable(db *sql.DB){
	stmnt, err := db.Prepare("CREATE TABLE IF NOT EXISTS admins (id INTEGER PRIMARY KEY, name TEXT, password TEXT)")
	if err != nil {
		panic(err)
	}
	defer stmnt.Close()
	_, err = stmnt.Exec()
	if err != nil {
		panic(err)
	}
}