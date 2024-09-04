package server

import (
	"agilityfc-bot/internal/bot"
	"agilityfc-bot/internal/dynamo"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Server struct holds the bot instance and DynamoDB service
type Server struct {
	Bot       *bot.Bot
	DynamoSvc *dynamo.DynamoDBService
}

// NewServer creates a new Server instance
func NewServer(bot *bot.Bot, dynamo *dynamo.DynamoDBService) *Server {
	return &Server{
		Bot:       bot,
		DynamoSvc: dynamo,
	}
}

// StartHTTPServer starts the HTTP server and listens for requests
func (s *Server) StartHTTPServer() {
	http.HandleFunc("/send-message-chat", s.handleSendMessage)
	http.HandleFunc("/send-message-user", s.handleSendUserMessage)
	http.HandleFunc("/save-user", s.handleSaveUser)
	http.HandleFunc("/save-snapshot", s.handleSaveSnapshot)

	fmt.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}
}

// handleSendMessage handles the /send-message-chat endpoint
func (s *Server) handleSendMessage(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	if message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	channelID := "1270770079486840873" // Replace with your Discord channel ID
	err := s.Bot.SendMessage(channelID, message)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Message sent: %s", message)
}

// handleSendMessage handles the /send-message-user endpoint
func (s *Server) handleSendUserMessage(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	if message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	channelID := "1271533422006767627" // Replace with your Discord channel ID
	err := s.Bot.SendMessage(channelID, message)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Message sent: %s", message)
}

// handleSaveUser handles the /save-user endpoint
func (s *Server) handleSaveUser(w http.ResponseWriter, r *http.Request) {
	var user dynamo.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Printf(err.Error())
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user.LastSeen = time.Now().UTC().Format(time.RFC3339)

	err = s.DynamoSvc.PutUser(user)
	if err != nil {
		fmt.Printf(err.Error())
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User saved: %s", user.UserID)
}

// handleSaveSnapshot handles the /save-snapshot endpoint
func (s *Server) handleSaveSnapshot(w http.ResponseWriter, r *http.Request) {
	var snapshot dynamo.Snapshot
	err := json.NewDecoder(r.Body).Decode(&snapshot)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	snapshot.Timestamp = time.Now().UTC().Format(time.RFC3339)

	err = s.DynamoSvc.PutSnapshot(snapshot)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Failed to save snapshot", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Snapshot saved: %s", snapshot.Timestamp)
}
