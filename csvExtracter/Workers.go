package csvExtracter

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Config va contenir les paramètres du programme
type Config struct {
	File_immeuble string
	File_export   string
	File_log      string
	Lst_Dprt      []string
}

// Retorune un pointeur sur config
type WorkerImmeuble struct {
	Config *Config
}

// NewWorkerImmeuble crée un nouvel objet WorkerImmeuble avec les paramètres fournis
func NewWorkerImmeuble(cfg Config) (*WorkerImmeuble, error) {
	workerImmeuble := &WorkerImmeuble{
		Config: &cfg,
	}
	return workerImmeuble, nil
}

// Créations des variables pour compter
var ComptElement int
var ComptTotal int

// SynoReaderRunner est la fonction d'execution des traitements autour du synoptique
func (wi *WorkerImmeuble) SuperreaderCSV() error {
	now := time.Now()

	// Ouverture du fichier log avec les permission d'ecriture
	logFile, err := os.OpenFile(wi.Config.File_log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier log : %v", err)
	}
	defer logFile.Close()
	logFile.WriteString("Le fichier log a été ouvert\n")

	// Ouverture du fichier source
	sourceFile, err := os.Open(wi.Config.File_immeuble)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier source : %v", err)
	}
	defer sourceFile.Close()
	logFile.WriteString(fmt.Sprintf("Le fichier source a été ouvert : %v\n", wi.Config.File_immeuble))

	// Création de l'instance de lecture et son séparateur
	readerInstance := csv.NewReader(sourceFile)
	readerInstance.Comma = ','

	// Ouverture du fichier export
	ExportedFile, err := os.Create(wi.Config.File_export + "Export_du_" + now.Format("02-01-2006_15-04-05") + ".csv")
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier export : %v", err)
	}
	defer ExportedFile.Close()
	logFile.WriteString(fmt.Sprintf("Le fichier d'export a été ouvert : %v\n", wi.Config.File_export))

	// Création de l'instance de l'écriture et son séparateur
	writerInstance := csv.NewWriter(ExportedFile)
	writerInstance.Comma = ';'
	defer writerInstance.Flush()

	// Écriture du header
	header := []string{"x", "y", "imb_id", "num_voie", "cp_no_voie", "type_voie", "nom_voie", "batiment", "code_insee", "code_poste", "nom_com", "type_imb"}
	writerInstance.Write(header)

	logFile.WriteString("Début de la lecture des enregistrements\n")
	// Lire et écrire les lignes
	// Lire et écrire les lignes
	for {
		// Lire une ligne du CSV
		record, err := readerInstance.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("erreur lors de la lecture de la ligne : %w", err)
		}
		ComptTotal++

		// Filtrer les lignes selon les codes département
		if len(record) >= 10 && len(record[8]) >= 2 {
			codeDpt := strings.TrimSpace(record[8][:2])

			// Vérifier les éléments dans Lst_Dprt
			for _, dept := range wi.Config.Lst_Dprt {
				if codeDpt == dept {
					ComptElement++

					// Vérifie que la ligne contient au moins six colonnes
					if len(record) >= 6 {
						// Crée un tableau contenant les 11 premières colonnes et la dernière
						selectedColumns := append(record[:11], record[len(record)-1])
						writerInstance.Write(selectedColumns)
					} else {
						// Gérer le cas où la ligne ne contient pas assez de colonnes
						writerInstance.Write(record)
					}
					break
				}
			}
		}
	}
	lap := time.Since(now)
	logFile.WriteString(fmt.Sprintf("Extraction finie | Nombre total d'enregistrement trouvés : %v | Nombre d'enregistrements Exporté : %v\n---FIN---\n\n", ComptTotal, ComptElement))

	fmt.Printf("Extraction terminée, le résultat est disponible ici (vous pouvez copier coller)\n: %s\nTemps de l'opération : %v\n", wi.Config.File_export+"Export_du_"+now.Format("02-01-2006_15-04-05")+".csv", lap)
	return nil
}

func (wi *WorkerImmeuble) ExtractStatisticsFromCSV() error {
	// Ouverture du fichier source
	sourceFile, err := os.Open(wi.Config.File_immeuble)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier source : %v", err)
	}
	defer sourceFile.Close()

	// Création de l'instance de lecture et son séparateur
	readerInstance := csv.NewReader(sourceFile)
	readerInstance.Comma = ','

	// Lire et écrire les lignes
	for {
		// Lire une ligne du CSV
		record, err := readerInstance.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("erreur lors de la lecture de la ligne : %v", err)
		}
		// Compteur d'élèments totaux
		ComptTotal++

		// Filtrer les lignes selon les codes département
		if len(record) >= 10 && len(record[8]) >= 2 {
			codeDpt := strings.TrimSpace(record[8][:2])

			// Vérifier les éléments dans Lst_Dprt
			for _, dept := range wi.Config.Lst_Dprt {
				if codeDpt == dept {
					ComptElement++
					break
				}
			}
		}
	}
	str := "%"
	prc := ComptElement * 100 / ComptTotal
	fmt.Printf("Il ya au total %v elements, dont %v traité. Soit un pourcentage de %v%s.", ComptTotal, ComptElement, prc, str)
	return nil
}

func (wi *WorkerImmeuble) ExtractDepartFromCSV() error {
	// Ouverture du fichier source
	sourceFile, err := os.Open(wi.Config.File_immeuble)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier source : %v", err)
	}
	defer sourceFile.Close()

	// Initialisation du lecteur
	readerInstance := csv.NewReader(sourceFile)
	readerInstance.Comma = ','

	// Utilisation d'une map pour éviter les doublons clé chaine de caractère et valeurs structure vide
	listeDepartements := make(map[string]struct{})

	// Lecture de chaque ligne du fichier
	for {
		// Lecture d'une ligne
		record, err := readerInstance.Read()

		// Si fin du fichier
		if err == io.EOF {
			break
		}

		// Si erreur lors de la lectuer d'une ligne
		if err != nil {
			fmt.Printf("Erreur lors de la lecture de la ligne : %v\n", err)
			continue
		}

		// Initialisation de la colonne département
		if len(record) >= 18 {
			if len(record[8]) >= 2 {
				codeDept := strings.TrimSpace(record[8][:2])
				// Ajout dans le map
				listeDepartements[codeDept] = struct{}{}
			} else {
				// Gérer le cas où la chaîne est trop courte
				continue
			}

		} else {
			continue
		}
	}

	// Affichage des départements
	fmt.Printf("Liste des départements :\n")
	for dpt := range listeDepartements {
		fmt.Print(dpt + "|")
	}
	return nil
}
