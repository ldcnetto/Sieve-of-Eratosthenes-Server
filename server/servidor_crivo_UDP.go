package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"sync"
)

const maxGoroutines = 1000 // Limite máximo de goroutines simultâneas

var sem = make(chan struct{}, maxGoroutines) // Semáforo para limitar goroutines

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

func handleUDPConnection(conn *net.UDPConn, addr *net.UDPAddr, buffer []byte) {
	defer func() { <-sem }() // Libera uma posição no semáforo quando a goroutine terminar

	var num int
	err := json.Unmarshal(buffer, &num)
	if err != nil {
		fmt.Println("Erro ao decodificar o valor de n:", err)
		return
	}

	primos := crivo(num)

	jsonData, err := json.Marshal(primos)
	if err != nil {
		fmt.Println("Erro ao gerar JSON:", err)
		return
	}

	_, err = conn.WriteToUDP(jsonData, addr)
	if err != nil {
		fmt.Println("Erro ao enviar dados para o cliente:", err)
		return
	}
}

func main() {
	addr := net.UDPAddr{
		Port: 8080,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor UDP:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Servidor UDP iniciado na porta 8080")

	var wg sync.WaitGroup

	for {
		buffer := make([]byte, 1024)

		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Erro ao ler dados do cliente:", err)
			continue
		}

		sem <- struct{}{} // Ocupa uma posição no semáforo
		wg.Add(1)

		go func() {
			defer wg.Done()
			handleUDPConnection(conn, addr, buffer[:n])
		}()
	}
}
