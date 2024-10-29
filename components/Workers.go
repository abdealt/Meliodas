package components

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
	Lst_Insee     []string
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
var CptEl int
var CptTo int

// SynoReaderRunner est la fonction d'execution des traitements autour du synoptique
func (wi *WorkerImmeuble) SuperreaderCSV() error {
	now := time.Now()

	// Ouverture du fichier source
	fileS, err := os.Open(wi.Config.File_immeuble)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier source : %w", err)
	}
	defer fileS.Close()

	// Création de l'instance de lecture et son séparateur
	r := csv.NewReader(fileS)
	r.Comma = ','

	// Ouverture du fichier export
	fileE, err := os.Create(wi.Config.File_export + "Export_du_" + now.Format("02-01-2006_15-04-05") + ".csv")
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier export : %w", err)
	}
	defer fileE.Close()

	// Création de l'instance de l'écriture et son séparateur
	w := csv.NewWriter(fileE)
	w.Comma = ';'
	defer w.Flush()

	// Écriture du header
	header := []string{"x", "y", "imb_id", "num_voie", "cp_no_voie", "type_voie", "nom_voie", "batiment", "code_insee", "code_poste", "nom_com", "catg_loc_imb", "imb_etat", "pm_ref", "pm_etat", "code_l331", "geom_mod", "type_imb"}
	w.Write(header)

	// Lire et écrire les lignes
	for {
		// Lire une ligne du CSV
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("erreur lors de la lecture de la ligne : %w", err)
		}

		// Filtrer les lignes selon les codes département et INSEE
		if len(record) >= 10 && len(record[9]) >= 2 {
			codeInsee := record[8]
			codeDpt := strings.TrimSpace(record[9][:2])

			// Vérifier les éléments dans Lst_Insee
			for _, insee := range wi.Config.Lst_Insee {
				if codeInsee == insee {
					w.Write(record)
					break
				}
			}

			// Vérifier les éléments dans Lst_Dprt
			for _, dept := range wi.Config.Lst_Dprt {
				if codeDpt == dept {
					w.Write(record)
					break
				}
			}
		}
	}

	fmt.Printf("Extraction terminée, le résultat est disponible ici : %s\n", wi.Config.File_export+"Export_du_"+now.Format("02-01-2006_15-04-05")+".csv")
	return nil
}

func (wi *WorkerImmeuble) ExtractStatisticsFromCSV() error {
	// Ouverture du fichier source
	fileS, err := os.Open(wi.Config.File_immeuble)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier source : %w", err)
	}
	defer fileS.Close()

	// Création de l'instance de lecture et son séparateur
	r := csv.NewReader(fileS)
	r.Comma = ','

	// Lire et écrire les lignes
	for {
		// Lire une ligne du CSV
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("erreur lors de la lecture de la ligne : %w", err)
		}
		// Compteur d'élèments totaux
		CptTo++

		// Filtrer les lignes selon les codes département et INSEE
		if len(record) >= 10 && len(record[9]) >= 2 {
			codeInsee := record[8]
			codeDpt := strings.TrimSpace(record[9][:2])

			// Vérifier les éléments dans Lst_Insee
			for _, insee := range wi.Config.Lst_Insee {
				if codeInsee == insee {
					CptEl++
					break
				}
			}

			// Vérifier les éléments dans Lst_Dprt
			for _, dept := range wi.Config.Lst_Dprt {
				if codeDpt == dept {
					CptEl++
					break
				}
			}
		}
	}
	str := "%"
	prc := CptEl * 100 / CptTo
	fmt.Printf("Il ya au total %v elements, dont %v traité. Soit un pourcentage de %v%s.", CptTo, CptEl, prc, str)
	return nil
}

func (wi *WorkerImmeuble) LogWriteInfo() error {
	now := time.Now()
	// Ouverture du fichier Log
	file, err := os.Open(wi.Config.File_log)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier log : %w", err)
	}
	defer file.Close()

	// Création du méssage
	message := fmt.Sprintf("Une extraction a été effectuée le : %s | depuis le fichier source %s | vers nouveau fichier %v.\n", now.Format("2006-01-02 15:04:05"), wi.Config.File_immeuble, wi.Config.File_export+"Export_du_"+now.Format("Mon Jan 2 15:04:05")+".csv")
	message += fmt.Sprintf("Les filtres actifs sont INSEE : %v et DPT :%v. Il un total de %v éléments, et %v qui sont extraits.", wi.Config.Lst_Insee, wi.Config.Lst_Dprt, CptTo, CptEl)

	// Ecriture du message dans le log
	file.WriteString(message)

	return nil
}

func (wi *WorkerImmeuble) ExtractDepartFromCSV() error {
	// Ouverture du fichier source
	fileS, err := os.Open(wi.Config.File_immeuble)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier source : %w", err)
	}
	defer fileS.Close()

	// Initialisation du lecteur
	r := csv.NewReader(fileS)
	r.Comma = ','

	// Utilisation d'une map pour éviter les doublons clé chaine de caractère et valeurs structure vide
	listeDepartements := make(map[string]struct{})

	// Lecture de chaque ligne du fichier
	for {
		// Lecture d'une ligne
		record, err := r.Read()

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
			if len(record[9]) >= 2 {
				codeDept := strings.TrimSpace(record[9][:2])
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
