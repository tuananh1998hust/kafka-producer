package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	kafka "github.com/segmentio/kafka-go"
)

// User :
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// HandlePost :
func HandlePost(w http.ResponseWriter, r *http.Request) {
	// kafkaURL := os.Getenv("KAFKA_URL")
	// topic := os.Getenv("KAFKA_TOPIC")
	kafkaURL := "localhost:29092"
	topic := "user-topic"
	// conn := getKafkaWriter(kafkaURL, topic)
	// defer conn.Close()
	newWriter := getKafkaWriter(kafkaURL, topic)

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	v, err := json.Marshal(user)

	if err != nil {
		log.Fatal(err)
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("address-%s", r.RemoteAddr)),
		Value: v,
	}

	// _, err = conn.WriteMessages(msg)
	err = newWriter.WriteMessages(context.Background(), msg)

	log.Println("Write data")

	if err != nil {
		log.Println("Kafka Write Message Error", err.Error())
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte("Ok"))
}

// func getKafkaWriter(kafkaURL, topic string) *kafka.Conn {
// 	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaURL, topic, 0)

// 	if err != nil {
// 		log.Panic("Connection Err", err)
// 	}

// 	return conn
// }

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaURL},
		Topic:   topic,
	})
}

func main() {
	r := mux.NewRouter()

	// POST /user
	r.HandleFunc("/user", HandlePost).Methods("POST")

	// Run the web server.
	log.Println("Server is start on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
