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
	fmt.Println(wi.Config.File_export)

	// Faire le traitement de lecture du fichier IMMEUBLE
	// Lire les parametres de config depuis wi.Config

	// Exporter les infos vers le CSV

	return nil
}
