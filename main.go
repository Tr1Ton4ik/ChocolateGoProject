package main

import (
	"chocolateproject/config"
	"chocolateproject/utils/commands"
	"chocolateproject/utils/databases"
	"chocolateproject/utils/types"

	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

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
	fmt.Println(`
Запуск программы

Справка по пользованию:
s - остановка программы
	`)
	if _, err := os.Stat("./db"); os.IsNotExist(err) {
		fmt.Println("Создаю папку под дб")
		os.Mkdir("./db", 0755)
	}
	if _, err := os.Stat("./db/chocolate.db"); os.IsNotExist(err) {
		fmt.Println("Создаю chocolate.db")
		os.Chdir("./db")
		file, _ := os.Create("chocolate.db")
		file.Close()
		os.Chdir("../")
	}
	go WaitForCommands()

	databases.CreateTables()
	//fillAllDatabases() //раскоментить для заполнения бд

	serverMux := http.NewServeMux()
	serverMux.Handle(frontPath+"img/", http.StripPrefix(frontPath+"img/", http.FileServer(http.Dir(strings.Trim(frontPath+"img/", "/")))))
	serverMux.Handle(frontPath+"js/", http.StripPrefix(frontPath+"js/", http.FileServer(http.Dir(strings.Trim(frontPath+"js/", "/")))))
	serverMux.Handle(frontPath+"html/", http.StripPrefix(frontPath+"html/", http.FileServer(http.Dir(strings.Trim(frontPath+"html/", "/")))))
	serverMux.Handle(frontPath+"css/", http.StripPrefix(frontPath+"css/", http.FileServer(http.Dir(strings.Trim(frontPath+"css/", "/")))))

	serverMux.HandleFunc("/", mainPageHandler)
	serverMux.HandleFunc("/cart.html", cartPageHandler)
	serverMux.HandleFunc("/product.html", productPageHandler)
	serverMux.HandleFunc("/contact.html", contactPageHanfler)
	serverMux.HandleFunc("/about.html", aboutPageHandler)
	serverMux.HandleFunc("/admin", adminPageHandler)
	serverMux.HandleFunc("/login", loginAdminPageHandler)

	// Запускаем веб-сервер на порту 8080 с нашим serverMux (в прошлых примерах был nil)
	fmt.Println("Запуск сервера ")
	fmt.Println("Сервер запущен: http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", serverMux)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
func mainPageHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("sqlite3", "./db/chocolate.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

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
func contactPageHanfler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("front/html/", "contact.html")
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		panic(err)
	}
	//выводим шаблон клиенту в браузер
	tmpl.ExecuteTemplate(w, "contact", nil)
}
func aboutPageHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("front/html/", "about.html")
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		panic(err)
	}
	//выводим шаблон клиенту в браузер
	tmpl.ExecuteTemplate(w, "about", nil)
}
func adminPageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		file, _, err := r.FormFile("image")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Сохраняем файл на диск
		r.ParseForm()
		name, price_str, description, category := r.Form.Get("name"), r.Form.Get("price"), r.Form.Get("description"), r.Form.Get("category")
		price, _ := strconv.Atoi(price_str)
		fileName := "./front/img/" + name + ".jpeg"
		f, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			panic(err)
		}

		err = addChocolate(name, price, description, category)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Продукт", name, "добавлен")
		}
	}
	path := filepath.Join("front/html/", "admin.html")
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		panic(err)
	}
	//выводим шаблон клиенту в браузер
	db, err := sql.Open("sqlite3", "./db/chocolate.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	categories, category := []string{}, ""
	rows, _ := db.Query("select name from categories")
	for rows.Next() {
		rows.Scan(&category)
		categories = append(categories, category)
	}
	tmpl.ExecuteTemplate(w, "admin", categories)
}
func loginAdminPageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		db, err := sql.Open("sqlite3", "./db/chocolate.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		r.ParseForm()
		name, password := r.Form.Get("name"), r.Form.Get("password")

		rows, _ := db.Query("select * from admins where (name = (?) AND password = (?))", name, config.Hash(password))
		i := 0
		for rows.Next() {
			i++
			break
		}
		if i >= 1 {
			http.SetCookie(w, &http.Cookie{
				Name:  "isAuthenticated",
				Value: "true",
			})
			http.Redirect(w, r, "/admin", http.StatusFound)
		}
	}
	path := filepath.Join("front/html/", "login.html")
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		panic(err)
	}
	//выводим шаблон клиенту в браузер
	tmpl.ExecuteTemplate(w, "login", nil)
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
		products.appendCategory(types.Category{Id: int(category_id), Name: category_name})
	}
	return products
}
func prepareForProductPage(productName string) types.Product {
	db, err := sql.Open("sqlite3", "./db/chocolate.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
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

func fillAllDatabases() {
	fillCategoryDataTable()
	fillAdminDataTable()
}

func addChocolate(name string, price int, description string, category string) error {
	db, err := sql.Open("sqlite3", "./db/chocolate.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var category_id int
	rows, _ := db.Query("select id from categories where name==(?)", category)
	for rows.Next() {
		rows.Scan(&category_id)
		break
	}

	stmnt, _ := db.Prepare("INSERT INTO chocolate (name, price, description, category_id) VALUES (?,?,?,?)")
	_, err = stmnt.Exec(name, price, description, category_id)
	return err
}
func addCategory(name string) error {
	db, err := sql.Open("sqlite3", "./db/chocolate.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	stmnt, _ := db.Prepare("INSERT INTO categories (name) VALUES (?)")
	_, err = stmnt.Exec(name)
	return err
}
func fillCategoryDataTable() {
	categories := []string{"Шоколадные плитки", "Батончики", "Конфеты", "Мороженое", "Пасты и карамели"}
	for _, item := range categories {
		if err := addCategory(item); err != nil {
			panic(err)
		}
	}
}
func fillAdminDataTable() {
	db, err := sql.Open("sqlite3", "./db/chocolate.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	stmnt, _ := db.Prepare("INSERT INTO admins (name, password) VALUES (?, ?)")
	_, err = stmnt.Exec(config.FirstAdmin, config.FirstPassword)
	if err != nil {
		panic(err)
	}
}
func WaitForCommands() {
	var command string
	for {
		fmt.Scan(&command)
		switch command {
		case "s":
			commands.Stop()
		default:
			fmt.Println("No such command")
		}
	}
}
