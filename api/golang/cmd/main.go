package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"os"
	"sync"

	"github.com/marceljk/pv_tracker/api/golang/internal/cronjob"
	firebaserealtimedb "github.com/marceljk/pv_tracker/api/golang/internal/firebase-realtime-db"
	"github.com/marceljk/pv_tracker/api/golang/internal/model"
	"github.com/marceljk/pv_tracker/api/golang/internal/solcast"
	sunnytripower "github.com/marceljk/pv_tracker/api/golang/internal/sunny-tripower"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func main() {
	var cfg model.Config
	err := model.LoadEnvs(&cfg)
	if err != nil {
		fmt.Printf("failed to load envs: %v", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	kickCronJobs(cfg)

	wg.Wait()

	// ticker := time.NewTicker(time.Nanosecond)
	// quit := make(chan struct{})
	// go func() {
	// 	defer wg.Done()
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			start(cfg)
	// 			wg.Done()
	// 		case <-quit:
	// 			ticker.Stop()
	// 		}
	// 	}
	// }()
	// wg.Wait()
}

func kickCronJobs(cfg model.Config) error {
	ctx := context.Background()

	// Firebase init
	config := &firebase.Config{
		ProjectID:   cfg.FirebaseProjectId,
		DatabaseURL: cfg.FirebaseDatabaseUrl,
	}
	opt := option.WithCredentialsFile(cfg.FirebaseCredentialsPath)
	db, err := firebaserealtimedb.NewFirebaseDbClient(ctx, config, opt)
	if err != nil {
		return fmt.Errorf("failed to initialize firebase db: %w", err)
	}

	// Inverter init
	repo := sunnytripower.NewRepo(cfg.SmaBaseUrl)

	solcastClient := solcast.NewSolcastClient(cfg.SolcastEndpoint, cfg.SolcastApiKey, false)

	c := cronjob.NewCronjob(repo, db, solcastClient)
	c.Start()
	return nil
}
