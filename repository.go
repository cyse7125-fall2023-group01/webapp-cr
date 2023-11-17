package main

import (
	"context"

	"fmt"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type HTTPCheckRepository interface {
	GetAll() ([]HTTPCheck, error)
	GetByID(id string, ctx context.Context) (*HTTPCheck, error)
	Create(check *HTTPCheck, ctx context.Context) error
	Update(id string, check *HTTPCheck, ctx context.Context) error
	Delete(id string, ctx context.Context) error
	PingDB() error
}

type HTTPCheckGORMRepository struct {
	db *gorm.DB
}

func (hc *HTTPCheck) ValidateURI() bool {
	return uriPattern.MatchString(hc.URI)
}

func (r *HTTPCheckGORMRepository) GetAll() ([]HTTPCheck, error) {
	var checks []HTTPCheck
	if err := r.db.Find(&checks).Error; err != nil {
		return nil, err
	}
	return checks, nil
}

func (r *HTTPCheckGORMRepository) GetByID(id string, ctx context.Context) (*HTTPCheck, error) {
	err, check := r.performDBOperation(ctx, "GetByID", id, &HTTPCheck{})
	return check, err
}

func (r *HTTPCheckGORMRepository) Create(check *HTTPCheck, ctx context.Context) error {
	if !check.ValidateURI() {
		return fmt.Errorf("Invalid URI format")
	}
	check.ID = uuid.New()
	err, _ := r.performDBOperation(ctx, "Create", "", check)
	return err
}

func (r *HTTPCheckGORMRepository) Update(id string, check *HTTPCheck, ctx context.Context) error {
	err, _ := r.performDBOperation(ctx, "Update", id, check)
	return err
}

func (r *HTTPCheckGORMRepository) Delete(id string, ctx context.Context) error {
	err, _ := r.performDBOperation(ctx, "Delete", id, &HTTPCheck{})
	return err
}
func (r *HTTPCheckGORMRepository) PingDB() error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	err = db.Ping()
	return err
}
func (r *HTTPCheckGORMRepository) performDBOperation(ctx context.Context, operation string, id string, check *HTTPCheck) (error, *HTTPCheck) {

	resultChan := make(chan *HTTPCheck)
	errChan := make(chan error)

	go func() {
		db := r.db.WithContext(ctx)
		switch operation {
		case "GetByID":
			var tempResultCheck HTTPCheck // Create a temporary variable to store the result
			if err := db.First(&tempResultCheck, "id = ?", id).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					resultChan <- nil // Record not found
				} else {
					errChan <- err // Other errors
				}
			} else {
				resultChan <- &tempResultCheck // Send the result
			}
		case "Create":

			errChan <- db.Create(check).Error
		case "Update":
			err, existing := r.performDBOperation(ctx, "GetByID", id, check)
			if err != nil {
				errChan <- err
				return
			}
			if existing == nil {
				errChan <- fmt.Errorf("HTTP Check not found")
				return
			}
			check.ID = existing.ID
			if err := db.Save(check).Error; err != nil {
				errChan <- err
			} else {
				resultChan <- nil
			}
		case "Delete":
			err, existing := r.performDBOperation(ctx, "GetByID", id, check)
			if err != nil {
				errChan <- err
				return
			}
			if existing == nil {
				errChan <- fmt.Errorf("HTTP Check not found")
				return
			}
			if err := db.Delete(existing).Error; err != nil {
				errChan <- err
			} else {
				resultChan <- nil
			}
		}
	}()

	select {
	case check := <-resultChan:
		return nil, check
	case err := <-errChan:
		return err, nil
	case <-ctx.Done():
		return fmt.Errorf("Request canceled"), nil
	}
}

type HTTPCheckRepositoryDecorator interface {
	GetAll() ([]HTTPCheck, error)
	GetByID(id string, ctx context.Context) (*HTTPCheck, error)
	Create(check *HTTPCheck, ctx context.Context) error
	Update(id string, check *HTTPCheck, ctx context.Context) error
	Delete(id string, ctx context.Context) error
	PingDB() error
}
