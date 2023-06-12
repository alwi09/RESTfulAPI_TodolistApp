package dto

import "todolist_gin_gorm/internal/model/entity"

type TodolistResponseCreate struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type TodolistResponseGetAll struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	More    int            `json:"more"`
	Data    []entity.Todos `json:"data"`
}

type TodolistResponseGetID struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    entity.Todos `json:"data"`
}

type TodolistResponseDelete struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type TodolistResponseUpdate struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CreateUserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type UserLoginResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}
