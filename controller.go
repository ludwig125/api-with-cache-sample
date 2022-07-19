package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Port    string
	Usecase SearchUsecase
}

type ServerConfig struct {
	Port string
}

func NewServer(config ServerConfig, usecase SearchUsecase) *Server {
	return &Server{
		Port:    config.Port,
		Usecase: usecase,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/getall", s.GetAllHandler)
	mux.HandleFunc("/search", s.SearchHandler)
	// mux.HandleFunc("/getitems", s.GetScoresHandler)
	// mux.HandleFunc("/find", s.FindHandler)
	// mux.HandleFunc("/update", s.UpdateHandler) //POST
	// mux.HandleFunc("/delete", s.DeleteHandler) //POST
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("ok")); err != nil {
			log.Println("error occurred:", err)
		}
	})
	srv := &http.Server{
		Addr:    "localhost:" + s.Port,
		Handler: mux,
	}
	fmt.Println("starting http server on :", s.Port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("Server closed with error: %v", err)
	}
	return nil
}

func (s *Server) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	is, err := s.Usecase.GetAll()
	if err != nil {
		e := "failed to GetAll"
		log.Printf("%s: %v", e, err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}
	responseByJSON(is, w)
}

func (s *Server) SearchHandler(w http.ResponseWriter, r *http.Request) {
	cond, err := NewSearchCondition(r.URL.Query().Get("price"),
		r.URL.Query().Get("expr"),
	)
	if err != nil {
		e := "failed to NewSearchCondition"
		log.Printf("%s: %v", e, err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	is, err := s.Usecase.Search(*cond)
	if err != nil {
		e := "failed to Search"
		log.Printf("%s: %v", e, err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}
	responseByJSON(is, w)
}

// // 取得方法
// // curl -X POST -H "Content-Type: application/json" -d '{"ids":[1,2,3,4,5]}' http://localhost:8080/getitems
// func (s *Server) GetScoresHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		e := fmt.Sprintf("invalid method %s, request must be POST", r.Method)
// 		http.Error(w, e, http.StatusBadRequest)
// 		return
// 	}

// 	if r.Header.Get("Content-Type") != "application/json" {
// 		http.Error(w, "POST request must be JSON. check Content-Type", http.StatusBadRequest)
// 		return
// 	}

// 	var ids RequestIDs
// 	dec := json.NewDecoder(r.Body)
// 	if err := dec.Decode(&ids); err != nil {
// 		// 参考：https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
// 		var syntaxError *json.SyntaxError
// 		var unmarshalTypeError *json.UnmarshalTypeError

// 		switch {
// 		case errors.As(err, &syntaxError):
// 			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
// 			http.Error(w, msg, http.StatusBadRequest)
// 		case errors.As(err, &unmarshalTypeError):
// 			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
// 			http.Error(w, msg, http.StatusBadRequest)
// 		case errors.Is(err, io.EOF):
// 			msg := "Request body must not be empty"
// 			http.Error(w, msg, http.StatusBadRequest)
// 		default:
// 			// 内部のロジックをさらさないようにあえて詳細は返さない
// 			log.Println(err.Error())
// 			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	items, err := s.Usecase.GetScores(ids.IDs)
// 	if err != nil {
// 		e := "failed to GetScores"
// 		log.Printf("%s: %v", e, err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		http.Error(w, e, http.StatusInternalServerError)
// 		return
// 	}
// 	responseByJSON(items, w)
// }

func responseByJSON(resp interface{}, w http.ResponseWriter) {
	jsonData, err := json.Marshal(resp)
	if err != nil {
		log.Println("failed to marshal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 省略可能
	fmt.Fprintln(w, string(jsonData))
}
