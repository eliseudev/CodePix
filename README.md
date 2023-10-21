# CodePix-Center-Bank
Para testar as implementações siga os passos abaixo no terminal.
Gerando a imagem no docker: docker compose up -d
Abrir o container com o bash: docker exec -it <contaierName/id> bash
Inicializando o gRPC and Kafka: go run main.go all
Url Kafka local: http://localhost:9021/
Testar Kafka no bash: kafka-topics --list --bootstrap-server=localhost:9092
