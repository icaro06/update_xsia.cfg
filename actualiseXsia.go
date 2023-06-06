/*
Met à jour un numéro d'alarme à partir d'un fichier de configuration xsia.cfg

	A.Villanueva
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*tester si un élément est dans une liste*/
func exists(mot string, liste []string) bool {
	if len(mot) < 1 {
		return false
	}

	parts := strings.Split(mot, ":")
	key := parts[0]

	for _, elemento := range liste {
		if elemento == key {
			return true
		}
	}
	return false
}

func main() {
	var err error
	alms := 0      //Nombre d'alarmes créées
	alm := 0       //premiere alarme ALM
	space := 0     //Espace entre les alarmes et les commandes
	lignes := 0    //Num de ligne
	buffer := "\n" //Buffer
	header := ""   //Header buffer
	header_types := []string{"IP", "PORT", "LOGIN", "PWD"}
	body_types := []string{"ALM_VALUE", "ALM_NUM", "XSIA_TRAME"}

	fmt.Println("Version 1")
	//Analysez s'il y a des arguments
	if len(os.Args) < 4 {
		fmt.Println("Erreur args \n Example : regenConfig {nom_fichier.cfg} {premiere alm} {space} ")
		os.Exit(0)
	}
	alm, err = strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Erreur arg[2] première alarme ", err)
		return
	}
	space, err = strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Erreur arg[3] space alarmes ", err)
		return
	}
	if space < 1 {
		fmt.Println("Erreur arg[3] space < 1 ! ", err)
		return
	}

	// Ouvrir le premier fichier en lecture
	entradaFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier d'entrée : ", err)
		return
	}
	defer entradaFile.Close()

	// Créer le deuxième fichier pour l'écriture
	salidaFile, err := os.Create(os.Args[1] + ".new") //New fichier
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier de sortie :", err)
		return
	}
	defer salidaFile.Close()

	// Lire le contenu du fichier d'entrée ligne par ligne
	scanner := bufio.NewScanner(entradaFile)

	for scanner.Scan() { //Analyser le fichier d'entrée pour écrire chaque ALM
		linea := scanner.Text()
		lignes++ //numéro de la ligne actuelle

		if exists(linea, header_types) { //HEADER types
			header += linea + "\n"
		}

		if exists(linea, body_types) || strings.Contains(linea, "/*") { //Ligne avec ALM_VALUE,XSIA_TRAME OR /*

			if err != nil {
				fmt.Println("Erreur num ligne :", lignes)
				return
			}

			// Écrire les entrées proportionnelles à la valeur de ALM_VALUE dans le fichier de sortie
			if strings.Contains(linea, "ALM_VALUE:") || strings.Contains(linea, "ALM_NUM:") { //Ligne avec ALM_VALUE
				fmt.Println(linea)
				buffer += "ALM_VALUE:" + strconv.Itoa(alm) + "\n"
				alms++
				alm = (alm + space) //alarme actuel
			} else {
				buffer += linea + "\n"
				if strings.Contains(linea, "XSIA_TRAME") { //++ space final
					buffer += "\n"
				}
			}

		}
	}
	//Ecrire Header IP ,PORT,LOGIN , PWD
	fmt.Fprintln(salidaFile, header)

	//Ecrire NBR_ALMs
	fmt.Fprintln(salidaFile, "NBR_ALM:", alms)

	//Ecrire Body
	fmt.Fprintln(salidaFile, buffer)

	if err := scanner.Err(); err != nil {
		fmt.Println("Erreur lors de la lecture du fichier d'entrée:", err)
		return
	}
	fmt.Println("Fichier de sortie créé avec succès. nombre d'alarmes ", alms)
}
