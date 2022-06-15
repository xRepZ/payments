package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	logger "github.com/sirupsen/logrus"
	"github.com/xRepZ/payments/internal"
)

const (
	login    = "admin"
	password = "admin"
	adminID  = "111"
)

var TokenGenerator = jwtauth.New("HS256", []byte("sssfsvasdf"), nil)

type AddRequest struct {
	UserId   int     `json:"userId"`
	Email    string  `json:"email"`
	Amount   float64 `json:"amount"`
	Сurrency string  `json:"currency"`
}

type UpdateRequest struct {
	Status string `json:"status"`
}

type Uadmin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var uadmin = Uadmin{
	Login:    "admin",
	Password: "admin",
}

type Server struct {
	storage internal.TransactionStorage
}

func NewServer(storage internal.TransactionStorage) *Server {
	return &Server{storage: storage}
}

func (s *Server) Login(w http.ResponseWriter, req *http.Request) {
	log := logger.New()
	user := &Uadmin{}
	err := json.NewDecoder(req.Body).Decode(user)
	if err != nil {
		http.Error(w, http.StatusText(402), 402)
		log.Infof("can't decode json: %s", err)
		return
	}
	if uadmin.Login != user.Login || uadmin.Password != user.Password {
		http.Error(w, http.StatusText(403), 403)
		log.Infof("Wrong login or password")
		return
	}
	_, tokenString, _ := TokenGenerator.Encode(map[string]interface{}{"id": adminID})
	err = json.NewEncoder(w).Encode(struct {
		Token string `json:"token"`
	}{Token: tokenString})
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Infof("can't encode json: %s", err)
		return
	}

}

func (s *Server) AddTransaction(w http.ResponseWriter, req *http.Request) {
	// парсим входящий json
	log := logger.New()
	data := &AddRequest{}
	err := json.NewDecoder(req.Body).Decode(data)
	if err != nil {
		http.Error(w, http.StatusText(402), 402)
		log.Infof("can't decode json: %s", err)
		return
	}
	//log.Infof("userId %v, email %s, amunt %v, currancy %s", data.UserId, data.Email, data.Amount, data.Сurrency)
	err = s.storage.AddTransaction(req.Context(), data.UserId, data.Email, data.Amount, data.Сurrency)
	if err != nil {
		log.Infof("can't add to db: %s", err)
		return
	}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Infof("can't encode json: %s", err)
		return
	}

	log.Info("done")
	// Получаем transaction_id из URL
	// id := chi.URLParam(req, "transaction_id")

	// отправляем ответ
	// err = json.NewEncoder(w).Encode(data)
}

func (s *Server) UpdateById(w http.ResponseWriter, req *http.Request) {
	log := logger.New()

	id := chi.URLParam(req, "transaction_id")
	intId, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, http.StatusText(403), 403)
		log.Infof("can't convert id to int: %s", err)
		return
	}
	if intId <= 0 {
		http.Error(w, http.StatusText(403), 403)
		log.Infof("id: %v, is missing", intId)
		return
	}
	_, claims, _ := jwtauth.FromContext(req.Context())
	ids := claims["id"]

	if ids == "" {
		// обработать
	}

	if ids != adminID {
		// return forbidden
	}

	status := &UpdateRequest{}
	err = json.NewDecoder(req.Body).Decode(status)
	if err != nil {
		log.Infof("can't decode json: %s", err)
		return
	}
	err = s.storage.UpdateById(req.Context(), intId, status.Status)
	if err != nil {
		log.Infof("can't add to db: %s", err)
		return
	}
	err = json.NewEncoder(w).Encode(status.Status)
	if err != nil {
		log.Infof("can't encode json: %s", err)
		return
	}

}

func (s *Server) CancelById(w http.ResponseWriter, req *http.Request) {
	log := logger.New()
	id := chi.URLParam(req, "transaction_id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, http.StatusText(403), 403)
		log.Infof("can't convert id to int: %s", err)
		return
	}
	if intId <= 0 {
		http.Error(w, http.StatusText(403), 403)
		log.Infof("id: %v, is missing", intId)
		return
	}

	err = s.storage.CancelById(req.Context(), intId)
	if err != nil {
		log.Infof("can't delet from db: %s", err)
		return
	}

}
func (s *Server) GetByMail(w http.ResponseWriter, req *http.Request) {
	log := logger.New()

	mail := chi.URLParam(req, "user_email")
	log.Info("done")
	// if uId <= 0 {
	// 	http.Error(w, http.StatusText(403), 403)
	// 	log.Infof("id: %v, is missing", uId)
	// 	return
	// }
	log.Info("done")
	userT, _ := s.storage.GetByMail(req.Context(), mail)

	log.Infof("done, %v", userT)
	//for i := range userT {
	//fmt.Println(userT[i])

	_ = json.NewEncoder(w).Encode(struct {
		Transactions []*internal.Transactions `json:"transactions"`
	}{Transactions: userT})

}

func (s *Server) GetByUserId(w http.ResponseWriter, req *http.Request) {
	log := logger.New()

	id := chi.URLParam(req, "user_id")
	uId, err := strconv.Atoi(id)
	log.Info("done")
	if err != nil {
		http.Error(w, http.StatusText(403), 403)
		log.Infof("can't convert id to int: %s", err)
		return
	}
	if uId <= 0 {
		http.Error(w, http.StatusText(403), 403)
		log.Infof("id: %v, is missing", uId)
		return
	}
	log.Info("done")
	//intId, _ := strconv.Atoi(id)
	//status "" -- ошибка
	userT, _ := s.storage.GetByUserId(req.Context(), uId)

	log.Infof("done, %v", userT)
	//for i := range userT {
	//fmt.Println(userT[i])

	err = json.NewEncoder(w).Encode(struct {
		Transactions []*internal.Transactions `json:"transactions"`
	}{Transactions: userT})
	if err != nil {
		//http.Error(w, http.StatusText(403), 403)
		log.Infof("can't ger user info: %s", err)
		return
	}

}

// GET /api/transaction/{transaction_id}
func (s *Server) GetById(w http.ResponseWriter, req *http.Request) {
	log := logger.New()

	id := chi.URLParam(req, "transaction_id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, http.StatusText(403), 403)
		log.Infof("can't convert id to int: %s", err)
		return
	}
	if intId <= 0 {
		http.Error(w, http.StatusText(403), 403)
		log.Infof("id: %v, is missing", intId)
		return
	}

	//intId, _ := strconv.Atoi(id)
	//status "" -- ошибка
	status, _ := s.storage.GetById(req.Context(), intId)

	_ = json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{Status: status})
}
