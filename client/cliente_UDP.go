package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	const n = 100000
	const serverAddr = "localhost:8080"
	const numInvocacoes = 10000

	var totalRTT time.Duration
	var rtts []float64

	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("Erro ao resolver endereço UDP:", err)
		return
	}

	for i := 0; i < numInvocacoes; i++ {
		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			fmt.Println("Erro ao conectar ao servidor:", err)
			return
		}

		start := time.Now()

		jsonData, err := json.Marshal(n)
		if err != nil {
			fmt.Println("Erro ao gerar JSON:", err)
			conn.Close()
			return
		}

		_, err = conn.Write(jsonData)
		if err != nil {
			fmt.Println("Erro ao enviar dados para o servidor:", err)
			conn.Close()
			return
		}

		buffer := make([]byte, 65507)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Erro ao receber dados do servidor:", err)
			conn.Close()
			return
		}

		var primos []int
		err = json.Unmarshal(buffer[:n], &primos)
		if err != nil {
			fmt.Println("Erro ao decodificar dados do servidor:", err)
			conn.Close()
			return
		}

		rtt := time.Since(start)
		totalRTT += rtt
		rtts = append(rtts, rtt.Seconds()*1000)

		conn.Close()
	}

	avgRTT := totalRTT.Seconds() * 1000 / float64(numInvocacoes)

	data := map[string]interface{}{
		"rtt_medio_milissegundo": avgRTT,
		"rtts_individuais":       rtts,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Erro ao gerar JSON:", err)
		return
	}

	filePath := "results/resultado_rtt_udp.json"

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Erro ao criar arquivo:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Erro ao escrever no arquivo:", err)
		return
	}

	fmt.Println("Dados de RTT salvos em: ", filePath)
}
