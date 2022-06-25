package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Server struct {
	Port    string
	Usecase ItemUsecase
}

type ServerConfig struct {
	Port string
}

func NewServer(config ServerConfig, usecase ItemUsecase) *Server {
	return &Server{
		Port:    config.Port,
		Usecase: usecase,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/getall", s.GetAllHandler)
	mux.HandleFunc("/getitems", s.GetItemsHandler)
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
	items, err := s.Usecase.GetAll()
	if err != nil {
		e := "failed to GetAll"
		log.Printf("%s: %v", e, err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}
	responseItemByJSON(items, w)
}

func (s *Server) GetItemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		e := fmt.Sprintf("invalid method %s, request must be POST", r.Method)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "POST request must be JSON. check Content-Type", http.StatusBadRequest)
		return
	}

	var ids IDs
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&ids); err != nil {
		// 参考：https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)
		default:
			// 内部のロジックをさらさないようにあえて詳細は返さない
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	items, err := s.Usecase.GetItems(ids.IDs)
	if err != nil {
		e := "failed to GetItems"
		log.Printf("%s: %v", e, err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}
	responseItemByJSON(items, w)
}

func responseItemByJSON(is []Item, w http.ResponseWriter) {
	jsonData, err := json.Marshal(is)
	if err != nil {
		log.Println("failed to marshal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 省略可能
	fmt.Fprintln(w, string(jsonData))
}