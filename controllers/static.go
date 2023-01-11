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

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func getDb(config PostgresConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.String()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return db, nil
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "password",
		Database: "brawl",
		SSLMode:  "disable",
	}
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func StaticHandlerJSON(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	cfg := DefaultPostgresConfig()
	db, err := getDb(cfg)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	db.Create(&Product{Code: "D42", Price: 100})

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
