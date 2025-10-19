package main

import (
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type Produit struct {
	Id               int
	Nom              string
	Description      string
	Prix             float64
	Reduction        float64
	Image            string
	Lareduc          bool
	PrixReduit       float64
	PourcentageReduc int
}

var imagePath string

var (
	produits = []Produit{
		{
			Id:               1,
			Nom:              "PALACE PULL A CAPUCHE UNISEXE CHASSEUR",
			Description:      "Pull unisexe confortable",
			Reduction:        0.20,
			Image:            "/static/img/products/19A.webp",
			Prix:             129.99,
			Lareduc:          true,
			PrixReduit:       129.99 * (1 - 0.20),
			PourcentageReduc: 20,
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
			PourcentageReduc: 10,
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
			Image:            "/static/img/products/16A.webp",
			Prix:             139.00,
			Lareduc:          true,
			PrixReduit:       139.00 * (1 - 0.15),
			PourcentageReduc: 15,
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
			PourcentageReduc: 5,
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
			PourcentageReduc: 25,
		},
		{
			Id:               7,
			Nom:              "PALACE PULL A CAPUCHE UNISEXE LONDON NOIR",
			Description:      "Pull Noir stylé",
			Reduction:        0,
			Image:            "/static/img/products/18A.webp",
			Prix:             259.00,
			Lareduc:          false,
			PrixReduit:       259.00,
			PourcentageReduc: 0,
		},
	}
	nextID = 8
)

func main() {

	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Println("Erreur template:", err)
		os.Exit(1)
	}

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		err := temp.ExecuteTemplate(w, "home", produits)
		if err != nil {
			temp.ExecuteTemplate(w, "erreurtemplates", nil)
		}
	})

	http.HandleFunc("/erreur", func(w http.ResponseWriter, r *http.Request) {
		err := temp.ExecuteTemplate(w, "erreur", produits)
		if err != nil {
			temp.ExecuteTemplate(w, "erreurtemplates", nil)
		}
	})

	http.HandleFunc("/produit", func(w http.ResponseWriter, r *http.Request) {
		idProduit := r.FormValue("id")
		produitId, err := strconv.Atoi(idProduit)
		if err != nil {
			temp.ExecuteTemplate(w, "notfound", nil)
			return
		}

		for _, product := range produits {
			if product.Id == produitId {
				temp.ExecuteTemplate(w, "produit", product)
				return
			}
		}

		temp.ExecuteTemplate(w, "notfound", nil)
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			nom := r.FormValue("nom")
			description := r.FormValue("description")
			prixStr := r.FormValue("prix")
			reducStr := r.FormValue("reduction")

			file, header, err := r.FormFile("image")
			if err != nil {
				header = &multipart.FileHeader{Filename: ""}
			} else {
				defer file.Close()
				importationimagequirendfou(file, header.Filename)
			}

			if nom == "" || description == "" || prixStr == "" {
				http.Error(w, "Champs obligatoires manquants", http.StatusBadRequest)
				return
			}

			prix, err := strconv.ParseFloat(prixStr, 64)
			if err != nil {
				http.Error(w, "Prix invalide", http.StatusBadRequest)
				return
			}

			reduc := 0.0
			if reducStr != "" {
				reduc, err = strconv.ParseFloat(reducStr, 64)
				if err != nil {
					http.Error(w, "Réduction invalide", http.StatusBadRequest)
					return
				}
			}

			if header.Filename != "" {
				imagePath = "/static/img/products/" + header.Filename
			} else {
				imagePath = "/static/img/products/default.png"
			}

			produit := Produit{
				Id:               nextID,
				Nom:              nom,
				Description:      description,
				Reduction:        reduc,
				Image:            imagePath,
				Prix:             prix,
				Lareduc:          reduc > 0,
				PrixReduit:       prix * (1 - reduc),
				PourcentageReduc: int(reduc * 100),
			}

			produits = append(produits, produit)
			nextID++

			http.Redirect(w, r, fmt.Sprintf("/produit?id=%d", produit.Id), http.StatusSeeOther)
			return
		}

		err := temp.ExecuteTemplate(w, "add", nil)
		if err != nil {
			temp.ExecuteTemplate(w, "notfound", nil)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	})

	fileServer := http.FileServer(http.Dir("./../assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	fmt.Println("Serveur lancé sur : http://localhost:8008/home")
	if err := http.ListenAndServe("localhost:8008", nil); err != nil {
		fmt.Println("Erreur serveur:", err)
		os.Exit(1)
	}
}

func importationimagequirendfou(fichier multipart.File, fichiernom string) {

	imageSouhaiter, err := os.Create("../assets/img/products/" + fichiernom)
	if err != nil {
		fmt.Println("Erreur création destination:", err)
		return
	} else {
		fmt.Println("Import d'image Réussis !")
	}
	defer imageSouhaiter.Close()

	_, err = io.Copy(imageSouhaiter, fichier)
	if err != nil {
		fmt.Println("Erreur copie:", err)
		return
	} else {
		fmt.Println("Copie d'image Réussis !")
	}
}
