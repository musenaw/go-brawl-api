package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"musenaw/go-brawl-api/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/go-redis/redis/v8"
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

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func StaticHandlerJSON(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	db := models.DB

	us := models.UserService{
		DB: db,
	}
	err := us.Migrate()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// db.Create(&Product{Code: "D42", Price: 100})

	resp := make(map[string]string)
	resp["message"] = "Working"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		panic(err)
	}
	w.Write(jsonResp)
}

// Service service
type Service struct {
	client *redis.Client
}

func NewRedisClient() Service {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	newService := Service{client: rdb}
	return newService
}

// Set sets key value
func (s *Service) Set(ctx context.Context, key string, value interface{}) error {
	exp := time.Duration(20 * time.Second)
	return s.client.Set(ctx, key, value, exp).Err()
}

// Get key value
func (s *Service) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

func GetPlayerInfo(w http.ResponseWriter, r *http.Request) {
	playerId := chi.URLParam(r, "playerId")

	ctx := context.Background()
	redisClient := NewRedisClient()
	val, _ := redisClient.Get(ctx, playerId)
	if val != "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(val))
		return
	}

	fmt.Print("LLEGOOOO")
	fmt.Println(val == "")

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
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		panic(err)
	}
	var newUserData models.UserInput
	bodyString := string(body)
	json.Unmarshal([]byte(bodyString), &newUserData)

	db := models.DB

	us := models.UserService{
		DB: db,
	}
	userData := models.User(newUserData)
	err = us.CreateOrUpdate(&userData)

	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	marhsaledUser, _ := json.Marshal(userData)
	errDB := redisClient.Set(ctx, playerId, marhsaledUser)
	fmt.Println(errDB)
	fmt.Println(marhsaledUser)
	if err := json.NewEncoder(w).Encode(userData); err != nil {
		panic(err)
	}
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
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
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
