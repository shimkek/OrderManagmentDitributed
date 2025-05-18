package main

import "context"

type store struct {
	//mongodb instance
}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(ctx context.Context) error {
	return nil
}

func (s *store) Update(ctx context.Context) error {
	return nil
}
func (s *store) Get(ctx context.Context) error {
	return nil
}
func (s *store) Delete(ctx context.Context) error {
	return nil
}
