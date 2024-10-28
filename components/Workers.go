package components

import "fmt"

type Config struct {
	File_immeuble string
	File_export   string
	File_log      string
	Lst_Dprt      []string
	Lst_Insee     []string
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

// SynoReaderRunner est la fonction d'execution des traitements autour du synoptique
func (wi *WorkerImmeuble) SuperreaderCSV() error {

	// Lire les parametres de config depuis wi.Config
	if wi.Config.File_immeuble == "" {
		return fmt.Errorf("aucun chemin de fichier source n'a été donné")
	}

	if wi.Config.File_export == "" {
		return fmt.Errorf("aucun chemin pour le fichier d'extraction n'a été donné")
	}

	if wi.Config.File_log == "" {
		return fmt.Errorf("aucun chemin de fichier log n'a été donné")

	}

	if len(wi.Config.Lst_Dprt) == 0 && len(wi.Config.Lst_Insee) == 0 {
		return fmt.Errorf("la liste d'extraction est vide")
	}

	// Faire le traitement de lecture du fichier IMMEUBLE

	// Exporter les infos vers le CSV

	return nil
}
