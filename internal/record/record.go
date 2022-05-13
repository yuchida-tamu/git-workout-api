package record

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/yuchida-tamu/git-workout-api/internal/user"
)

type Record struct {
	ID          string
	DateCreated string
	MessageBody string
	Author      string
}

type Store interface {
	GetRecordsByAuthor(context.Context, string) ([]Record, error)
	GetRecordById(context.Context, string) (Record, error)
	PostRecord(context.Context, Record) (Record, error)
	UpdateRecord(ctx context.Context, ID string, rcd Record) (Record, error)
	DeleteRecord(context.Context, string) error

	GetUser(context.Context, string) (user.User, error)
}

type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetRecordsByAuthor(ctx context.Context, ID string) ([]Record, error) {
	rcd, err := s.Store.GetRecordsByAuthor(ctx, ID)
	if err != nil {
		fmt.Println(err)
		return []Record{}, err
	}

	return rcd, nil
}

func (s *Service) GetRecordById(ctx context.Context, ID string) (Record, error) {
	rcd, err := s.Store.GetRecordById(ctx, ID)
	if err != nil {
		fmt.Println(err)
		return Record{}, err
	}
	return rcd, nil
}

func (s *Service) PostRecord(ctx context.Context, rcd Record) (Record, error) {
	// validate uuid format
	if _, err := uuid.FromString(rcd.Author); err != nil {
		fmt.Println(err)
		return Record{}, err
	}
	// check if the user already exists
	if _, err := s.Store.GetUser(ctx, rcd.Author); err != nil {
		fmt.Print(err)
		return Record{}, err
	}

	postedRecord, err := s.Store.PostRecord(ctx, rcd)
	if err != nil {
		fmt.Println(err)
		return Record{}, err
	}

	return postedRecord, nil
}

func (s *Service) UpdateRecord(ctx context.Context, ID string, rcd Record) (Record, error) {
	// validate uuid format
	if _, err := uuid.FromString(rcd.Author); err != nil {
		fmt.Println(err)
		return Record{}, err
	}
	// check if the user already exists
	if _, err := s.Store.GetUser(ctx, rcd.Author); err != nil {
		fmt.Print(err)
		return Record{}, err
	}

	updatedRecord, err := s.Store.UpdateRecord(ctx, ID, rcd)
	if err != nil {
		fmt.Println(err)
		return updatedRecord, err
	}

	return updatedRecord, nil
}

func (s *Service) DeleteRecord(ctx context.Context, ID string) error {
	err := s.Store.DeleteRecord(ctx, ID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
