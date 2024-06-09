package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Ao utilizar a instrução gorm.Model, os campos ID, CreatedAt, UpdatedAt e DeletedAt são gerados automáticamente
// além disso ao utilizar o comando de db.Delete o registro apenas será marcado com a data que isso aconteceu
type Produto struct {
	Nome         string
	Preco        float64
	CategoriaID  int
	Categoria    Categoria
	NumeroSerial NumeroSerial
	gorm.Model
}

type Categoria struct {
	Nome string
	gorm.Model
}

type NumeroSerial struct {
	Numero    string
	ProdutoID int
	gorm.Model
}

func main() {
	// Para trabalhar com data e hora aparentemente é necessário utilizar configuração avançada na conexão parseTime
	dsn := "root:root@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Produto{}, &Categoria{}, &NumeroSerial{})

	portatil := Categoria{Nome: "Portátil"}
	db.Save(&portatil)

	perifericos := Categoria{Nome: "Periféricos"}
	db.Save(&perifericos)

	// Passando um slice de produtos
	produtos := []Produto{
		{Nome: "Mouse", Preco: 160.00, Categoria: perifericos},
		{Nome: "Monitor", Preco: 890.00, Categoria: perifericos},
		{Nome: "Placa mãe", Preco: 1280.00, Categoria: perifericos},
	}

	db.Create(&produtos)

	//Criando um produto
	db.Create(&Produto{
		Nome:      "Notebook",
		Preco:     2560.30,
		Categoria: portatil,
	})

	// Ao criar ele já retorna os valores do banco atualizados, inclusive com o ID gerado
	for _, p := range produtos {
		fmt.Printf("ID: %v , Nome: %v\n", p.ID, p.Nome)
	}

	var prod1 Produto
	db.First(&prod1, 2)
	fmt.Println(&prod1)

	var prod2 Produto
	db.First(&prod2, "nome = ?", "Monitor")
	fmt.Println(&prod2)

	var lista []Produto
	db.Find(&lista)

	for _, p := range lista {
		fmt.Printf("ID: %v , Nome: %v\n", p.ID, p.Nome)
	}

	produtos = []Produto{}
	// Limitanto a quantidade de retorno
	db.Limit(2).Find(&produtos)
	for _, p := range produtos {
		fmt.Printf("ID: %v , Nome: %v\n", p.ID, p.Nome)
	}
	fmt.Println("======")
	// Paginando
	db.Limit(2).Offset(2).Find(&produtos)
	for _, p := range produtos {
		fmt.Printf("ID: %v , Nome: %v\n", p.ID, p.Nome)
	}

	fmt.Println("======")
	// Aplicando uma condição para a busca
	db.Where("preco > ?", 500).Find(&produtos)
	for _, p := range produtos {
		fmt.Printf("ID: %v , Nome: %v, Preco: %.2f\n", p.ID, p.Nome, p.Preco)
	}

	fmt.Println("======")
	// Aplicando uma condição para a busca
	db.Where("nome LIKE ?", "%Mo%").Find(&produtos)
	for _, p := range produtos {
		fmt.Printf("ID: %v , Nome: %v, Preco: %.2f\n", p.ID, p.Nome, p.Preco)
	}

	prod1 = Produto{}
	db.First(&prod1, "nome = ?", "Mouse")
	prod1.Nome = "Mouse Razer"
	db.Save(&prod1)

	fmt.Println("======")
	prod2 = Produto{}
	db.First(&prod2, "nome LIKE ?", "%Mouse%")
	fmt.Printf("Nome: %v, Preco: %.2f\n", prod2.Nome, prod2.Preco)

	var periferico Categoria
	db.Find(&periferico, "nome = ?", "Periféricos")

	produtos = []Produto{}
	fmt.Println("======")
	fmt.Println(periferico.Nome)
	db.Where("categoria_id = ?", periferico.ID).Find(&produtos)
	for _, p := range produtos {
		fmt.Printf("ID: %v , Nome: %v, Preco: %.2f\n", p.ID, p.Nome, p.Preco)
	}

	fmt.Println("======")
	fmt.Println("Todos os Produtos")
	produtos = []Produto{}
	db.Preload("Categoria").Find(&produtos)
	for _, p := range produtos {
		fmt.Printf("Produto: %v, Categoria: %v\n", p.Nome, p.Categoria.Nome)
		if p.NumeroSerial.ProdutoID == 0 {
			var numeroSerial = NumeroSerial{}
			numeroSerial.Numero = fmt.Sprintf("%d%d", p.ID, p.CategoriaID)
			numeroSerial.ProdutoID = int(p.ID)
			db.Save(&numeroSerial)
		}
	}

	fmt.Println("======")
	fmt.Println("Todos os Produtos com Numero Serial")
	db.Preload("NumeroSerial").Find(&produtos)
	for _, p := range produtos {
		fmt.Printf("ID: %v , Nome: %v, Serial: %v\n", p.ID, p.Nome, p.NumeroSerial.Numero)
	}
}
