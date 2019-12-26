// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphql

import (
	"fmt"
	"io"
	"strconv"

	"github.com/qilin/crm-api/internal/db/domain/store"
)

type AuthMutation struct {
	PasswordUpdate *PasswordUpdateResponse `json:"passwordUpdate"`
}

type AuthQuery struct {
	Profile *User `json:"profile"`
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

type FriendGame struct {
	Game    store.Game `json:"game"`
	Friends []*User    `json:"friends"`
}

type PasswordUpdateResponse struct {
	Status AuthenticatedRequestStatus `json:"status"`
}

type StoreQuery struct {
	Game       store.Game        `json:"game"`
	GameBySlug store.Game        `json:"gameBySlug"`
	Games      []store.Game      `json:"games"`
	Module     store.Module      `json:"module"`
	StoreFront *store.StoreFront `json:"storeFront"`
}

type User struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
	Language  string `json:"language"`
}

type ViewerQuery struct {
	Games        []store.Game  `json:"games"`
	FriendsGames []*FriendGame `json:"friendsGames"`
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

type RequestStatus string

const (
	RequestStatusOk                  RequestStatus = "OK"
	RequestStatusForbidden           RequestStatus = "FORBIDDEN"
	RequestStatusNotFound            RequestStatus = "NOT_FOUND"
	RequestStatusBadRequest          RequestStatus = "BAD_REQUEST"
	RequestStatusServerInternalError RequestStatus = "SERVER_INTERNAL_ERROR"
)

var AllRequestStatus = []RequestStatus{
	RequestStatusOk,
	RequestStatusForbidden,
	RequestStatusNotFound,
	RequestStatusBadRequest,
	RequestStatusServerInternalError,
}

func (e RequestStatus) IsValid() bool {
	switch e {
	case RequestStatusOk, RequestStatusForbidden, RequestStatusNotFound, RequestStatusBadRequest, RequestStatusServerInternalError:
		return true
	}
	return false
}

func (e RequestStatus) String() string {
	return string(e)
}

func (e *RequestStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RequestStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RequestStatus", str)
	}
	return nil
}

func (e RequestStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type RoleEnum string

const (
	RoleEnumAdmin RoleEnum = "ADMIN"
	RoleEnumUser  RoleEnum = "USER"
)

var AllRoleEnum = []RoleEnum{
	RoleEnumAdmin,
	RoleEnumUser,
}

func (e RoleEnum) IsValid() bool {
	switch e {
	case RoleEnumAdmin, RoleEnumUser:
		return true
	}
	return false
}

func (e RoleEnum) String() string {
	return string(e)
}

func (e *RoleEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RoleEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RoleEnum", str)
	}
	return nil
}

func (e RoleEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
