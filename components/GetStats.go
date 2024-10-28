package components

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// ExtractStatisticsFromCSV lit le fichier CSV source et extrait les données en fonction de la configuration
func (wi *WorkerImmeuble) ExtractStatisticsFromCSV() {
	// Récupération du chemin du fichier CSV (SOURCE) depuis la configuration
	FilePath := wi.Config.File_immeuble

	// Récupération des listes de codes INSEE et de départements depuis la configuration
	InseeList := wi.Config.Lst_Insee
	DepartList := wi.Config.Lst_Dprt

	// Ouverture du fichier CSV (Source)
	csvFile, err := os.Open(FilePath)
	if err != nil {
		fmt.Printf("Erreur survenue lors de l'ouverture du fichier Source : %v\n", err)
		return
	}
	defer csvFile.Close() // On s'assure de sa fermeture à la fin de la fonction

	// Création d'un reader pour lire le fichier source
	r := csv.NewReader(csvFile)
	r.Comma = ',' // Définit le séparateur du reader du fichier source

	// Initialisation des compteurs
	var comptTotal int
	var comptElement int

	// Boucle pour lire et traiter chaque enregistrement
	for {
		// Lire une ligne du fichier source
		record, err := r.Read()

		// Si l'erreur est la fin du fichier, on arrête
		if err == io.EOF {
			break
		}

		// Si une erreur est rencontrée, on arrête
		if err != nil {
			fmt.Printf("Erreur lors de la lecture d'une ligne suivante : %v\n", err)
			return // Sortir de la fonction en cas d'erreur
		}

		// Incrémentation du compteur du total d'éléments
		comptTotal++

		// Vérification du nombre de colonnes
		if len(record) < 10 {
			continue // Passer à l'enregistrement suivant si moins de 10 colonnes
		}

		// Vérification du code postal
		codeDept := strings.TrimSpace(record[9])
		if len(codeDept) < 2 {
			continue // Passer à l'enregistrement suivant si le code postal est invalide
		}

		// Récupérer le département
		codeDept = codeDept[:2] // On garde seulement les deux premiers caractères

		// Vérification du département
		for _, dept := range DepartList {
			if codeDept == strings.TrimSpace(dept) {
				comptElement++ // Incrémentation de l'élément pour le département
				break
			}
		}

		// Vérification du code INSEE
		codeInsee := strings.TrimSpace(record[8]) // 9e colonne
		for _, insee := range InseeList {
			if codeInsee == strings.TrimSpace(insee) {
				comptElement++ // Incrémentation de l'élément pour le code INSEE
				break
			}
		}
	}

	// Impression des résultats
	fmt.Printf("Total d'éléments traités : %d\n Éléments extraits : %d", comptElement, comptTotal)
}
