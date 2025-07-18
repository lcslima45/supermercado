package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"supermercado/models"
	"text/template"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStatement := `
	DROP TABLE IF EXISTS products;

	CREATE TABLE products (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		barcode TEXT,
		image TEXT,
		brand TEXT,
		CONSTRAINT barcode_validate CHECK (
			barcode GLOB '[0-9]*' AND LENGTH(barcode) >= 6 AND LENGTH(barcode) <= 15
		)
	);
	`

	_, err = db.Exec(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}

	sqlStatement = `
	DROP TABLE IF EXISTS prices;
	CREATE TABLE prices (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		productId INTEGER, 
		price INTEGER,
		createdAt TEXT DEFAULT (DATETIME('now')),
		FOREIGN KEY (productId) REFERENCES products(id)
	);
	`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}

	file, _ := os.ReadFile("../scraper/products.json")

	var data models.Data

	err = json.Unmarshal(file, &data)

	if err != nil {
		log.Fatal(err)
	}

	for _, p := range data.Data {
		sqlStatement := fmt.Sprintf(`
		INSERT INTO products(name, barcode, image, brand)
		VALUES ("%s", "%s", "%s", "%s");`,
			p.Descricao,
			p.CodigoBarras,
			p.Imagem,
			p.Marca,
		)
		result, err := db.Exec(sqlStatement)
		if err != nil {
			log.Fatal(err)
		}

		lastId, err := result.LastInsertId()

		if err != nil {
			log.Println(err)
		}

		sqlStatement = fmt.Sprintf(`
		INSERT INTO prices(productId, price, createdAt)
		VALUES (%d, %s, DATETIME('now'));`,
			lastId,
			p.Preco, // já está em centavos como string: "949"
		)

		_, err = db.Exec(sqlStatement)
		if err != nil {
			log.Fatal(err)
		}
	}
	sqlStatement = `
	INSERT INTO prices(productId, price, createdAt) 
	VALUES (1, 630, DATETIME('now'));
	 `
	_, err = db.Exec(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Tabela 'products' e 'prices' criadas com sucesso e dados inseridos.")

	log.Println("Table 'products' created succesfully")

	prd := models.GetProductsWithPrices(db)

	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"formatPrice": func(cents int) string {
			return fmt.Sprintf("%.2f", float64(cents)/100)
		},
	})

	r.LoadHTMLGlob("../templates/template.html")

	r.GET("/products", func(c *gin.Context) {
		if err != nil {
			c.String(http.StatusInternalServerError, "Erro ao buscar produtos")
			return
		}
		c.HTML(http.StatusOK, "template.html", gin.H{
			"Products": prd,
		})
	})
	log.Println(prd)

	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}
