package commands

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strings"
)

func CollectProduct() {
	name := getName()
	category := getCategory()
	description := getDescription()
	price := getPrice()
	fmt.Println("Вы выбрали", name, category, description, price)
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
	fmt.Println("Введите описание товара для того, чтобы программа поняла, о завершении текста в конце описания добавьте фразу qwerty123# : ")
	return ""
}
func getCategory() string {
	var (
		id             int
		name, category string
	)
	all_categories := []string{}
	reader := bufio.NewReader(os.Stdin)

	db, _ := sql.Open("sqlite3", "./db/chocolate.db")
	rows, _ := db.Query("select * from categories")

	for rows.Next() {
		rows.Scan(&id, &name)
		all_categories = append(all_categories, name)
	}

	for {
		fmt.Println("Введите категорию товара, возможные варианты : ", all_categories)

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
