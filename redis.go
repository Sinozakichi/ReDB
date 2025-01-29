package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

const keyname = "cards"

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 的地址
		Password: "",               // 如果有密碼則填入
		DB:       0,                // 使用的資料庫編號（預設為 0）
	})

	// 測試連線
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("無法連線到 Redis: %v", err)
	}
	log.Println("Redis 已成功連線！")
}

func HandleCardsRequestToRedis(w http.ResponseWriter, r *http.Request) {

	// 從 Redis 嘗試獲取資料
	cardsJSON, err := redisClient.Get(ctx, keyname).Result()
	if err != nil {
		http.Error(w, "Failed to fetch cards from Redis", http.StatusInternalServerError)
		return
	}

	// 實作 Redis CRUD 操作
	// GET(SELECT)、POST(INSERT)、PUT(UPDATE)、DELETE(DELETE)
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cardsJSON))

	case http.MethodPost:
		card, cards, err := parseRedisDataToObject(w, r, cardsJSON)
		if err == nil {
			//拿取新的 ID(auto_increment的效果)
			card.ID, err = getNewCardID()
			// 添加新卡片
			if err == nil {
				cards = append(cards, card)
			}
		}
		updateRedis(ctx, keyname, cards, w)
		w.WriteHeader(http.StatusCreated)

	case http.MethodPut:
		updatedCard, cards, err := parseRedisDataToObject(w, r, cardsJSON)
		if err == nil {
			// 更新卡片
			for i, card := range cards {
				if card.ID == updatedCard.ID {
					cards[i] = updatedCard
					break
				}
			}
		}
		updateRedis(ctx, keyname, cards, w)
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		deleteCard, cards, err := parseRedisDataToObject(w, r, cardsJSON)
		if err == nil {
			// 刪除卡片
			for i, card := range cards {
				if card.ID == deleteCard.ID {
					cards = append(cards[:i], cards[i+1:]...)
					break
				}
			}
		}
		updateRedis(ctx, keyname, cards, w)
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func parseRedisDataToObject(w http.ResponseWriter, r *http.Request, cardsJSON string) (Card, []Card, error) {
	// 整理 Redis 中撈取的資料與傳入的資料
	// card為要異動的卡片的資料
	// cards為所有卡片的資料
	var card Card
	var cards []Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return card, cards, err
	}
	if cardsJSON != "" {
		if err := json.Unmarshal([]byte(cardsJSON), &cards); err != nil {
			http.Error(w, "Failed to parse cards from Redis", http.StatusInternalServerError)
			return card, cards, err
		}
	}
	return card, cards, nil
}

func updateRedis(ctx context.Context, tablename string, cards []Card, w http.ResponseWriter) error {
	updatedCardsJSON, err := json.Marshal(cards)
	if err != nil {
		http.Error(w, "Failed to encode cards to JSON", http.StatusInternalServerError)
		return err
	}
	// 更新 Redis
	if err := redisClient.Set(ctx, tablename, updatedCardsJSON, 0).Err(); err != nil {
		http.Error(w, "Failed to update cards in Redis", http.StatusInternalServerError)
		return err
	}
	return nil

}

// 取得新的 ID
// 在 Redis 中，沒有像 MySQL 那樣自動遞增 (AUTO_INCREMENT)的機制，但可以透過 Redis 的 遞增操作 (INCR 或 INCRBY) 來模擬類似的功能。這樣，每次新增資料時，都能確保 id 是唯一且自動遞增的。
func getNewCardID() (int, error) {
	id, err := redisClient.Incr(ctx, "card_id_counter").Result()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
