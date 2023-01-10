package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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

type Battle struct {
	Duration   int
	Mode       string
	Result     string
	StarPlayer map[string]any
	Teams      [][]map[string]any
	Type       string
}

type Event struct {
	Id   int
	Map  string
	Mode string
}

type BattleData struct {
	Battle     Battle
	BattleTime string
	Event      Event
}

type Paging struct {
	Cursors map[string]any
}
type BattleLog struct {
	Items  []BattleData
	Paging Paging
}

func GetPlayerBattlelog(w http.ResponseWriter, r *http.Request) {
	playerId := chi.URLParam(r, "playerId")

	url := fmt.Sprintf("https://api.brawlstars.com/v1/players/%%23%s/battlelog", playerId)
	battleType := strings.Trim(r.URL.Query().Get("type"), " ")
	mode := strings.Trim(r.URL.Query().Get("mode"), " ")
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
	var myStoredVariable BattleLog
	bodyString := string(body)
	json.Unmarshal([]byte(bodyString), &myStoredVariable)

	battleData := myStoredVariable.Items

	result := make([]BattleData, 0)
	for _, v := range battleData {

		battle := v.Battle
		if battleType != "" && battle.Type != battleType {
			continue
		}

		event := v.Event
		if mode != "" && event.Mode != mode {
			continue
		}

		result = append(result, v)

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}
