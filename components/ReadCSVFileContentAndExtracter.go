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

func ReadCSVFileContentAndExtracter() {
	// On initialise le temps qui servira plutard pour l'horodatage du fichier exporté
	now := time.Now()

	// Charger les variables d'environnement à partir du fichier .env
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Printf("Erreur lors du chargement du fichier .env : %v\n", err)
	// 	return
	// }

	// Récupérer le chemin du fichier CSV (SOURCE) depuis les variables d'environnement
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

	// Séparation des INSEE et des Département (dans le cas ou il y a plusieurs qui sont saisies)
	DepartList := strings.Split(DepartID, ",")
	InseeList := strings.Split(CityINSEE, ",")

	// Récupérer le chemin du fichier d'extraction depuis les variables d'environnement
	ExtractFilePath := os.Getenv("EXTRACT_FILE")
	if ExtractFilePath == "" {
		fmt.Printf("Aucun chemin pour le fichier d'extraction (EXTRACT_FILE) n'est fourni dans le fichier de configuration .env\n")
		return
	}

	// Ouvrir le fichier CSV (Source)
	csvFile, err := os.Open(FilePath)
	if err != nil {
		fmt.Printf("Erreur survenue lors de l'ouverture du fichier Source : %v\n", err)
		return
	}
	// On s'assur de sa fermeture a la fin de la fonction
	defer csvFile.Close()

	// Création d'un reader pour lire le fichier source
	r := csv.NewReader(csvFile)
	r.Comma = ',' // Définit le séparateur du fichier source

	// Déterminer le nom complet du fichier d'extraction
	if CityINSEE == "" {
		CompleteExtractFileName = "Extraction_du_" + now.Format("2006-01-02_15-04-05") + "_PAR_DPT_" + DepartID
	} else {
		CompleteExtractFileName = "Extraction_du_" + now.Format("2006-01-02_15-04-05") + "_PAR_INSEE_" + CityINSEE
	}

	// Créer le fichier d'extraction
	csvExtractedFile, err := os.Create(ExtractFilePath + CompleteExtractFileName + ".csv")
	if err != nil {
		fmt.Printf("Erreur survenue lors de l'ouverture du fichier Extraction : %v\n", err)
		return
	}
	// On s'assur de sa fermeture a la fin de la fonction
	defer csvExtractedFile.Close()

	// Créer un writer pour écrire dans le fichier d'extraction
	w := csv.NewWriter(csvExtractedFile)
	// Définir le séparateur du fichier d'extraction
	w.Comma = ';'

	// On s'assure que ce qui doit etre écrit est ecris a la fin de la fonction
	defer w.Flush()

	// Écrire le header dans le fichier d'extraction
	header := []string{"x", "y", "imb_id", "num_voie", "cp_no_voie", "type_voie", "nom_voie", "batiment", "code_insee", "code_poste", "nom_com", "catg_loc_imb", "imb_etat", "pm_ref", "pm_etat", "code_l331", "geom_mod", "type_imb"}
	w.Write(header)

	var codeDept string

	// Lire et traiter les enregistrements du fichier source
	for {
		ComptTotal++
		// Lire une ligne du fichier source
		record, err := r.Read()

		// Vérifier si on a atteint la fin du fichier
		if err == io.EOF {
			break
		}

		// Gestion des erreurs de lecture
		if err != nil {
			fmt.Printf("Erreur lors de la lecture d'une ligne : %v\n", err)
			continue
		}

		// Vérifier si la ligne contient suffisamment de colonnes (au moins 10)
		if len(record) >= 10 {
			// Condition pour récupérer les 2 premier caractère (du code postal)
			if len(record[9]) >= 2 {
				codeDept = strings.TrimSpace(record[9][:2])
				for _, dept := range DepartList {
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

			codeInsee := strings.TrimSpace(record[8]) // 9e colonne
			for _, insee := range InseeList {
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
