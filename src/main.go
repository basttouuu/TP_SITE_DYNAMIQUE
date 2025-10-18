package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

type Produit struct {
	Id               int
	Nom              string
	Description      string
	Reduction        float64
	Image            string
	Prix             float64
	Lareduc          bool
	PrixReduit       float64
	PourcentageReduc int
}

var produits = []Produit{
	{
		Id:               1,
		Nom:              "PALACE PULL A CAPUCHE UNISEXE CHASSEUR",
		Description:      "Pull unisexe confortable",
		Reduction:        0.20,
		Image:            "/static/img/products/19A.webp",
		Prix:             129.99,
		Lareduc:          true,
		PrixReduit:       129.99 * (1 - 0.20),
		PourcentageReduc: int(0.20 * 100),
	},
	{
		Id:               2,
		Nom:              "PALACE PULL A CAPUCHON MARINE",
		Description:      "Pull marine stylé",
		Reduction:        0.10,
		Image:            "/static/img/products/21A.webp",
		Prix:             119.00,
		Lareduc:          true,
		PrixReduit:       119.00 * (1 - 0.10),
		PourcentageReduc: int(0.10 * 100),
	},
	{
		Id:               3,
		Nom:              "PALACE PULL CREW PASSEPOSE NOIR",
		Description:      "Pull noir classique",
		Reduction:        0.00,
		Image:            "/static/img/products/22A.webp",
		Prix:             99.50,
		Lareduc:          false,
		PrixReduit:       99.50,
		PourcentageReduc: 0,
	},
	{
		Id:               4,
		Nom:              "PALACE WASHED TERRY 1/4 PLACKET HOOD MOJITO",
		Description:      "Hoodie vert mojito",
		Reduction:        0.15,
		Image:            "/static/img/products/34B.webp",
		Prix:             139.00,
		Lareduc:          true,
		PrixReduit:       139.00 * (1 - 0.15),
		PourcentageReduc: int(0.15 * 100),
	},
	{
		Id:               5,
		Nom:              "PALACE PANTALON BOSSY JEAN STONE",
		Description:      "Jean stone coupe bossy",
		Reduction:        0.05,
		Image:            "/static/img/products/34B.webp",
		Prix:             149.90,
		Lareduc:          true,
		PrixReduit:       149.90 * (1 - 0.05),
		PourcentageReduc: int(0.05 * 100),
	},
	{
		Id:               6,
		Nom:              "PALACE PANTALON CARGO GORE-TEX R-TEK NOIR",
		Description:      "Cargo Gore-Tex noir",
		Reduction:        0.25,
		Image:            "/static/img/products/33B.webp",
		Prix:             199.00,
		Lareduc:          true,
		PrixReduit:       199.00 * (1 - 0.25),
		PourcentageReduc: int(0.25 * 100),
	},
}

func main() {

	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Println("Erreur template:", err)
		os.Exit(1)
	}

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		err := temp.ExecuteTemplate(w, "home", produits)
		if err != nil {
			http.Error(w, "Erreur template home", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/produit", func(w http.ResponseWriter, r *http.Request) {
		idProduit := r.FormValue("id")
		produitId, err := strconv.Atoi(idProduit)
		if err != nil {
			http.Error(w, "Erreur: id du produit invalide", http.StatusBadRequest)
			return
		}

		for _, product := range produits {
			if product.Id == produitId {
				temp.ExecuteTemplate(w, "produit", product)
				return
			}
		}

		http.Error(w, "Produit non trouvé", http.StatusNotFound)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	})

	fs := http.FileServer(http.Dir("./../assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Serveur lancé sur : http://localhost:8008/home")
	if err := http.ListenAndServe("localhost:8008", nil); err != nil {
		fmt.Println("Erreur serveur:", err)
		os.Exit(1)
	}
}
