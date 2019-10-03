// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphql

import (
	"fmt"
	"io"
	"strconv"
)

type AuthMutation struct {
	Signup *SignupOut `json:"signup"`
}

type AuthQuery struct {
	Signin *SigninOut `json:"signin"`
}

type CursorIn struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Cursor string `json:"cursor"`
}

type CursorOut struct {
	Count  int    `json:"count"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	IsEnd  bool   `json:"isEnd"`
	Cursor string `json:"cursor"`
}

type MsMutation struct {
	New *NewOut `json:"new"`
}

type MsQuery struct {
	Search *SearchOut `json:"search"`
}

type NewOut struct {
	Status NewOutStatus `json:"status"`
	ID     string       `json:"id"`
}

type SearchOut struct {
	Status SearchOutStatus `json:"status"`
	ID     []string        `json:"id"`
	Cursor *CursorOut      `json:"cursor"`
}

type SigninOut struct {
	Status SigninOutStatus `json:"status"`
	Token  string          `json:"token"`
}

type SignupOut struct {
	Status SignupOutStatus `json:"status"`
}

type AuthenticatedRequestStatus string

const (
	AuthenticatedRequestStatusOk                  AuthenticatedRequestStatus = "OK"
	AuthenticatedRequestStatusForbidden           AuthenticatedRequestStatus = "FORBIDDEN"
	AuthenticatedRequestStatusNotFound            AuthenticatedRequestStatus = "NOT_FOUND"
	AuthenticatedRequestStatusBadRequest          AuthenticatedRequestStatus = "BAD_REQUEST"
	AuthenticatedRequestStatusServerInternalError AuthenticatedRequestStatus = "SERVER_INTERNAL_ERROR"
)

var AllAuthenticatedRequestStatus = []AuthenticatedRequestStatus{
	AuthenticatedRequestStatusOk,
	AuthenticatedRequestStatusForbidden,
	AuthenticatedRequestStatusNotFound,
	AuthenticatedRequestStatusBadRequest,
	AuthenticatedRequestStatusServerInternalError,
}

func (e AuthenticatedRequestStatus) IsValid() bool {
	switch e {
	case AuthenticatedRequestStatusOk, AuthenticatedRequestStatusForbidden, AuthenticatedRequestStatusNotFound, AuthenticatedRequestStatusBadRequest, AuthenticatedRequestStatusServerInternalError:
		return true
	}
	return false
}

func (e AuthenticatedRequestStatus) String() string {
	return string(e)
}

func (e *AuthenticatedRequestStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AuthenticatedRequestStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AuthenticatedRequestStatus", str)
	}
	return nil
}

func (e AuthenticatedRequestStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type NewOutStatus string

const (
	NewOutStatusOk                  NewOutStatus = "OK"
	NewOutStatusForbidden           NewOutStatus = "FORBIDDEN"
	NewOutStatusBadRequest          NewOutStatus = "BAD_REQUEST"
	NewOutStatusServerInternalError NewOutStatus = "SERVER_INTERNAL_ERROR"
)

var AllNewOutStatus = []NewOutStatus{
	NewOutStatusOk,
	NewOutStatusForbidden,
	NewOutStatusBadRequest,
	NewOutStatusServerInternalError,
}

func (e NewOutStatus) IsValid() bool {
	switch e {
	case NewOutStatusOk, NewOutStatusForbidden, NewOutStatusBadRequest, NewOutStatusServerInternalError:
		return true
	}
	return false
}

func (e NewOutStatus) String() string {
	return string(e)
}

func (e *NewOutStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = NewOutStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid NewOutStatus", str)
	}
	return nil
}

func (e NewOutStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type OrderIn string

const (
	OrderInAsc  OrderIn = "ASC"
	OrderInDesc OrderIn = "DESC"
)

var AllOrderIn = []OrderIn{
	OrderInAsc,
	OrderInDesc,
}

func (e OrderIn) IsValid() bool {
	switch e {
	case OrderInAsc, OrderInDesc:
		return true
	}
	return false
}

func (e OrderIn) String() string {
	return string(e)
}

func (e *OrderIn) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OrderIn(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OrderIn", str)
	}
	return nil
}

func (e OrderIn) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SearchOutStatus string

const (
	SearchOutStatusOk                  SearchOutStatus = "OK"
	SearchOutStatusForbidden           SearchOutStatus = "FORBIDDEN"
	SearchOutStatusNotFound            SearchOutStatus = "NOT_FOUND"
	SearchOutStatusBadRequest          SearchOutStatus = "BAD_REQUEST"
	SearchOutStatusServerInternalError SearchOutStatus = "SERVER_INTERNAL_ERROR"
)

var AllSearchOutStatus = []SearchOutStatus{
	SearchOutStatusOk,
	SearchOutStatusForbidden,
	SearchOutStatusNotFound,
	SearchOutStatusBadRequest,
	SearchOutStatusServerInternalError,
}

func (e SearchOutStatus) IsValid() bool {
	switch e {
	case SearchOutStatusOk, SearchOutStatusForbidden, SearchOutStatusNotFound, SearchOutStatusBadRequest, SearchOutStatusServerInternalError:
		return true
	}
	return false
}

func (e SearchOutStatus) String() string {
	return string(e)
}

func (e *SearchOutStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SearchOutStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SearchOutStatus", str)
	}
	return nil
}

func (e SearchOutStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SigninOutStatus string

const (
	SigninOutStatusOk                  SigninOutStatus = "OK"
	SigninOutStatusBadRequest          SigninOutStatus = "BAD_REQUEST"
	SigninOutStatusServerInternalError SigninOutStatus = "SERVER_INTERNAL_ERROR"
)

var AllSigninOutStatus = []SigninOutStatus{
	SigninOutStatusOk,
	SigninOutStatusBadRequest,
	SigninOutStatusServerInternalError,
}

func (e SigninOutStatus) IsValid() bool {
	switch e {
	case SigninOutStatusOk, SigninOutStatusBadRequest, SigninOutStatusServerInternalError:
		return true
	}
	return false
}

func (e SigninOutStatus) String() string {
	return string(e)
}

func (e *SigninOutStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SigninOutStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SigninOutStatus", str)
	}
	return nil
}

func (e SigninOutStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SignupOutStatus string

const (
	SignupOutStatusOk                  SignupOutStatus = "OK"
	SignupOutStatusBadRequest          SignupOutStatus = "BAD_REQUEST"
	SignupOutStatusServerInternalError SignupOutStatus = "SERVER_INTERNAL_ERROR"
	SignupOutStatusUserExists          SignupOutStatus = "USER_EXISTS"
)

var AllSignupOutStatus = []SignupOutStatus{
	SignupOutStatusOk,
	SignupOutStatusBadRequest,
	SignupOutStatusServerInternalError,
	SignupOutStatusUserExists,
}

func (e SignupOutStatus) IsValid() bool {
	switch e {
	case SignupOutStatusOk, SignupOutStatusBadRequest, SignupOutStatusServerInternalError, SignupOutStatusUserExists:
		return true
	}
	return false
}

func (e SignupOutStatus) String() string {
	return string(e)
}

func (e *SignupOutStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SignupOutStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SignupOutStatus", str)
	}
	return nil
}

func (e SignupOutStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
