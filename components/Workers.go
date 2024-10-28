package components

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Config struct {
	File_immeuble string
	File_export   string
	File_log      string
	Lst_Insee     []string
	Lst_Dprt      []string
}

type WorkerImmeuble struct {
	Config *Config
}

func NewWorkerImmeuble(cfg Config) (*WorkerImmeuble, error) {
	workerImmeuble := &WorkerImmeuble{
		Config: &cfg,
	}
	return workerImmeuble, nil
}

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

	// Exporter les infos vers le CSV

	// Création d'un nom unique pour le fichier exporté

	// Ouverture du fichier export
	fileE, err := os.Create(wi.Config.File_export + "Export_du_" + now.Format("Mon Jan 2 15:04:05") + ".csv")
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier export : %w", err)
	}
	defer fileE.Close()

	// Création de l'instance de l'ecriture et son séparateur
	w := csv.NewWriter(fileE)
	w.Comma = ';'
	defer w.Flush()

	// Ecriture du header
	header := []string{"x", "y", "imb_id", "num_voie", "cp_no_voie", "type_voie", "nom_voie", "batiment", "code_insee", "code_poste", "nom_com", "catg_loc_imb", "imb_etat", "pm_ref", "pm_etat", "code_l331", "geom_mod", "type_imb"}
	w.Write(header)

	// Lire et écrire les lignes
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("erreur lors de la lecture de la ligne : %w", err)
		}
		// Filtrer les lignes selon les codes département et INSEE
		if len(record) >= 10 {
			codeInsee := record[8]
			codeDpt := strings.TrimSpace(record[9][:2])

			for _, insee := range wi.Config.Lst_Insee {
				if codeInsee == insee {
					w.Write(record)
					break

				} else {
					continue
				}
			}
			for _, dept := range wi.Config.Lst_Dprt {
				if codeDpt == dept {
					w.Write(record)
					break
				} else {
					continue
				}
			}
		}
	}
	fmt.Printf("Extraction terminée, le résultat est disponible dans le fichier : %s\n", wi.Config.File_export+".csv")
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
		CptTo++
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("erreur lors de la lecture de la ligne : %w", err)
		}
		// Filtrer les lignes selon les codes département et INSEE
		if len(record) >= 10 {
			codeInsee := record[8]                      // Colonne du code insee
			codeDpt := strings.TrimSpace(record[8][:2]) // Les 2 premier caractère du code insee sont le département

			for _, insee := range wi.Config.Lst_Insee {
				if codeInsee == insee {
					CptEl++
					break

				} else {
					continue
				}
			}
			for _, dept := range wi.Config.Lst_Dprt {
				if codeDpt == dept {
					CptEl++
					break
				} else {
					continue
				}
			}
		}
	}
	fmt.Printf("Il y'a %v éléments totaux. Il y'a %v éléments traités.", CptTo, CptEl)
	return nil
}
