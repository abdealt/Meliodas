package components

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Exportation de variables
var FilePath string
var ExtractFilePath string
var CompleteExtractFileName string
var CityINSEE string
var DepartID string

// ReadCSVFileContent charge et lit le fichier CSV spécifié, puis compare la 9e colonne avec CITY_INSEE
// et exporte les lignes correspondantes dans un nouveau fichier CSV
func ReadCSVFileContent() {
	// On initialise le temps
	now := time.Now()

	// Charger les variables d'environnement à partir du fichier .env
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Erreur lors du chargement du fichier .env : %v\n", err)
		return
	}

	// Récupérer le chemin du fichier CSV depuis les variables d'environnement
	FilePath = os.Getenv("SOURCE_FILE")
	if FilePath == "" {
		fmt.Printf("Aucun chemin de fichier fourni dans le fichier de configuration .env\n")
		return
	}

	// Récupérer la variable CITY_INSEE ou DEPARTMENT_ID depuis le fichier .env
	CityINSEE = os.Getenv("CITY_INSEE")
	DepartID = os.Getenv("DEPARTMENT_ID")
	if CityINSEE == "" && DepartID == "" {
		fmt.Printf("Aucun code INSEE (CITY_INSEE) ou Département n'est fourni dans le fichier de configuration .env\n")
		return
	}

	// Récupérer le chemin du fichier d'extraction depuis les variables d'environnement
	ExtractFilePath := os.Getenv("EXTRACT_FILE")
	if ExtractFilePath == "" {
		fmt.Printf("Aucun chemin pour le fichier d'extraction (EXTRACT_FILE) n'est fourni dans le fichier de configuration .env\n")
		return
	}

	// Ouvrir le fichier CSV (Source)
	csvFile, err := os.Open(FilePath)
	if err != nil {
		fmt.Printf("Erreur lors de l'ouverture du fichier : %v\n", err)
		return
	}
	defer csvFile.Close()

	// Création d'un reader, élément qui permet la lecture d'enregistrement
	r := csv.NewReader(csvFile)
	r.Comma = ',' // Définit le séparateur

	if CityINSEE == "" {
		CompleteExtractFileName = "Extraction_du_" + now.Format("2006-01-02_15-04-05") + "_PAR_DPT_" + DepartID
	} else {
		CompleteExtractFileName = "Extraction_du_" + now.Format("2006-01-02_15-04-05") + "_PAR_INSEE_" + CityINSEE
	}

	// Création du fichier CSV (Extraction)
	csvExtractedFile, err := os.Create(ExtractFilePath + CompleteExtractFileName + ".csv")
	if err != nil {
		fmt.Printf("Erreur lors de l'ouverture du fichier : %v\n", err)
		return
	}
	// On s'assur de sa fermeture a la fin de la fonction
	defer csvExtractedFile.Close()

	// Création d'un writer, élément qui permet d'écrire
	w := csv.NewWriter(csvExtractedFile)
	// On s'assure que ce qui doit etre ecris est ecris a la fin de la fonction
	defer w.Flush()

	// Definition du header du nouveau csv (la ou les données seront extraites)
	header := []string{"x", "y", "imb_id", "num_voie", "cp_no_voie", "type_voie", "nom_voie", "batiment", "code_insee", "code_poste", "nom_com", "catg_loc_imb", "imb_etat", "pm_ref", "pm_etat", "code_l331", "geom_mod", "type_imb"}
	w.Write(header)

	// Création de la boucle de lecture dans le fichier SOURCE_FILE
	for {
		// Lire une ligne dans le fichier CSV (Source)
		record, err := r.Read()

		// On vérifie si on est a la fin du fichier
		if err == io.EOF {
			// Fin du fichier atteinte
			break
		}

		// On verifie si une erreur est présente
		if err != nil {
			fmt.Printf("Erreur lors de la lecture d'une ligne : %v\n", err)
			continue // Passe à la ligne suivante
		}

		// On vérifie si la ligne en cours contient suffisamment de colonnes (au moins 10)
		if len(record) >= 10 {
			codeInsee := strings.TrimSpace(record[8])    // 9e colonne
			codeDept := strings.TrimSpace(record[9][:2]) // 10e colonne, on recupère les 2 premier caractères

			if codeInsee == CityINSEE || codeDept == DepartID {
				// Ligne correspondante trouvée, on l'écrit dans le nouveau CSV
				w.Write(record)
				continue // Passe à la ligne suivante
			} else {
				continue // Passe à la ligne suivante
			}
		}
	}
	fmt.Printf("Extraction terminée, le résultat est disponible dans le fichier : %s\n", ExtractFilePath+CompleteExtractFileName+".csv")
}
