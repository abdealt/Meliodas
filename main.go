package main

import (
	"github.com/abdealt/meliodas/components"
)

func main() {
	// Appel de la fonction exportée ReadFileContent
	components.ReadCSVFileContentAndExtracter()
	// Appel de la fonction de création du log
	components.LogWriter()
	// Appel de la fonction statistique
	components.GetStats()
}
