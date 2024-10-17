package components

import (
	"fmt"
	"os"
	"time"
)

func logWriter() {
	//Initialisation d'une variable now
	now := time.Now()

	// Ajout des informations aux fichier log
	logFilePath := "/logs/log.txt"
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Une erreur s'est produite lors de l'ouverture du fichier log : %v", err)
		return
	}
	defer logFile.Close()
	logFile.WriteString(fmt.Sprintf("Une extraction a été effectuée le : %s | depuis le fichier source %s | au nouveau fichier %v | les filtre actifs sont INSEE : %v et DPT :%v \n", now.Format("2006-01-02 15:04:05"), components.filePath, components.extractFilePath+components.completeExtractFileName+".csv", components.cityINSEE, components.departID))
}
