package main

import "github.com/shimkek/omd-common/api"

type CreateOrderRequest struct {
	Order         *api.Order `json:"order"`
	RedirectToURL string     `json:"redirectToURL"`
}
