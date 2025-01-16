package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {
	// Ouvrir l'image d'entrée
	file, err := os.Open("input.jpg") // Remplacez par le chemin de votre image
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de l'image: %v", err)
	}
	defer file.Close()

	// Décoder l'image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Erreur lors du décodage de l'image: %v", err)
	}

	// Se connecter au serveur
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Erreur de connexion au serveur: %v", err)
	}
	defer conn.Close()

	goRoutineNumberStr := os.Args[1]
	goRoutineNumber, err := strconv.Atoi(goRoutineNumberStr)
	if err != nil {
		log.Fatalf("Erreur lors du passage en int de l'argument%v", err)
	}

	// Envoyer l'image au serveur
	err = sendGoRoutineNumber(conn, goRoutineNumber)
	if err != nil {
		log.Fatalf("Erreur lors de l'envoi de l'image: %v", err)
	}

	// Envoyer l'image au serveur
	err = sendImage(conn, img)
	if err != nil {
		log.Fatalf("Erreur lors de l'envoi de l'image: %v", err)
	}

	// Recevoir l'image floutée
	blurredImg, err := receiveImage(conn)
	if err != nil {
		log.Fatalf("Erreur lors de la réception de l'image: %v", err)
	}

	// Sauvegarder l'image floutée sur le disque
	outFile, err := os.Create("output.jpg")
	if err != nil {
		log.Fatalf("Erreur lors de la création du fichier de sortie: %v", err)
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, blurredImg, nil)
	if err != nil {
		log.Fatalf("Erreur lors de l'encodage de l'image de sortie: %v", err)
	}

	fmt.Println("Image floutée sauvegardée sous 'output.jpg'")
}

// Fonction pour envoyer l'image au serveur
func sendImage(conn net.Conn, img image.Image) error {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	if err != nil {
		return err
	}

	// Envoyer l'image via la connexion
	_, err = conn.Write(buf.Bytes())
	return err
}

// Fonction pour envoyer le nombre de GoRoutine
func sendGoRoutineNumber(conn net.Conn, val int) error {
	err := binary.Write(conn, binary.BigEndian, int32(val)) // Utiliser int32 pour garantir la taille de l'entier
	return err
}

// Fonction pour recevoir l'image du serveur
func receiveImage(conn net.Conn) (image.Image, error) {
	img, _, err := image.Decode(conn)
	if err != nil {
		return nil, fmt.Errorf("échec du décodage de l'image : %v", err)
	}
	return img, nil
}
