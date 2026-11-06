package main

import (
	"context"
	"fmt"
	"labflow/bootstrap"
	"labflow/mail"
	"labflow/queue"
	"labflow/seeders"
	"labflow/server"
	"labflow/worker"
	"log"
	"os"

	"github.com/hibiken/asynq"
)

func main() {
	bootstrap.LoadEnvVariables()
	db, err := bootstrap.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database")
	}
	redis, err := bootstrap.InitRedis()
	if err != nil {
		log.Fatalf("Failed to connect to redis")
	}

	defer func() {
		if err := redis.Close(); err != nil {
			log.Fatalf("Failed to close redis connection: %v", err)
		}
	}()
	s3client := bootstrap.ConnectToR2(context.Background())
	bootstrap.InitBuckets(context.Background(), s3client)

	mailer := mail.NewSender()

	redisURL := os.Getenv("REDIS_URL")

	enqueuer := queue.NewEnqueuer(redisURL)
	defer func() {
		if err := enqueuer.Close(); err != nil {
			log.Printf("failed to close enqueuer: %v", err)
		}
	}()

	asynqRedisOpt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		log.Fatalf("failed to parse REDIS_URL for asynq: %v", err)
	}
	asynqServer := asynq.NewServer(
		asynqRedisOpt,
		asynq.Config{
			Concurrency: 5,
			Queues:      map[string]int{"email": 10},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(worker.TypeSendQuotationEmail, worker.HandleSendQuotationEmailTask(mailer))
	mux.HandleFunc(worker.TypeSendOTPEmail, worker.HandleSendOTPEmailTask(mailer))

	if err := asynqServer.Start(mux); err != nil {
		log.Fatalf("[worker] asynq server error: %v", err)
	}
	defer asynqServer.Shutdown()

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		seeders.RunSeeders(db)
		log.Println("Database seeded.")
		return
	}

	server.Serve(bootstrap.CookieKey(), db, redis, mailer, enqueuer, s3client)

	fmt.Println("server started")
}
