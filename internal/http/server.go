package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/xRepZ/payments/internal"
)

type AddRequest struct {
	UserId   int     `json:"userId"`
	Email    string  `json:"email"`
	Amount   float64 `json:"amount"`
	Сurrency string  `json:"currency"`
}

type Server struct {
	storage internal.TransactionStorage
}

func NewServer(storage internal.TransactionStorage) *Server {
	return &Server{storage: storage}
}

func (s *Server) AddTransaction(w http.ResponseWriter, req *http.Request) {
	// парсим входящий json
	// data := &AddRequest{}
	// err := json.NewDecoder(req.Body).Decode(data)
	// if err != nil {

	// }

	// Получаем transaction_id из URL
	// id := chi.URLParam(req, "transaction_id")

	// отправляем ответ
	// err = json.NewEncoder(w).Encode(ourStruct)
}

// GET /api/transaction/{transaction_id}
func (s *Server) GetById(w http.ResponseWriter, req *http.Request) {
	//log := logger.New()

	id := chi.URLParam(req, "transaction_id")
	// if id <= 0 {
	// 	// TODO add error
	// }

	intId, _ := strconv.Atoi(id)

	status, _ := s.storage.GetById(req.Context(), intId)

	_ = json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{Status: status})
}
