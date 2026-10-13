package server

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/takashi002/Go-nami/internal/store"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

type Server struct {
	store *store.Store
	tmpl  *template.Template
	mux   *http.ServeMux
}

func New(st *store.Store) (*Server, error) {
	funcs := template.FuncMap{
		"moodLabel": moodLabel,
		"moodJP":    moodJP,
		"formatTime": func(t time.Time) string {
			return t.Local().Format("2006.01.02 15:04")
		},
	}

	tmpl, err := template.New("").Funcs(funcs).ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, err
	}

	s := &Server{store: st, tmpl: tmpl, mux: http.NewServeMux()}
	s.routes()
	return s, nil
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) routes() {
	static, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	s.mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))

	s.mux.HandleFunc("GET /{$}", s.handleHome)
	s.mux.HandleFunc("GET /journal", s.handleJournal)
	s.mux.HandleFunc("POST /entries", s.handleCreateEntry)
	s.mux.HandleFunc("POST /entries/{id}/delete", s.handleDeleteEntry)
	s.mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	s.render(w, "home", map[string]any{
		"Title": "Nami — 波の手帖",
	})
}

func (s *Server) handleJournal(w http.ResponseWriter, r *http.Request) {
	s.render(w, "journal", map[string]any{
		"Title":   "手帖 — Nami",
		"Entries": s.store.List(),
	})
}

func (s *Server) handleCreateEntry(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	body := strings.TrimSpace(r.FormValue("body"))
	if body == "" {
		http.Error(w, "本文を書いてください", http.StatusBadRequest)
		return
	}
	if len([]rune(body)) > 500 {
		http.Error(w, "500文字以内で書いてください", http.StatusBadRequest)
		return
	}

	mood, ok := store.ValidMood(r.FormValue("mood"))
	if !ok {
		mood = store.MoodCalm
	}

	s.store.Add(body, mood)

	if r.Header.Get("HX-Request") == "true" {
		s.render(w, "entries", map[string]any{"Entries": s.store.List()})
		return
	}
	http.Redirect(w, r, "/journal", http.StatusSeeOther)
}

func (s *Server) handleDeleteEntry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	s.store.Delete(id)

	if r.Header.Get("HX-Request") == "true" {
		s.render(w, "entries", map[string]any{"Entries": s.store.List()})
		return
	}
	http.Redirect(w, r, "/journal", http.StatusSeeOther)
}

func (s *Server) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("template %s: %v", name, err)
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func moodLabel(m store.Mood) string {
	switch m {
	case store.MoodCalm:
		return "凪"
	case store.MoodStirred:
		return "寄せ"
	case store.MoodStorm:
		return "嵐"
	default:
		return string(m)
	}
}

func moodJP(m store.Mood) string {
	return moodLabel(m)
}
