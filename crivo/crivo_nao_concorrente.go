package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"
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

func main() {
	const n = 1000000

	// Mede o tempo de execução
	start := time.Now()
	primos := crivo(n)
	duration := time.Since(start).Seconds()

	// Cria o mapa para armazenar os dados
	data := map[string]interface{}{
		"tempo_execucao":    duration,
		"quantidade_primos": len(primos),
		"limite_superior":   n,
		"primos":            primos, // Adiciona a lista de primos ao JSON
	}

	// Converte o mapa para JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Erro ao gerar JSON:", err)
		return
	}

	filePath := "results/resultado_com_primos.json"

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
