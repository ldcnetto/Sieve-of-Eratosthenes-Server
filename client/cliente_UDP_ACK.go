package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

type Message struct {
	Type string      `json:"type"` // Tipo da mensagem: "result" ou "ack"
	Data interface{} `json:"data"`
}

func main() {
	const n = 100000
	const serverAddr = "localhost:8081"
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
		defer conn.Close()

		start := time.Now()

		jsonData, err := json.Marshal(n)
		if err != nil {
			fmt.Println("Erro ao gerar JSON:", err)
			return
		}

		for {
			// Envia os dados para o servidor
			_, err = conn.Write(jsonData)
			if err != nil {
				fmt.Println("Erro ao enviar dados para o servidor:", err)
				return
			}

			// Configura um timeout para a leitura
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))

			// Buffer para receber a resposta
			buffer := make([]byte, 65507)
			n, err := conn.Read(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					fmt.Println("Timeout, reenviando dados...")
					continue
				}
				fmt.Println("Erro ao receber dados do servidor:", err)
				return
			}

			// Decodifica a mensagem recebida
			var msg Message
			err = json.Unmarshal(buffer[:n], &msg)
			if err != nil {
				fmt.Println("Erro ao decodificar mensagem do servidor:", err)
				return
			}

			switch msg.Type {
			case "result":
				continue
			case "ack":
				break
			default:
				fmt.Println("Tipo de mensagem desconhecido:", msg.Type)
				return
			}

			break
		}

		rtt := time.Since(start)
		totalRTT += rtt
		rtts = append(rtts, rtt.Seconds()*1000) // Armazena o RTT em milissegundos com casas decimais
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

	filePath := "results/resultado_rtt_udp_ack.json"

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
