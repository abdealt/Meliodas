package components

import (
	"fmt"
)

func GetStats() {
	// Récupérations des variables de comptage du fichier ReadCSVFileContentAndExtracter
	fmt.Printf("Il y a %v éléments totaux dans le fichier source. Sur tous ces éléments, il y a %v éléments exportés. \n----------------------------------------------------------------------------------------------------------", ComptTotal, ComptElement)
}
