package components

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	// "github.com/joho/godotenv"
)

// Déclaration des variables
var FilePath string
var ExtractFilePath string
var CompleteExtractFileName string
var CityINSEE string
var DepartID string
var ComptTotal int
var ComptElement int
var Stat string
var codeDept string

func ReadCSVFileContentAndExtracter() {
	// On initialise le temps qui servira plutard pour l'horodatage du fichier exporté
	now := time.Now()

	/*
		Les variables du fichier .env sont chargé grace à cmd/root.go (depuis le projet du CLI)
	*/

	// Récupéreration le chemin du fichier CSV (SOURCE) depuis les variables d'environnement
	FilePath = os.Getenv("SOURCE_FILE")
	if FilePath == "" {
		fmt.Printf("Aucun chemin de fichier fourni dans le fichier de configuration .env\n")
		return
	}

	// Récupéreration la variable CITY_INSEE ou DEPARTMENT_ID depuis le fichier .env
	CityINSEE = os.Getenv("CITY_INSEE")
	DepartID = os.Getenv("DEPARTMENT_ID")
	if CityINSEE == "" && DepartID == "" {
		fmt.Printf("Aucun code INSEE (CITY_INSEE) ou Département n'est fourni dans le fichier de configuration .env\n")
		return
	}

	// Séparation des INSEE et des Département (dans le cas ou il y a plusieurs qui sont saisies)
	DepartList := strings.Split(DepartID, ",")
	InseeList := strings.Split(CityINSEE, ",")

	// Récupérer le chemin d'extraction depuis les variables d'environnement (La ou sera extrait le fichier)
	ExtractFilePath := os.Getenv("EXTRACT_FILE")
	if ExtractFilePath == "" {
		fmt.Printf("Aucun chemin pour le fichier d'extraction (EXTRACT_FILE) n'est fourni dans le fichier de configuration .env\n")
		return
	}

	// Ouvrerture du fichier CSV (Source)
	csvFile, err := os.Open(FilePath)
	if err != nil {
		fmt.Printf("Erreur survenue lors de l'ouverture du fichier Source : %v\n", err)
		return
	}
	// On s'assure de sa fermeture a la fin de la fonction
	defer csvFile.Close()

	// Création d'un reader pour lire le fichier source
	r := csv.NewReader(csvFile)
	// Définit le séparateur du reader du fichier source
	r.Comma = ','

	// Création du nom du fichier d'extractions
	if CityINSEE == "" {
		CompleteExtractFileName = "Extraction_du_" + now.Format("2006-01-02_15-04-05") + "_PAR_DPT_" + DepartID
	} else {
		CompleteExtractFileName = "Extraction_du_" + now.Format("2006-01-02_15-04-05") + "_PAR_INSEE_" + CityINSEE
	}

	// Si les deux variable ne sont pas vides
	if CityINSEE != "" && DepartID != "" {
		CompleteExtractFileName = "Extraction_du_" + now.Format("2006-01-02_15-04-05") + "_PAR_INSEE_" + CityINSEE + "_PAR_DPT_" + DepartID
	}

	// Création du fichier d'extraction
	csvExtractedFile, err := os.Create(ExtractFilePath + CompleteExtractFileName + ".csv")
	if err != nil {
		fmt.Printf("Erreur survenue lors de l'ouverture du fichier Extraction : %v\n", err)
		return
	}
	// On s'assur de sa fermeture a la fin de la fonction
	defer csvExtractedFile.Close()

	// Créer un writer pour écrire dans le fichier d'extraction
	w := csv.NewWriter(csvExtractedFile)
	// Définir le séparateur pour le writer du fichier d'extraction
	w.Comma = ';'

	// On s'assure que ce qui doit etre écrit est ecris a la fin de la fonction
	defer w.Flush()

	// Écrire le header dans le fichier d'extraction
	header := []string{"x", "y", "imb_id", "num_voie", "cp_no_voie", "type_voie", "nom_voie", "batiment", "code_insee", "code_poste", "nom_com", "catg_loc_imb", "imb_etat", "pm_ref", "pm_etat", "code_l331", "geom_mod", "type_imb"}
	w.Write(header)

	// Boucle pour lire et traiter chaque enregistrement
	for {
		// Incrémentation du compteur du total d'élements
		ComptTotal++

		// Lire une ligne du fichier source
		record, err := r.Read()

		// Si l'érreur est la fin du fichier alors on arrête
		if err == io.EOF {
			break
		}

		// Si un erreur est rencontrait on continue quand même
		if err != nil {
			fmt.Printf("Erreur lors de la lecture d'une ligne suivante : %v\n", err)
			continue
		}

		// Pour chaque élément de la boucle on vérifier si il a au moins 10 colonnes
		if len(record) >= 10 {
			// Si la colonne 9 (Code Postale) contient 2 ou plus caractères
			if len(record[9]) >= 2 {
				// On récupère le département dans une variable
				codeDept = strings.TrimSpace(record[9][:2])

				// Boucle pour parcourir la liste des départements à éxtraire
				for _, dept := range DepartList {
					// Si le département de l'enregistrement correspond au département à extraire
					if codeDept == strings.TrimSpace(dept) {
						// Incrémentation de element
						ComptElement++
						// Ecriture de l'enregistrement trouvé
						w.Write(record)
						break
					}
				}
			} else {
				continue
			}

			// On definit le code INSEE de l'enregistrement en cours
			codeInsee := strings.TrimSpace(record[8]) // 9e colonne

			// Boucle pour parcourir la liste des codes INSEE à éxtraire
			for _, insee := range InseeList {
				// Si le code INSEE de l'enregistrement correspond au code INSEE à extraire
				if codeInsee == strings.TrimSpace(insee) {
					// Incrémentation de element
					ComptElement++
					// Ecriture de l'enregistrement trouvé
					w.Write(record)
					break
				}
			}
		}
	}
	fmt.Printf("Extraction terminée, le résultat est disponible dans le fichier : %s\n", ExtractFilePath+CompleteExtractFileName+".csv")
}
