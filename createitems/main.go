package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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
	INSERT INTO products(name, barcode, image, brand) VALUES ("coca cola", "7894900709841", "https://cdn-cosmos.bluesoft.com.br/products/7894900709841", "coca-cola")
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
