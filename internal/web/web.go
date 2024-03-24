package web

import (
	"encoding/json"
	"errors"
	"github.com/djcrock/fwip/internal/model"
	"github.com/djcrock/fwip/internal/repository"
	"github.com/djcrock/fwip/internal/web/static"
	"log"
	"net/http"
	"strconv"
)

type server struct {
	logger   *log.Logger
	repoPool *repository.Pool
}

func NewApp(
	logger *log.Logger,
	pool *repository.Pool,
) http.Handler {
	server := &server{
		logger:   logger,
		repoPool: pool,
	}
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", static.FileServer))
	mux.Handle("GET /style.css", static.FileServer)
	mux.Handle("GET /main.js", static.FileServer)
	mux.HandleFunc("GET /", static.HandleIndex)
	mux.HandleFunc("GET /titles", server.handleGetTitles)
	mux.HandleFunc("GET /titles/{id}", server.handleGetTitle)
	mux.HandleFunc("GET /services", server.handleGetServices)
	mux.HandleFunc("GET /services/{id}", server.handleGetService)
	mux.HandleFunc("POST /users", server.handlePostUsers)
	mux.HandleFunc("GET /users", server.handleGetUsers)
	mux.HandleFunc("GET /users/{id}", server.handleGetUser)
	mux.HandleFunc("POST /users/{id}/watch_history", server.handlePostUserWatchHistory)
	mux.HandleFunc("GET /users/{id}/watch_history", server.handleGetUserWatchHistory)

	// TODO: add a "-dev" flag to control this
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		//w.Header().Add("Access-Control-Allow-Credentials", "true")
		mux.ServeHTTP(w, r)
	})
}

