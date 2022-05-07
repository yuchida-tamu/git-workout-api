package record

import (
	"context"
	"fmt"
)

type Record struct {
	ID          string
	DateCreated string
	Message     string
	UserId      string
}

type Store interface {
	GetRecordsByUserId(context.Context, string) ([]Record, error)
	GetRecordById(context.Context, string) (Record, error)
	PostRecord(context.Context, Record) (Record, error)
	UpdateRecord(ctx context.Context, ID string, rcd Record) (Record, error)
	DeleteRecord(context.Context, string) error
}

type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetRecordsByUserId(ctx context.Context, ID string) ([]Record, error) {
	rcd, err := s.Store.GetRecordsByUserId(ctx, ID)
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
	postedRecord, err := s.Store.PostRecord(ctx, rcd)
	if err != nil {
		fmt.Println(err)
		return postedRecord, err
	}

	return postedRecord, nil
}

func (s *Service) UpdateRecord(ctx context.Context, ID string, rcd Record) (Record, error) {
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
