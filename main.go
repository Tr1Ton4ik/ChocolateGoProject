package main

import (
	"chocolateproject/config"
	"chocolateproject/utils/commands"
	"chocolateproject/utils/databases"
	"chocolateproject/utils/types"

	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const frontPath string = "/front/"

var db, err = sql.Open("sqlite3", "./db/chocolate.db")

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
filldb - заполнение бд
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

	databases.CreateTables(db)

	serverMux := http.NewServeMux()
	serverMux.Handle(frontPath+"img/", http.StripPrefix(frontPath+"img/", http.FileServer(http.Dir(strings.Trim(frontPath+"img/", "/")))))
	serverMux.Handle(frontPath+"js/", http.StripPrefix(frontPath+"js/", http.FileServer(http.Dir(strings.Trim(frontPath+"js/", "/")))))
	serverMux.Handle(frontPath+"html/", http.StripPrefix(frontPath+"html/", http.FileServer(http.Dir(strings.Trim(frontPath+"html/", "/")))))
	serverMux.Handle(frontPath+"css/", http.StripPrefix(frontPath+"css/", http.FileServer(http.Dir(strings.Trim(frontPath+"css/", "/")))))

	serverMux.HandleFunc("/", mainPageHandler)
	serverMux.HandleFunc("/cart", cartPageHandler)
	serverMux.HandleFunc("/product", productPageHandler)
	serverMux.HandleFunc("/contact", contactPageHanfler)
	serverMux.HandleFunc("/about", aboutPageHandler)
	serverMux.HandleFunc("/admin", adminPageHandler)
	serverMux.HandleFunc("/login", loginAdminPageHandler)
	serverMux.HandleFunc("/add_product", addProductHandler)
	serverMux.HandleFunc("/add_admin", addAdminHandler)
	serverMux.HandleFunc("/delete_product", deleteProductHandler)
	serverMux.HandleFunc("/search_product", searchProductHandler)

	// Запускаем веб-сервер на порту 8080 с нашим serverMux (в прошлых примерах был nil)
	fmt.Println("Запуск сервера ")
	APP_IP := os.Getenv("APP_IP")
    APP_PORT := os.Getenv("APP_PORT")

	fmt.Println("Сервер запущен: ", APP_IP+":"+APP_PORT)
	
	err := http.ListenAndServe(APP_IP+":"+APP_PORT, serverMux)
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
	path := filepath.Join("front/html/", "admin.html")
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		panic(err)
	}
	//выводим шаблон клиенту в браузер
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
func addProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
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

		for i := 0; i < 5; i++ { // Пытаемся 5 раз выполнить операцию
			tx, err := db.Begin()
			if err != nil {
				http.Error(w, "Ошибка начала транзакции", http.StatusInternalServerError)
				return
			}
			err = addChocolate(name, price, description, category, tx)
			if err != nil {
				if strings.Contains(err.Error(), "database is locked") {
					tx.Rollback()
					time.Sleep(100 * time.Millisecond) // Ждем 100 мс перед повторной попыткой
					continue
				}
				http.Error(w, "Ошибка добавления продукта", http.StatusInternalServerError)
				tx.Rollback()
			}

			err = tx.Commit()
			if err == nil {
				fmt.Println("Добавлен продукт", name)
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			}
		}
		fmt.Print(err)
		http.Error(w, "Не удалось добавить продукт после 5 попыток", http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}
func addAdminHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin", http.StatusSeeOther)

	r.ParseForm()
	name, password := r.Form.Get("name"), config.Hash(r.Form.Get("password"))

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Ошибка создания транзакции")
	}
	for i := 0; i < 5; i++ {
		err = addAdmin(name, password, tx)
		if err == nil {
			err = tx.Commit()
			if err == nil {
				fmt.Println("Добавлен администратор", name)
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			} else {
				fmt.Println("Ошибка сохранения изменений addadminhandler, попрбую снова")
			}
		} else {
			fmt.Println("Ошибка добавления администратора addadminhandler, попрбую снова")
		}
	}

}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		http.Redirect(w, r, "/admin", http.StatusSeeOther)

		r.ParseForm()
		name, category := r.Form.Get("name"), r.Form.Get("category")
		var category_id int
		tx, err := db.Begin()
		if err != nil {
			fmt.Println("Ошибка создания транзакции deleteProductHandler")
			return
		}

		rows, _ := tx.Query("SELECT id FROM categories WHERE UPPER(name) = UPPER(?)", category)
		for rows.Next() {
			rows.Scan(&category_id)
			break
		}
		for i := 0; i < 2; i++ {
			_, err := tx.Exec("DELETE FROM chocolate WHERE UPPER(name) = UPPER(?) AND category_id = ?", name, category_id)
			if err != nil{
				fmt.Println("Возникли проблема при удалении продукта deleteProductHandler, пробую еще раз")
				tx.Rollback()
			}
		}
		err = tx.Commit()
		if err != nil{
			fmt.Println("Ошибка сохранения изменений deleteProductHandler", err)
            return
		}
		fmt.Println("Продукт", name, "удален из категории", category)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)

	}
}
func searchProductHandler(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/admin", http.StatusSeeOther)

	search_name := r.URL.Query().Get("name")

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Ошибка создания транзакции searchProductHandler")
		return
	}

	type Product struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}
	products := []Product{}
	var query_name, query_category string
	var query_category_id int

	fmt.Print("Ищем продукты по имени: ", search_name)
	rows, _ := tx.Query("SELECT name, category_id FROM chocolate where UPPER(name) LIKE UPPER(? || '%')", search_name)
	for rows.Next() {
		rows.Scan(&query_name, &query_category_id)
		category_rows, _ := tx.Query("SELECT name FROM categories WHERE id = ?", query_category_id)
		for category_rows.Next() {
			category_rows.Scan(&query_category)
			products = append(products, Product{Name: query_name, Category: query_category})
			break
		}
	}
	fmt.Println()
	err = tx.Commit()
	if err != nil {
		fmt.Println("Ошибка сохранения изменений searchProductHandler")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)

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
	fmt.Println("Заполняю бд")
	fillCategoryDataTable()
	fillAdminDataTable()
	fmt.Println("Бд заполнена")
}

func addChocolate(name string, price int, description string, category string, tx *sql.Tx) error {
	var category_id int
	rows, err := tx.Query("SELECT id from categories where name==(?)", category)
	for rows.Next() {
		rows.Scan(&category_id)
		break
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO chocolate (name, description, price, category_id) VALUES (?, ?, ?, ?)", name, description, price, category_id)
	return err
}
func addCategory(name string, tx *sql.Tx) error {
	_, err = tx.Exec("INSERT INTO categories (name) VALUES (?)", name)
	return err
}
func addAdmin(name string, password string, tx *sql.Tx) error {
	_, err = tx.Exec("INSERT INTO admins (name, password) VALUES (?, ?)", name, password)
	return err
}
func fillCategoryDataTable() {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Ошибка начала транзакции fillcategory")
	}

	categories := []string{"Шоколадные плитки", "Батончики", "Конфеты", "Мороженое", "Пасты и карамели"}
	for _, item := range categories {
		err := addCategory(item, tx)
		if err != nil {
			fmt.Println("Ошибка при выполнении операции addcategory, пытаюсь еще раз")
			tx.Rollback()
		}
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println("Ошибка сохранения изменений fillcategorydatabase")
		return
	}
	fmt.Println("Таблица категорий заполнена")
}
func fillAdminDataTable() {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Ошибка начала транзакции filladmindatatable")
	}
	for i := 0; i < 5; i++ {
		_, err = tx.Exec("INSERT INTO admins (name, password) VALUES (?, ?)", config.FirstAdmin, config.FirstPassword)
		if err == nil {
			break
		}
		fmt.Println("Ошибка выполнения операции filladmindatatable, пытаюсь еще раз")
	}
	err = tx.Commit()
	if err == nil {
		fmt.Println("Таблица админов заполнена")
		return
	}
	fmt.Println("Ошибка сохранения изменений filladmindatatable")
}
func WaitForCommands() {
	var command string
	for {
		fmt.Scan(&command)
		switch command {
		case "s":
			commands.Stop()
		case "filldb":
			fillAllDatabases()
		default:
			fmt.Println("No such command")
		}
	}
}
