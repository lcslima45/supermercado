package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Data struct {
	Data []Products `json:"data"`
}

type Products struct {
	Descricao    string `json:"descricao"`
	Preco        string `json:"preco"`
	Imagem       string `json:"imagem"`
	CodigoBarras string `json:"codigo_barras"`
	Marca        string `json:"marca"`
}

func TreatImage(imageUrl string) string {
	return "https://produtos.vipcommerce.com.br/250x250/" + imageUrl
}

func TreatPrice(price string) string {
	return strings.ReplaceAll(price, ".", "")
}

func TreatData(dat *Data) {
	for i := range dat.Data {
		dat.Data[i].Preco = TreatPrice(dat.Data[i].Preco)
		dat.Data[i].Imagem = TreatImage(dat.Data[i].Imagem)
	}
}

func SaveData(dat Data) {
	file, _ := os.Create("products.json")
	defer file.Close()
	encoder := json.NewEncoder(file)
	err := encoder.Encode(dat)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpc3MiOiJ2aXBjb21tZXJjZSIsImF1ZCI6ImFwaS1hZG1pbiIsInN1YiI6IjZiYzQ4NjdlLWRjYTktMTFlOS04NzQyLTAyMGQ3OTM1OWNhMCIsInZpcGNvbW1lcmNlQ2xpZW50ZUlkIjpudWxsLCJpYXQiOjE3NTI3OTQwNDAsInZlciI6MSwiY2xpZW50IjpudWxsLCJvcGVyYXRvciI6bnVsbCwib3JnIjoiMTgwIn0.LHqEKsuU5KOj9OQaOK03HAdUph9JR8FLIT1OEsp3Cp1UJjvolb75o5aN_ZmNtgGenvhOcn-TWW1xYeXRlMSScQ")
		r.Headers.Set("Content-Type", "application/json")
		r.Headers.Set("Domainkey", "lojaonline.pinheirosupermercado.com.br")
		r.Headers.Set("Referer", "https://lojaonline.pinheirosupermercado.com.br/")
		r.Headers.Set("Organizationid", "180")
	})

	var result Data
	c.OnResponse(func(r *colly.Response) {
		err := json.Unmarshal(r.Body, &result)
		if err != nil {
			log.Fatal(err)
		}
	})

	err := c.Visit("https://services.vipcommerce.com.br/api-admin/v1/org/180/filial/1/centro_distribuicao/1/loja/classificacoes_mercadologicas/departamentos/4/produtos?page=2&")

	TreatData(&result)

	SaveData(result)

	if err != nil {
		log.Fatal(err)
	}

}
