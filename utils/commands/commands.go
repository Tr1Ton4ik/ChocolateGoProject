package commands

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strings"
)

func NewProduct() {
	db, _ := sql.Open("sqlite3", "./db/chocolate.db")
	var category_id int
	name, category, description, price := collectProduct()

	rows, _ := db.Query("select id from categories where name==(?)", category)
	rows.Next()
	rows.Scan(&category_id)

	stmnt, _ := db.Prepare("INSERT INTO chocolate (name, price, description, category_id) VALUES (?,?,?,?)")
	_, err := stmnt.Exec(name, price, description, category_id)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Продукт %v успешно добавлен", name)
	}
}
func collectProduct() (string, string, string, int) {
	name := getName()
	category := getCategory()
	description := getDescription()
	price := getPrice()
	fmt.Printf(`Вы создали товар:
Имя: %v
Категория: %v
Описание: %v
Цена: %v руб.
Добавте изображение товара под названием %v.jpeg в папку front/img/
`, name, category, description, price, name)
	return name, category, description, price
}

func getName() string {
	var name string
	fmt.Println("Введите название товара: ")

	reader := bufio.NewReader(os.Stdin)
	name, _ = reader.ReadString('\n')
	name = strings.Trim(name, "\n")

	if name == "s" {
		Stop()
	}
	rune_name := []rune(name)
	name = strings.ToUpper(string(rune_name[0])) + strings.ToLower(string(rune_name[1:]))
	return name
}

func getPrice() int {
	fmt.Println("Введите цену товара в рублях без запятых и пробелов: ")
	var price int
	fmt.Scan(&price)
	return price
}
func getDescription() string {
	fmt.Println("Введите описание товара для того, чтобы программа поняла, о завершении текста в конце описания нажмите tab : ")
	reader := bufio.NewReader(os.Stdin)
	description, _ := reader.ReadString('\t')
	return description
}
func getCategory() string {
	var (
		id             int
		name, category string
	)
	all_categories := []string{}

	db, _ := sql.Open("sqlite3", "./db/chocolate.db")
	rows, _ := db.Query("select * from categories")

	for rows.Next() {
		rows.Scan(&id, &name)
		all_categories = append(all_categories, name)
	}

	for {
		fmt.Println("Введите категорию товара, возможные варианты : ", all_categories)

		reader := bufio.NewReader(os.Stdin)
		category, _ = reader.ReadString('\n')
		category = strings.Trim(category, "\n")
		category_rune := []rune(category)
		category = strings.ToUpper(string(category_rune[0])) + strings.ToLower(string(category_rune[1:]))

		if slices.Contains(all_categories, category) {
			return category
		} else if category == "s" {
			Stop()
		} else {
			fmt.Println("Такой категории нет, попробуйте снова")
		}
	}
}
func Stop() {
	fmt.Print("The stop is complete")
	os.Exit(0)
}
