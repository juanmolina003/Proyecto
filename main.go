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

func main() {
	// Configuración de archivo de logs
	http.Handle("/metrics", promhttp.Handler())
	f, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	log.SetOutput(f)

	http.HandleFunc("/", handler)

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Servidor iniciado con sistema de alertas activo.")
	server.ListenAndServe()
}
