package components

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var FilePath string
var ExtractFilePath string
var CompleteExtractFileName string
var CityINSEE string
var DepartID string

// ReadCSVFileContentAndExtracter lit le fichier CSV source et extrait les données en fonction de la configuration
func (wi *WorkerImmeuble) ReadCSVFileContentAndExtracter() {
	// On initialise le temps qui servira plus tard pour l'horodatage du fichier exporté
	now := time.Now()

	// Récupération du chemin du fichier CSV (SOURCE) depuis la configuration
	FilePath = wi.Config.File_immeuble

	// Récupération des listes de codes INSEE et de départements depuis la configuration
	InseeList := wi.Config.Lst_Insee
	DepartList := wi.Config.Lst_Dprt

	// Récupérer le chemin d'extraction depuis la configuration
	ExtractFilePath := wi.Config.File_export

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

	// Création du nom du fichier d'extraction
	var CompleteExtractFileName string
	CityINSEE := strings.Join(InseeList, ",") // Reconstituer la chaîne de codes INSEE
	DepartID := strings.Join(DepartList, ",") // Reconstituer la chaîne de départements

	if CityINSEE == "" {
		CompleteExtractFileName = "Extraction_du_" + now.Format("2006-01-02_15-04-05") + "_PAR_DPT_" + DepartID
	} else {
		CompleteExtractFileName = "Extraction_du_" + now.Format("2006-01-02_15-04-05") + "_PAR_INSEE_" + CityINSEE
	}

	// Si les deux variables ne sont pas vides
	if CityINSEE != "" && DepartID != "" {
		CompleteExtractFileName = "Extraction_du_" + now.Format("2006-01-02_15-04-05") + "_PAR_INSEE_" + CityINSEE + "_PAR_DPT_" + DepartID
	}

	// Création du fichier d'extraction
	csvExtractedFile, err := os.Create(ExtractFilePath + CompleteExtractFileName + ".csv")
	if err != nil {
		fmt.Printf("Erreur survenue lors de la création du fichier Extraction : %v\n", err)
		return
	}
	defer csvExtractedFile.Close() // On s'assure de sa fermeture à la fin de la fonction

	// Créer un writer pour écrire dans le fichier d'extraction
	w := csv.NewWriter(csvExtractedFile)
	w.Comma = ';'   // Définir le séparateur pour le writer du fichier d'extraction
	defer w.Flush() // On s'assure que ce qui doit être écrit est écrit à la fin de la fonction

	// Écrire le header dans le fichier d'extraction
	header := []string{"x", "y", "imb_id", "num_voie", "cp_no_voie", "type_voie", "nom_voie", "batiment", "code_insee", "code_poste", "nom_com", "catg_loc_imb", "imb_etat", "pm_ref", "pm_etat", "code_l331", "geom_mod", "type_imb"}
	w.Write(header)

	for {

		// Lire une ligne du fichier source
		record, err := r.Read()

		// Si l'erreur est la fin du fichier, on arrête
		if err == io.EOF {
			break
		}

		// Si une erreur est rencontrée, on continue quand même
		if err != nil {
			fmt.Printf("Erreur lors de la lecture d'une ligne suivante : %v\n", err)
			continue
		}

		// Pour chaque élément de la boucle, on vérifie s'il a au moins 10 colonnes
		if len(record) >= 10 {
			// Si la colonne 9 (Code Postal) contient 2 ou plus caractères
			if len(record[9]) >= 2 {
				codeDept := strings.TrimSpace(record[9][:2]) // Récupérer le département

				// Boucle pour parcourir la liste des départements à extraire
				for _, dept := range DepartList {
					// Si le département de l'enregistrement correspond au département à extraire
					if codeDept == strings.TrimSpace(dept) {

						// Écriture de l'enregistrement trouvé
						w.Write(record)
						break
					}
				}
			} else {
				continue
			}

			// On définit le code INSEE de l'enregistrement en cours
			codeInsee := strings.TrimSpace(record[8]) // 9e colonne

			// Boucle pour parcourir la liste des codes INSEE à extraire
			for _, insee := range InseeList {
				// Si le code INSEE de l'enregistrement correspond au code INSEE à extraire
				if codeInsee == strings.TrimSpace(insee) {

					// Écriture de l'enregistrement trouvé
					w.Write(record)
					break
				}
			}
		}
	}
	fmt.Printf("Extraction terminée, le résultat est disponible dans le fichier : %s\n", ExtractFilePath+CompleteExtractFileName+".csv")
}
