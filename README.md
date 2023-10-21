# CodePix-Center-Bank
Para testar as implementações siga os passos abaixo no terminal.</br>
Gerando a imagem no docker: docker compose up -d</br>
Abrir o container com o bash: docker exec -it <contaierName/id> bash</br>
Inicializando o gRPC and Kafka: go run main.go all</br>
Url Kafka local: http://localhost:9021/</br>
Testar Kafka no bash: kafka-topics --list --bootstrap-server=localhost:9092</br>
