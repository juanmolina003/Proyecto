package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// --- Función para enviar la alerta ---
func enviarAlertaDiscord(mensaje string) {
	webhookURL := "https://discordapp.com/api/webhooks/1454198381495976271/dUKGCz5e4Wmj4NekTOjrgaPLka7Uq6MuukAr5-Lbvm9syiEKZjFZ26y9_QBJqbHee4H5"

	payload := map[string]string{
		"content": "🚨 **ALERTA DEVSECOPS**: " + mensaje,
	}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("Error enviando alerta: %v", err)
		return
	}
	defer resp.Body.Close()
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/simular-fallo" {
		// Log crítico
		log.Printf("level=critical msg='Evento anómalo detectado' method=%s path=%s", r.Method, r.URL.Path)

		// Lanzamos la alerta al detectar el fallo
		enviarAlertaDiscord("Se ha detectado un acceso al endpoint de fallo desde " + r.RemoteAddr)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error 500: Fallo detectado y alerta enviada.")
		return
	}

	if r.URL.Path == "/" {
		log.Printf("level=info method=%s path=%s", r.Method, r.URL.Path)
		http.ServeFile(w, r, "static/logo.png")
		return
	}

	log.Printf("level=warning msg='Ruta no encontrada' path=%s", r.URL.Path)
	http.NotFound(w, r)
}

// ... (resto de tu código arriba)

func main() {
	// 1. Configuración del handler
	http.HandleFunc("/", handler)
	http.Handle("/metrics", promhttp.Handler())

	// 2. CORRECCIÓN G302: Cambiamos 0644 por 0600 para mayor seguridad
	f, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Error al abrir archivo de logs: %v\n", err)
		return
	}
	defer f.Close()
	log.SetOutput(f)

	// 3. Configuración SEGURA del servidor
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	fmt.Println("Servidor iniciado con sistema de alertas activo en http://localhost:8080")
	log.Println("level=info msg='Iniciando servidor seguro'")

	// 4. CORRECCIÓN G104: Controlamos el error de ListenAndServe
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error crítico al arrancar el servidor: %v", err)
	}
}
