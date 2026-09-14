package graphql

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

import "github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"

type Resolver struct {
	uc task.UseCase
}

func NewResolver(uc task.UseCase) *Resolver {
	return &Resolver{uc: uc}
}
