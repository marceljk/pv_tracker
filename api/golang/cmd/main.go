package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"os"
	"sync"

	firebase "firebase.google.com/go/v4"
	"github.com/marceljk/pv_tracker/api/golang/internal/cronjob"
	firebaserealtimedb "github.com/marceljk/pv_tracker/api/golang/internal/firebase-realtime-db"
	"github.com/marceljk/pv_tracker/api/golang/internal/model"
	"github.com/marceljk/pv_tracker/api/golang/internal/solcast"
	"github.com/marceljk/pv_tracker/api/golang/internal/varta"
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

	if err := startCronJobs(cfg); err != nil {
		fmt.Printf("failed to start: %v", err)
		os.Exit(1)
	}

	wg.Wait()
}

func startCronJobs(cfg model.Config) error {
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

	// Varta init
	vartaClient := varta.NewVartaRepo(cfg.VartaUsername, cfg.VartaPassword)
	if err := vartaClient.Login(); err != nil {
		return fmt.Errorf("failed to login to varta: %w", err)
	}

	solcastClient := solcast.NewSolcastClient(cfg.SolcastEndpoint, cfg.SolcastApiKey, false)

	c := cronjob.NewCronjob(vartaClient, db, solcastClient)
	c.Start()
	return nil
}
