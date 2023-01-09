package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func StaticHandlerJSON(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Working"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

type PlayerInfo struct {
	Species     string
	Description string
}

func GetPlayerInfo(w http.ResponseWriter, r *http.Request) {
	playerId := chi.URLParam(r, "playerId")

	url := fmt.Sprintf("https://api.brawlstars.com/v1/players/%%23%s", playerId)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	apiKey := goDotEnvVariable("API_KEY")
	req.Header.Add("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println(resp.Status)
		// log.Fatalf("Error happened Err: %s", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	var myStoredVariable map[string]any
	bodyString := string(body)
	json.Unmarshal([]byte(bodyString), &myStoredVariable)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(myStoredVariable)
}

func GetPlayerBattlelog(w http.ResponseWriter, r *http.Request) {
	playerId := chi.URLParam(r, "playerId")

	url := fmt.Sprintf("https://api.brawlstars.com/v1/players/%%23%s/battlelog", playerId)
	battleType := r.URL.Query().Get("type")
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	apiKey := goDotEnvVariable("API_KEY")
	req.Header.Add("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println(resp.Status)
		// log.Fatalf("Error happened Err: %s", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	var myStoredVariable map[string]any
	bodyString := string(body)
	json.Unmarshal([]byte(bodyString), &myStoredVariable)

	battles := myStoredVariable["items"].([]interface{})

	var result []map[string]any
	for _, v := range battles {
		values := v.(map[string]any)
		battle := values["battle"].(map[string]any)
		if battle["type"] == battleType {
			result = append(result, values)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}
