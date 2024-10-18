package components

import (
	"fmt"
	"os"
)

func GetStats() {
	// Ajout des informations aux fichier log
	logFilePath := os.Getenv("LOG_FILE")
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Une erreur s'est produite lors de l'ouverture du fichier log : %v", err)
		return
	}
	defer logFile.Close()

	// Ecriture des statistique dans le Log
	logFile.WriteString(fmt.Sprintf("Il y a %v élements totaux dans le fichier source, sur tous ces éléments, il y'a %v éléments exportés", ComptTotal, ComptElement))
}
