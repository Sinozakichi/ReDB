package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Card 表示 cards 表中的一條記錄
// 結構體中的 json 標籤用於定義結構體欄位在進行 JSON 編碼和解碼時的名稱映射。也就是說，當你將一個 Card 結構體轉換成 JSON 格式，或者從 JSON 格式解析成 Card 結構體時，這些標籤會指定如何將欄位名稱對應到 JSON 鍵（key）。
type Card struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Level     int    `json:"level"`
	Attribute string `json:"attribute"`
	Race      string `json:"race"`
	Attack    int    `json:"attack"`
	Defense   int    `json:"defense"`
	Effect    string `json:"effect"`
}

// Database 是資料庫操作的介面
type Database interface {
	Close() error
	CreateTable() error
	InsertCard(card Card) error
	GetAllCards() ([]Card, error)
	UpdateCard(card Card) error
	DeleteCard(id int) error
}

func main() {

	// STEP1.初始化資料庫
	//db, err := NewSQLiteDatabase("./test.db") //SQLite
	db, err := NewMySQLDatabase("root", "Eo@e368619220", "localhost:3306", "go_redb") //MySQL(Local)
	//db, err := NewMySQLDatabase("test", "test", "123.192.158.69:3306", "go_redb") //MySQL(Remote)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 建立表格
	// err = db.CreateTable()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Table創建!")

	// 插入資料 //SQLite
	// cardsToInsert := []Card{
	// 	{Name: "青眼白龍", Level: 8, Attribute: "光", Race: "龍", Attack: 3000, Defense: 2500, Effect: ""},
	// 	{Name: "黑魔導", Level: 7, Attribute: "闇", Race: "魔法使", Attack: 2500, Defense: 2100, Effect: ""},
	// 	{Name: "真紅眼黑龍", Level: 7, Attribute: "闇", Race: "龍", Attack: 2400, Defense: 2000, Effect: ""},
	// }

	// for _, card := range cardsToInsert {
	// 	err = db.InsertCard(card)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// fmt.Println("資料新增成功!")

	// STEP2.查詢資料(Select)
	datas, err := db.GetAllCards()
	if err != nil {
		panic(err)
	}
	fmt.Println("Datas in database:")
	for _, card := range datas {
		fmt.Printf("ID: %d, Name: %s, Level: %d, Attribute: %s, Race: %s, Attack: %d, Defense: %d, Effect: %s\n",
			card.ID, card.Name, card.Level, card.Attribute, card.Race, card.Attack, card.Defense, card.Effect)
	}

	// STEP3.啟動 HTTP 伺服器
	// 只能處理"/cards"路徑的請求，無法向下兼容如"/cards/id"等的請求
	http.HandleFunc("/cards", func(w http.ResponseWriter, r *http.Request) {
		//允許跨域請求
		w.Header().Set("Access-Control-Allow-Origin", "*") // 允許所有來源
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		//HTTP 方法「簡單方法」與「非簡單方法」
		//簡單方法：請求方法為GET、HEAD、POST，且請求頭（Header）僅包含：Accept、Accept-Language、Content-Language、Content-Type（僅限於值：application/x-www-form-urlencoded, multipart/form-data, text/plain）
		//非簡單方法：請求方法為PUT、DELETE、OPTIONS、CONNECT、TRACE、PATCH，或者請求頭中包含自定義請求頭（例如：Content-Type: application/json、Authorization），或者請求攜帶跨域的憑證資訊（例如 Cookies）
		//若屬非簡單方法，瀏覽器會先發送一個 OPTIONS 請求，詢問伺服器是否允許進行跨域請求，伺服器回應允許後，瀏覽器才會發送真正的請求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent) // 回應成功的204狀態碼
			return
		}
		HandleCardsRequest(db, w, r)
	})

	fmt.Println("Server is running on port 5500...")
	log.Fatal(http.ListenAndServe(":5500", nil))

}

// HandleCardsRequest 處理前端頁面的 CRUD請求
func HandleCardsRequest(db Database, w http.ResponseWriter, r *http.Request) {
	// 實作 CRUD 操作
	// GET(SELECT)、POST(INSERT)、PUT(UPDATE)、DELETE(DELETE)
	switch r.Method {
	case http.MethodGet:
		// 取得所有卡片資料
		cards, err := db.GetAllCards()
		if err != nil {
			http.Error(w, "Failed to fetch cards", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(cards)
	case http.MethodPost:
		// 新增卡片
		var card Card
		if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if err := db.InsertCard(card); err != nil {
			http.Error(w, "Failed to insert card", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case http.MethodPut:
		// 更新卡片
		var card Card
		if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if err := db.UpdateCard(card); err != nil {
			http.Error(w, "Failed to update card", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		// 刪除卡片
		var card Card
		if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if err := db.DeleteCard(card.ID); err != nil {
			http.Error(w, "Failed to delete card", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
