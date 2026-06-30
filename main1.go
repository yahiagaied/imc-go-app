package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func accueil(w http.ResponseWriter, r *http.Request) {
	// Charger le template HTML pour la page d'accueil
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Exécuter le template sans données
	tmpl.Execute(w, nil)
}

func traitement(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode est POST
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Récupérer les valeurs du formulaire
	tailleStr := r.FormValue("taille")
	poidsStr := r.FormValue("poids")
	age := r.FormValue("age")
	sexe := r.FormValue("sexe")

	// Convertir taille et poids en float64
	taille, err1 := strconv.ParseFloat(tailleStr, 64)
	poids, err2 := strconv.ParseFloat(poidsStr, 64)

	if err1 != nil || err2 != nil || taille <= 0 || poids <= 0 {
		http.Error(w, "Veuillez entrer des valeurs valides pour la taille et le poids.", http.StatusBadRequest)
		return
	}

	// Calculer l'IMC : poids (kg) / (taille (m))^2
	imc := poids / (taille * taille)

	// Déterminer la catégorie d'IMC
	var categorie string
	switch {
	case imc < 18.5:
		categorie = "Sous-poids"
	case imc >= 18.5 && imc < 25:
		categorie = "Poids normal"
	case imc >= 25 && imc < 30:
		categorie = "Surpoids"
	default:
		categorie = "Obésité"
	}

	// Charger le template de traitement
	tmpl, err := template.ParseFiles("templates/traitement.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Préparer les données à passer au template
	data := struct {
		Age       string
		Sexe      string
		Imc       float64
		Categorie string
	}{
		Age:       age,
		Sexe:      sexe,
		Imc:       imc,
		Categorie: categorie,
	}

	// Exécuter le template avec les résultats
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", accueil)
	http.HandleFunc("/traitement", traitement)
	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
