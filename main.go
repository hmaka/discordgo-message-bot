package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Register the two new handler functions and corresponding URL patterns with
	// the servemux, in exactly the same way that we did before.

	pubsub := Newpubsub()

	defer pubsub.close()

	go launchWebServer(&pubsub)

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}
	Token := os.Getenv("DISCORD_TOKEN")
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		messageCreate(s, m, &pubsub)
	})

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

func launchWebServer(pubsub *Pubsub) {
	//messages := make([]string, 0)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		home(w, r, pubsub.subscribe())
	})

	server := &http.Server{Addr: ":8080", Handler: mux}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-sc
		log.Println("Shutting down server...")
		server.Shutdown(context.Background())
	}()

	log.Println("Starting server on :8080")
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
func home(w http.ResponseWriter, r *http.Request, msg_channel <-chan string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	defer func() {
		w.Write([]byte("event: close\ndata: Closing connection\n\n"))
		w.(http.Flusher).Flush()
	}()

	for msg := range msg_channel {
		// Send event to the client
		w.Write([]byte(msg + "\n\n"))

		// Flush the response writer to send the event immediately
		w.(http.Flusher).Flush()

	}

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate, pubsub *Pubsub) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	println(m.Content)

	// code to send message through channel to webserver
	pubsub.notifyAll(m.Author.GlobalName + ": " + m.Content)
}
