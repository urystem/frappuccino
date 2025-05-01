package models

import "errors"

var (
	ErrBadInput = errors.New("bad input") // 422
	ErrNotFound = errors.New("not found") // 404
	ErrConflict = errors.New("conflict")  // 409 used for post ing and menu
	// ErrContentType = errors.New("")

	ErrBadInputItems       = errors.New("items invalid")    // 400
	ErrNotFoundItems       = errors.New("items not found")  // 404 //for menu ings and product items
	ErrOrderNotEnoughItems = errors.New("items not enough") // 500 used for not enough invents for order
	ErrOrderStatusClosed     = errors.New("order is already closed")
)

// 200 OK
// 201 Created
// 202 Accepted in async ---
// 204 No Content

// 400 BadRequest
// 404 not found
// 409 Conflict
// 415 Unsupported Media Type /* is not json*/
// 422 	Unprocessable Entity Ошибка валидации тела запроса
// 424
// 500
