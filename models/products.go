package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Barcode string `json:"barcode"`
	Image   string `json:"image"`
	Brand   string `json:"brand"`
}

type Data struct {
	Data []ProductsVipcommerce `json:"data"`
}

type ProductsVipcommerce struct {
	Descricao    string `json:"descricao"`
	Preco        string `json:"preco"`
	Imagem       string `json:"imagem"`
	CodigoBarras string `json:"codigo_barras"`
	Marca        string `json:"marca"`
}

type Price struct {
	PriceID   int    `json:"id"`
	ProductID int    `json:"productId"`
	Price     int    `json:"price"`
	CreatedAt string `json:"createdAt"`
}

type ProductWithPrice struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Barcode   string `json:"barcode"`
	Image     string `json:"image"`
	Brand     string `json:"brand"`
	Price     int    `json:"price"`
	CreatedAt string `json:"createdAt"` // se quiser mostrar a data do preço mais recente
}

func GetProductsWithPrices(db *sql.DB) []ProductWithPrice {
	query := `
	SELECT p.id, p.name, p.barcode, p.image, p.brand, pr.price from products p 
	INNER JOIN prices pr ON p.id = pr.productId;
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var products []ProductWithPrice
	for rows.Next() {
		var p ProductWithPrice
		err := rows.Scan(&p.ID, &p.Name, &p.Barcode, &p.Image, &p.Brand, &p.Price)
		if err != nil {
			log.Fatal(err)
		}
		products = append(products, p)
	}

	return products
}
