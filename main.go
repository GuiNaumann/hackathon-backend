package main

import (
	setup "hackathon-backend/infrastructure"
	"hackathon-backend/settings_loader"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func init() {
	dir, _ := os.Getwd()
	log.Printf("Diretório atual: %s", dir)

	// Carrega as variáveis do .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: arquivo .env não encontrado")
	}
}

func initLogger() {
	logDir := "./logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			log.Fatalf("Erro ao criar a pasta de logs: %v", err)
		}
	}

	logFile, err := os.OpenFile(logDir+"/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Erro ao abrir/criar arquivo de log: %v", err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	router := mux.NewRouter()
	initLogger()
	log.Println("Servidor iniciado com sucesso!")

	// Configurar o middleware CORS
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://localhost:8081",
			"http://localhost:5173",
		},
		AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	// Carregar as configurações
	settings := settings_loader.NewSettingsLoader()

	// Configura o projeto chamando o Setup da infraestrutura
	setupConfig, err := setup.Setup(router, settings)
	if err != nil {
		log.Fatalf("Erro ao configurar a infraestrutura: %v", err)
	}
	defer setupConfig.CloseDB()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Servidor HTTP em goroutine separada
	go func() {
		port := settings.GetAppPort()
		log.Printf("🚀 Servidor rodando na porta %s", port)
		if err := http.ListenAndServe(":"+port, handler); err != nil {
			log.Fatal(err)
		}
	}()

	// Aguardar sinal de parada
	<-sigChan
	log.Println("✅ Servidor encerrado")
}
