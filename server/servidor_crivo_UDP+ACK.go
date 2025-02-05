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

type Message struct {
	Type string      `json:"type"` // Tipo da mensagem: "result" ou "ack"
	Data interface{} `json:"data"`
}

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

	resultMsg := Message{
		Type: "result",
		Data: primos,
	}
	resultJSON, err := json.Marshal(resultMsg)
	if err != nil {
		fmt.Println("Erro ao gerar JSON do resultado:", err)
		return
	}
	_, err = conn.WriteToUDP(resultJSON, addr)
	if err != nil {
		fmt.Println("Erro ao enviar resultado para o cliente:", err)
		return
	}

	ackMsg := Message{
		Type: "ack",
		Data: "ACK",
	}
	ackJSON, err := json.Marshal(ackMsg)
	if err != nil {
		fmt.Println("Erro ao gerar JSON do ACK:", err)
		return
	}
	_, err = conn.WriteToUDP(ackJSON, addr)
	if err != nil {
		fmt.Println("Erro ao enviar ACK para o cliente:", err)
		return
	}
}

func main() {
	addr := net.UDPAddr{
		Port: 8081,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor UDP:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Servidor UDP iniciado na porta 8081")

	var wg sync.WaitGroup

	for {
		buffer := make([]byte, 65507)

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
