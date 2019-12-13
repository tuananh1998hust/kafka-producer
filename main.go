package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	kafka "github.com/segmentio/kafka-go"
)

// HandlePost :
func HandlePost(w http.ResponseWriter, r *http.Request) {
	kafkaURL := os.Getenv("KAFKA_URL")
	topic := os.Getenv("KAFKA_TOPIC")
	// kafkaURL := "localhost:9092"
	// topic := "user-topic"
	conn := getKafkaWriter(kafkaURL, topic)
	defer conn.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("address-%s", r.RemoteAddr)),
		Value: body,
	}

	_, err = conn.WriteMessages(msg)

	if err != nil {
		log.Println("Kafka Write Message Error", err.Error())
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte("Ok"))
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Conn {
	conn, _ := kafka.DialLeader(context.Background(), "tcp", kafkaURL, topic, 0)

	return conn
}

func main() {
	r := mux.NewRouter()

	// POST /user
	r.HandleFunc("/user", HandlePost).Methods("POST")

	// Run the web server.
	log.Println("Server is start on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
