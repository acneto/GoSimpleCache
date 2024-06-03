package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	http.HandleFunc("/get", handleGet(connectTCPServer))
	http.HandleFunc("/set", handleSet(connectTCPServer))

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}

func handleGet(connectTCPServer func(string) (string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "Key is required", http.StatusBadRequest)
			return
		}

		response, err := connectTCPServer("GET " + key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}
}

func handleSet(connectTCPServer func(string) (string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		if key == "" || value == "" {
			http.Error(w, "Key and value are required", http.StatusBadRequest)
			return
		}

		_, err := connectTCPServer("SET " + key + " " + value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK\n"))
	}
}

func connectTCPServer(command string) (string, error) {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return "", fmt.Errorf("error connecting: %v", err)
	}

	_, err = conn.Write([]byte(command))
	if err != nil {
		fmt.Println("Error sending message:", err.Error())
		return "", fmt.Errorf("error sending message: %v", err)
	}

	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		fmt.Println("Error reading response:", err.Error())
		return "", fmt.Errorf("error reading response: %v", err)
	}

	return string(response[:n]), nil
}