func (s *server) handleGetTitles(w http.ResponseWriter, r *http.Request) {
	repo := s.repoPool.GetRepository(r.Context())
	defer s.repoPool.PutRepository(repo)

	var titles []*model.Title
	var err error
	serviceIdStr := r.URL.Query().Get("service")
	if serviceIdStr != "" {
		serviceId, err := strconv.ParseInt(serviceIdStr, 10, 64)
		if err != nil {
			s.logger.Printf("invalid service id: %v", err)
			http.Error(w, "invalid service id", http.StatusInternalServerError)
		}
		titles, err = repo.GetTitlesByService(serviceId)
	} else {
		titles, err = repo.GetTitles()
	}

	if err != nil {
		s.logger.Printf("failed to retrieve titles: %v", err)
		http.Error(w, "failed to retrieve titles", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&titles)
	if err != nil {
		s.logger.Printf("failed to serialize titles: %v", err)
	}
}

func (s *server) handleGetTitle(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		s.logger.Printf("invalid id: %v", err)
		http.Error(w, "invalid id", http.StatusInternalServerError)
	}

	repo := s.repoPool.GetRepository(r.Context())
	defer s.repoPool.PutRepository(repo)

	title, err := repo.GetTitle(id)
	if err != nil {
		if errors.Is(err, repository.ErrNoSuchTitle) {
			s.logger.Printf("title not found: `%d`", id)
			http.Error(w, "title not found", http.StatusNotFound)
		} else {
			s.logger.Printf("failed to retrieve title `%d`: %v", id, err)
			http.Error(w, "failed to retrieve title", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&title)
	if err != nil {
		s.logger.Printf("failed to serialize title: %v", err)
	}
}

func (s *server) handleGetServices(w http.ResponseWriter, r *http.Request) {
	repo := s.repoPool.GetRepository(r.Context())
	defer s.repoPool.PutRepository(repo)

	services, err := repo.GetServices()
	if err != nil {
		s.logger.Printf("failed to retrieve services: %v", err)
		http.Error(w, "failed to retrieve services", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&services)
	if err != nil {
		s.logger.Printf("failed to serialize services: %v", err)
	}
}

func (s *server) handleGetService(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		s.logger.Printf("invalid id: %v", err)
		http.Error(w, "invalid id", http.StatusInternalServerError)
	}

	repo := s.repoPool.GetRepository(r.Context())
	defer s.repoPool.PutRepository(repo)

	service, err := repo.GetService(id)
	if err != nil {
		if errors.Is(err, repository.ErrNoSuchService) {
			s.logger.Printf("service not found: `%d`", id)
			http.Error(w, "service not found", http.StatusNotFound)
		} else {
			s.logger.Printf("failed to retrieve service `%d`: %v", id, err)
			http.Error(w, "failed to retrieve service", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&service)
	if err != nil {
		s.logger.Printf("failed to serialize service: %v", err)
	}
}

func (s *server) handlePostUsers(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		s.logger.Printf("malformed user: %v", err)
		http.Error(w, "malformed user", http.StatusBadRequest)
		return
	}
	if user.Username == "" {
		s.logger.Printf("malformed user: missing username")
		http.Error(w, "malformed user: missing username", http.StatusBadRequest)
		return
	}

	repo := s.repoPool.GetRepository(r.Context())
	defer s.repoPool.PutRepository(repo)

	_, err = repo.PutUser(user)
	if err != nil {
		s.logger.Printf("failed to create user: %v", err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		s.logger.Printf("failed to serialize user: %v", err)
	}
}

func (s *server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	repo := s.repoPool.GetRepository(r.Context())
	defer s.repoPool.PutRepository(repo)

	users, err := repo.GetUsers()
	if err != nil {
		s.logger.Printf("failed to retrieve users: %v", err)
		http.Error(w, "failed to retrieve users", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&users)
	if err != nil {
		s.logger.Printf("failed to serialize users: %v", err)
	}
}

func (s *server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		s.logger.Printf("invalid id: %v", err)
		http.Error(w, "invalid id", http.StatusInternalServerError)
	}

	repo := s.repoPool.GetRepository(r.Context())
	defer s.repoPool.PutRepository(repo)

	user, err := repo.GetUser(id)
	if err != nil {
		if errors.Is(err, repository.ErrNoSuchUser) {
			s.logger.Printf("user not found: `%d`", id)
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			s.logger.Printf("failed to retrieve user `%d`: %v", id, err)
			http.Error(w, "failed to retrieve user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		s.logger.Printf("failed to serialize user: %v", err)
	}
}

func (s *server) handlePostUserWatchHistory(w http.ResponseWriter, r *http.Request) {
	var watchHistory *model.WatchHistory
	err := json.NewDecoder(r.Body).Decode(&watchHistory)
	if err != nil {
		s.logger.Printf("malformed watch history: %v", err)
		http.Error(w, "malformed watch history", http.StatusBadRequest)
		return
	}

	repo := s.repoPool.GetRepository(r.Context())
	defer s.repoPool.PutRepository(repo)

	err = repo.PutWatchHistory(watchHistory)
	if err != nil {
		s.logger.Printf("failed to create watch history: %v", err)
		http.Error(w, "failed to create watch history", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(&watchHistory)
	if err != nil {
		s.logger.Printf("failed to serialize watch history: %v", err)
	}
}

func (s *server) handleGetUserWatchHistory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		s.logger.Printf("invalid id: %v", err)
		http.Error(w, "invalid id", http.StatusInternalServerError)
	}

	repo := s.repoPool.GetRepository(r.Context())
	defer s.repoPool.PutRepository(repo)

	user, err := repo.GetUser(id)
	if err != nil {
		if errors.Is(err, repository.ErrNoSuchUser) {
			s.logger.Printf("user not found: `%d`", id)
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			s.logger.Printf("failed to retrieve user `%d`: %v", id, err)
			http.Error(w, "failed to retrieve user", http.StatusInternalServerError)
		}
		return
	}

	watchHistory, err := repo.GetUserWatchHistory(user.Id)
	if err != nil {
		s.logger.Printf("failed to retrieve watch history: %v", err)
		http.Error(w, "failed to retrieve watch history", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&watchHistory)
	if err != nil {
		s.logger.Printf("failed to serialize watch history: %v", err)
	}
}
