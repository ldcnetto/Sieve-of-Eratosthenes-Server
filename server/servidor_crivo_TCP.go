package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
)

func crivo(n int) []int {
	numeros := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		numeros[i] = true
	}

	limite := int(math.Sqrt(float64(n)))

	for i := 2; i <= limite; i++ {
		if numeros[i] {
			for j := i * i; j <= n; j += i {
				numeros[j] = false
			}
		}
	}

	var primos []int
	for i := 2; i <= n; i++ {
		if numeros[i] {
			primos = append(primos, i)
		}
	}

	return primos
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	var n int
	err := json.NewDecoder(conn).Decode(&n)
	if err != nil {
		fmt.Println("Erro ao decodificar o valor de n:", err)
		return
	}

	primos := crivo(n)

	jsonData, err := json.Marshal(primos)
	if err != nil {
		fmt.Println("Erro ao gerar JSON:", err)
		return
	}

	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Println("Erro ao enviar dados para o cliente:", err)
		return
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor TCP:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor TCP iniciado na porta 8082")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go handleTCPConnection(conn)
	}
}
