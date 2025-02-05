# Sieve-of-Eratosthenes-Server

```markdown
# Repositório para Testes de RTT (Round-Trip Time) com TCP e UDP

Este repositório contém códigos em Go para servidores e clientes projetados para medir o Round-Trip Time (RTT) de comunicação através de TCP e UDP. O objetivo é avaliar o desempenho da rede sob diferentes protocolos.

## Estrutura do Repositório

O repositório está organizado da seguinte forma:
```

├── README.md (Este arquivo)
├── crivo/ (Códigos dos algoritmos)
│ ├── crivo_concorrente.go (Servidor TCP)
│ ├── crivo_nao_concorrente.go (Servidor UDP)
│ └── results/ (Resultados com os valores primos)
│ └── resultado_com_primos.json
├── servers/ (Códigos dos servidores)
│ ├── tcp_server.go (Servidor TCP)
│ ├── udp_server.go (Servidor UDP)
│ └── udp_server_ack.go (Servidor UDP com ACK)
└── clients/ (Códigos dos clientes)
├── tcp_client.go (Cliente TCP)
├── udp_client.go (Cliente UDP)
├── udp_client_ack.go (Cliente UDP com ACK e tratamento de perdas)
└── results/ (Resultados)
├── resultado_rtt_tcp.json
├── resultado_rtt_udp.json
└── resultado_rtt_udp_ack.json

````

## Descrição dos Códigos

### Servidores (`servers/`)

*   **`tcp_server.go`**:  Implementa um servidor TCP que escuta na porta 8082. Recebe um inteiro `n` do cliente via JSON, calcula os números primos até `n` usando o Crivo de Eratóstenes, e envia a lista de primos de volta para o cliente em formato JSON.  Utiliza goroutines para lidar com múltiplas conexões simultaneamente.

*   **`udp_server.go`**:  Implementa um servidor UDP que escuta na porta 8080. Recebe um inteiro `n` do cliente via JSON, calcula os números primos até `n` usando o Crivo de Eratóstenes, e envia a lista de primos de volta para o cliente em formato JSON. Usa um semáforo para limitar o número de goroutines concorrentes.

*   **`udp_server_ack.go`**: Implementa um servidor UDP que escuta na porta 8081. Semelhante ao `udp_server.go`, mas envia uma mensagem de "ack" (acknowledgment) para o cliente após enviar os resultados, indicando que a mensagem foi recebida.

### Clientes (`clients/`)

*   **`tcp_client.go`**:  Implementa um cliente TCP que se conecta ao `tcp_server.go`. Envia um inteiro `n` para o servidor, recebe a lista de números primos em formato JSON e calcula o RTT para cada invocação. Ao final, calcula o RTT médio e salva os RTTs individuais em um arquivo JSON (`resultado_rtt_tcp.json`).

*   **`udp_client.go`**: Implementa um cliente UDP que se comunica com `udp_server.go`. Envia um inteiro `n` para o servidor, recebe a lista de números primos em formato JSON e calcula o RTT.  Calcula o RTT médio e salva os RTTs individuais em um arquivo JSON (`resultado_rtt_udp.json`).

*   **`udp_client_ack.go`**: Implementa um cliente UDP que se comunica com `udp_server_ack.go`.  Envia um inteiro `n` ao servidor e espera receber os resultados e um "ACK". Utiliza um timeout e reenvia os dados caso o ACK não seja recebido dentro do tempo limite, simulando um tratamento básico de perdas de pacotes.  Calcula o RTT médio e salva os RTTs individuais em um arquivo JSON (`resultado_rtt_udp_ack.json`).

## Como Executar

### Pré-requisitos

*   Go instalado (versão 1.16 ou superior)

### Executando os Servidores

1.  Navegue até a pasta `servers/`:

    ```bash
    cd servers/
    ```

2.  Execute o servidor TCP:

    ```bash
    go run tcp_server.go
    ```

3.  Execute o servidor UDP:

    ```bash
    go run udp_server.go
    ```

4.  Execute o servidor UDP com ACK:

    ```bash
    go run udp_server_ack.go
    ```

    **Importante:** Execute cada servidor em uma janela de terminal separada.

### Executando os Clientes

1.  Navegue até a pasta `clients/`:

    ```bash
    cd clients/
    ```

2.  Execute o cliente TCP:

    ```bash
    go run tcp_client.go
    ```

3.  Execute o cliente UDP:

    ```bash
    go run udp_client.go
    ```

4.  Execute o cliente UDP com ACK:

    ```bash
    go run udp_client_ack.go
    ```

    **Importante:** Certifique-se de que o servidor correspondente esteja em execução antes de executar o cliente.

### Observações

*   Os clientes são configurados para enviar `numInvocacoes` (definido como 10000) requisições ao servidor.
*   Os clientes salvam os resultados do RTT (RTT médio e RTTs individuais) em arquivos JSON na pasta `clients/`.
*   A constante `n` (tamanho máximo para calcular os primos) está definida como 100000 nos clientes.  Aumentar esse valor pode aumentar o tempo de processamento, mas pode ser necessário para obter resultados mais representativos.
*   Os servidores utilizam um semáforo para controlar o número de goroutines concorrentes (`maxGoroutines` = 1000).  Isso ajuda a prevenir o consumo excessivo de recursos do sistema, especialmente ao lidar com um grande número de requisições simultâneas.
*   O `udp_client_ack.go` implementa uma forma básica de tratamento de perdas. Em um ambiente real, mecanismos mais robustos de controle de congestionamento e retransmissão seriam necessários.
*   O tamanho do buffer de leitura do UDP é grande (65507) para acomodar mensagens potencialmente grandes retornadas pelo servidor. Ajuste conforme a necessidade.

## Lógica dos Códigos

### Crivo de Eratóstenes

Tanto os servidores quanto os clientes utilizam a função `crivo(n)` para calcular os números primos até `n`.  Este método implementa o algoritmo do Crivo de Eratóstenes, que é um algoritmo eficiente para encontrar todos os números primos até um determinado limite.

### JSON para Serialização

A comunicação entre cliente e servidor (tanto TCP quanto UDP) utiliza JSON (JavaScript Object Notation) para serializar e desserializar os dados. A biblioteca `encoding/json` do Go é usada para codificar os dados a serem enviados e decodificar os dados recebidos.

### RTT

O RTT (Round-Trip Time) é calculado no cliente medindo o tempo entre o envio da requisição para o servidor e o recebimento da resposta.  Essa medida fornece uma estimativa do tempo necessário para um pacote viajar do cliente para o servidor e de volta.

### ACK (Acknowledgment)

O servidor UDP com ACK (`udp_server_ack.go`) envia uma mensagem de confirmação (ACK) ao cliente após enviar os resultados. Isso permite que o cliente verifique se a mensagem foi entregue com sucesso. O cliente (`udp_client_ack.go`) espera por essa mensagem e, se não a receber dentro de um determinado período, reenvia a requisição.

Este repositório fornece um ponto de partida para avaliar e comparar o desempenho de TCP e UDP sob diferentes condições de rede. Os códigos podem ser modificados e estendidos para incluir funcionalidades adicionais, como tratamento de erros mais robusto, controle de congestionamento e simulação de perdas de pacotes.
````
