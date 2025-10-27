package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"main/pkg/game"
)

type Game struct {
	Word           string
	Letters        []string
	Wrong          int
	Alphabet       []string
	GuessedLetters map[string]bool
	Corrects       int
	Finished       bool
	Message        string
	SessionID      string
}

var (
	gameState = make(map[string]*Game)
)

func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func createNewGame(difficulty string) *Game {
	word := strings.ToLower(game.Getword(difficulty))
	letters := make([]string, len(word))
	for i := range word {
		letters[i] = "_"
	}

	alphabet := []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}

	sessionID := generateSessionID()
	return &Game{
		Word:           word,
		Letters:        letters,
		Wrong:          0,
		Alphabet:       alphabet,
		GuessedLetters: make(map[string]bool),
		Corrects:       0,
		Finished:       false,
		SessionID:      sessionID,
	}
}

func getOrCreateSession(w http.ResponseWriter, r *http.Request) *Game {
	cookie, err := r.Cookie("session_id")
	var sessionID string

	if err != nil || cookie.Value == "" {
		g := createNewGame("easy")
		gameState[g.SessionID] = g

		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   g.SessionID,
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		})
		return g
	}

	sessionID = cookie.Value
	g, exists := gameState[sessionID]

	if !exists {
		g = createNewGame("easy")
		gameState[g.SessionID] = g

		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   g.SessionID,
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		})
	}

	return g
}

func CreateServer() *http.Server {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/newgame", func(w http.ResponseWriter, r *http.Request) {
		g := createNewGame("easy")
		gameState[g.SessionID] = g

		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   g.SessionID,
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		g := getOrCreateSession(w, r)

		if r.Method == http.MethodPost && !g.Finished {
			guess := strings.ToLower(r.FormValue("letter"))

			if g.GuessedLetters[guess] {
				tmpl := template.Must(template.ParseFiles("./static/index.html"))
				tmpl.Execute(w, g)
				return
			}

			correct := false
			for i, ch := range g.Word {
				if string(ch) == guess {
					g.Letters[i] = guess
					correct = true
					g.Corrects++
				}
			}
			if !correct {
				g.Wrong++
			}

			g.GuessedLetters[guess] = true

			if g.Corrects == len(g.Word) {
				fmt.Printf("You won! The word was %s\n", g.Word)
				g.Message = "You won! The word was " + g.Word
				g.Finished = true
			}

			if g.Wrong == 8 {
				fmt.Printf("You lose! The word was %s\n", g.Word)
				g.Message = "You lose! The word was " + g.Word
				g.Finished = true
			}
		}

		funcMap := template.FuncMap{
			"lower": strings.ToLower,
		}
		tmpl := template.Must(template.New("index.html").Funcs(funcMap).ParseFiles("./static/index.html"))
		tmpl.Execute(w, g)
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
}
