package main

import (
	"context"

	"log"

	"time"
)

type LoggingDecorator struct {
	Repo HTTPCheckRepository
}

func (ld *LoggingDecorator) PingDB() error {
	err := ld.Repo.PingDB()
	if err != nil {
		log.Printf("PingDB failed: %v", err)
	} else {
		log.Println("PingDB succeeded")
	}
	return err
}
func (ld *LoggingDecorator) GetAll() ([]HTTPCheck, error) {
	start := time.Now()
	checks, err := ld.Repo.GetAll()
	elapsed := time.Since(start)
	log.Printf("GetAll took %s", elapsed)
	return checks, err
}

func (ld *LoggingDecorator) GetByID(id string, ctx context.Context) (*HTTPCheck, error) {
	start := time.Now()
	check, err := ld.Repo.GetByID(id, ctx)
	elapsed := time.Since(start)
	log.Printf("GetByID took %s", elapsed)
	return check, err
}

func (ld *LoggingDecorator) Create(check *HTTPCheck, ctx context.Context) error {
	start := time.Now()
	err := ld.Repo.Create(check, ctx)
	elapsed := time.Since(start)
	log.Printf("Create took %s", elapsed)
	return err
}

func (ld *LoggingDecorator) Update(id string, check *HTTPCheck, ctx context.Context) error {
	start := time.Now()
	err := ld.Repo.Update(id, check, ctx)
	elapsed := time.Since(start)
	log.Printf("Update took %s", elapsed)
	return err
}

func (ld *LoggingDecorator) Delete(id string, ctx context.Context) error {
	start := time.Now()
	err := ld.Repo.Delete(id, ctx)
	elapsed := time.Since(start)
	log.Printf("Delete took %s", elapsed)
	return err
}
