package main

import (
	"fmt"
	"goTask/client"
	"time"
)

func main() {
	client_ := client.NewClient()
	client_.Run()
	defer client_.Close()
	timeout := time.After(500 * time.Second)

	for {
		select {
		case <-client_.Out:
			return
		case <-timeout:
			fmt.Println("Saliendo del bucle después de 5 segundos")
		default:
			// Continuar con la próxima iteración si no se recibe la señal
		}
	}
}
