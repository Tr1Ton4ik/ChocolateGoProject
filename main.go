package main

import (
	"chocolateproject/utils/databases"
	"chocolateproject/utils/types"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var db *sql.DB
var err error

const frontPath string = "/front/"

type AllForMain struct {
	Products   []types.Product
	Categories []types.Category
}

func (a *AllForMain) appendProduct(p types.Product) {
	a.Products = append(a.Products, p)
}
func (a *AllForMain) appendCategory(c types.Category) {
	a.Categories = append(a.Categories, c)
}
func main() {
	go waitForStop()
	db, err = sql.Open("sqlite3", "./db/chocolate.db")
	if err != nil {
		panic(err)
	}
	databases.CreateTables(db)
	//fillAllDatebases() //раскоментить для заполнения бд

	serverMux := http.NewServeMux()
	serverMux.Handle(frontPath+"img/", http.StripPrefix(frontPath+"img/", http.FileServer(http.Dir(strings.Trim(frontPath+"img/", "/")))))
	serverMux.Handle(frontPath+"js/", http.StripPrefix(frontPath+"js/", http.FileServer(http.Dir(strings.Trim(frontPath+"js/", "/")))))
	serverMux.Handle(frontPath+"html/", http.StripPrefix(frontPath+"html/", http.FileServer(http.Dir(strings.Trim(frontPath+"html/", "/")))))
	serverMux.Handle(frontPath+"css/", http.StripPrefix(frontPath+"css/", http.FileServer(http.Dir(strings.Trim(frontPath+"css/", "/")))))

	serverMux.HandleFunc("/", mainPageHandler)
	serverMux.HandleFunc("/cart.html", cartPageHandler)
	serverMux.HandleFunc("/product.html", productPageHandler)

	// Запускаем веб-сервер на порту 8080 с нашим serverMux (в прошлых примерах был nil)
	fmt.Println("Запуск сервера ")
	fmt.Println("Сервер запущен: http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", serverMux)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("front/html/", "index.html")
	//создаем html-шаблон
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		panic(err)
	}
	category := r.URL.Query().Get("category")
	if category == "" {
		category = "Все"
	}
	//выводим шаблон клиенту в браузер
	tmpl.ExecuteTemplate(w, "index", prepareForMainPage(category, db))
}
func cartPageHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("front/html/", "cart.html")
	//создаем html-шаблон
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		panic(err)
	}
	//выводим шаблон клиенту в браузер
	tmpl.ExecuteTemplate(w, "cart", nil)
}
func productPageHandler(w http.ResponseWriter, r *http.Request) {
	var product_name string
	switch r.Method {
	case http.MethodGet:
		product_name = r.URL.Query().Get("name")
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
	path := filepath.Join("front/html/", "product.html")
	//создаем html-шаблон
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		panic(err)
	}
	//выводим шаблон клиенту в браузер
	tmpl.ExecuteTemplate(w, "product", prepareForProductPage(product_name))
}
func prepareForMainPage(categoryForFind string, db *sql.DB) AllForMain {
	var (
		id, price, category_id           int64
		name, description, category_name string
		rows                             *sql.Rows
	)
	if categoryForFind == "Все" {
		rows, _ = db.Query("select * from chocolate")
	} else {
		rows, _ = db.Query("select * from categories where name==(?)", categoryForFind)
		rows.Next()
		rows.Scan(&category_id, &category_name)
		categoryForFind = strconv.Itoa(int(category_id))
		rows, _ = db.Query("select * from chocolate where category_id==(?)", categoryForFind)
	}
	products := AllForMain{}
	for rows.Next() {
		if err := rows.Scan(&id, &name, &price, &description, &category_id); err != nil {
			panic(err)
		}
		category, _ := db.Query("select name from categories where id == (?)", category_id)
		category.Next()
		category.Scan(&category_name)
		products.appendProduct(types.Product{Id: int(id), Name: name, Price: int(price), Description: description, Category: category_name})
	}
	rows, _ = db.Query("select * from categories")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&category_id, &category_name)
		products.appendCategory(types.Category{Id: int(category_id),Name:  category_name})
	}
	return products
}
func prepareForProductPage(productName string) types.Product {
	var (
		id, price, category_id           int
		name, description, category_name string
	)
	rows, _ := db.Query("select * from chocolate where name==(?)", productName)
	rows.Next()
	rows.Scan(&id, &name, &price, &description, &category_id)
	rows, _ = db.Query("select name from categories where id==(?)", category_id)
	rows.Next()
	rows.Scan(&category_name)
	answer := types.Product{Id: int(id), Name: name, Price: int(price), Description: description, Category: category_name}
	return answer
}

func fillAllDatebases() {
	fillProductsDateTable()
	fillCategoryDateTable()
}

func addChocolate(name string, price int, description string, category_id int) error {
	stmnt, _ := db.Prepare("INSERT INTO chocolate (name, price, description, category_id) VALUES (?,?,?,?)")
	_, err := stmnt.Exec(name, price, description, category_id)
	return err
}
func addCategory(name string) error {
	stmnt, _ := db.Prepare("INSERT INTO categories (name) VALUES (?)")
	_, err := stmnt.Exec(name)
	return err
}
func fillProductsDateTable() {
	for i := 1; i <= 5; i++ {
		strFormatI := strconv.Itoa(i)
		err := addChocolate("name"+strFormatI, i*100, "description"+strFormatI, i)
		if err != nil {
			panic(err)
		}
	}
}
func fillCategoryDateTable() {
	categories := []string{"Шоколадные плитки", "Батончики", "Конфеты", "Мороженое", "Пасты и карамели"}
	for _, item := range categories {
		if err := addCategory(item); err != nil {
			panic(err)
		}
	}
}
func waitForStop() {
	var stop string
	for {
		fmt.Scan(&stop)
		if strings.ToLower(stop) == "s" {
			fmt.Println("Successfully stopped")
			os.Exit(0)
		}
	}
}
