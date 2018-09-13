package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dwarvesf/yggdrasil/toolkit"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	consul "github.com/hashicorp/consul/api"
	kafka "github.com/segmentio/kafka-go"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	consulClient, err := consul.NewClient(&consul.Config{
		Address: fmt.Sprintf("consul:8500"),
	})
	if err != nil {
		panic(err)
	}
	kafkaAddr, kafkaPort, err := toolkit.GetServiceAddress(consulClient, "kafka")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	type Request struct {
		Service string                 `json:"service"`
		Payload map[string]interface{} `json:"payload"`
	}

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body) //<--- here!
		if err != nil {
			logger.Log("error", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var req Request
		if err := json.Unmarshal(body, &req); err != nil {
			logger.Log("error", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, err := json.Marshal(req.Payload)
		if err != nil {
			logger.Log("error", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		kw := kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{fmt.Sprintf("%v:%v", kafkaAddr, kafkaPort)},
			Topic:   req.Service,
		})
		kw.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(req.Service),
				Value: data,
			},
		)
		kw.Close()

		logger.Log("status", fmt.Sprintf("sent [%v] request to kafka", req.Service))
		w.Write([]byte("ok"))
	})

	fmt.Println(fmt.Sprintf("Listening on port :%v", os.Getenv("PORT")))
	http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), r)
}
