package main

import (
	"github.com/abdealt/meliodas/components"
)

func main() {
	// Appel de la fonction exportée ReadFileContent
	components.ReadCSVFileContent()
	// Appel de la fonction de création du log
	components.LogWriter()
}
