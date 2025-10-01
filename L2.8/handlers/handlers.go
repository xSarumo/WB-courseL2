package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"ntp-service/service"
	"time"
)

func TimeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("Method isn`t get")
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	host := r.URL.Query().Get("host")
	if host == "" {
		host = "pool.ntp.org"
	}

	currentTime, err := service.GetCurrentTime(host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(currentTime.Format(time.RFC3339))

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"ntp-host":     host,
			"extract-time": currentTime.Format(time.RFC3339),
			"server_time":  time.Now().Format(time.RFC3339),
		})

}
