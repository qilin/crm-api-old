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
	Signin  *SigninOut  `json:"signin"`
	Me      *User       `json:"me"`
	Signout *SignoutOut `json:"signout"`
}

type Covers struct {
	FloorSmall    *Image `json:"floor_small"`
	MainLittle    *Image `json:"main_little"`
	MainBig       *Image `json:"main_big"`
	FloorMedium   *Image `json:"floor_medium"`
	FloorWide     *Image `json:"floor_wide"`
	FloorLarge    *Image `json:"floor_large"`
	FloorSmallest *Image `json:"floor_smallest"`
	FloorWidest   *Image `json:"floor_widest"`
	BackgroundBig *Image `json:"background_big"`
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

type Game struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Description string     `json:"description"`
	Publisher   *Publisher `json:"publisher"`
	Covers      *Covers    `json:"covers"`
	Screenshots []*Image   `json:"screenshots"`
	Tags        []*Tag     `json:"tags"`
	Genre       Genres     `json:"genre"`
	Rating      int        `json:"rating"`
}

type Image struct {
	URL string `json:"url"`
}

type Publisher struct {
	Title string `json:"title"`
}

type SigninOut struct {
	Status SigninOutStatus `json:"status"`
	Token  string          `json:"token"`
}

type SignoutOut struct {
	Status AuthenticatedRequestStatus `json:"status"`
}

type SignupOut struct {
	Status SignupOutStatus `json:"status"`
}

type StoreQuery struct {
	Games []*Game `json:"games"`
}

type Tag struct {
	Name *string  `json:"name"`
	Type *TagType `json:"type"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
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

type Genres string

const (
	GenresBoard     Genres = "Board"
	GenresCards     Genres = "Cards"
	GenresCasino    Genres = "Casino"
	GenresFarm      Genres = "Farm"
	GenresRacing    Genres = "Racing"
	GenresShooter   Genres = "Shooter"
	GenresFindItems Genres = "FindItems"
	GenresPuzzle    Genres = "Puzzle"
	GenresRpg       Genres = "RPG"
	GenresSimulator Genres = "Simulator"
	GenresStrategy  Genres = "Strategy"
)

var AllGenres = []Genres{
	GenresBoard,
	GenresCards,
	GenresCasino,
	GenresFarm,
	GenresRacing,
	GenresShooter,
	GenresFindItems,
	GenresPuzzle,
	GenresRpg,
	GenresSimulator,
	GenresStrategy,
}

func (e Genres) IsValid() bool {
	switch e {
	case GenresBoard, GenresCards, GenresCasino, GenresFarm, GenresRacing, GenresShooter, GenresFindItems, GenresPuzzle, GenresRpg, GenresSimulator, GenresStrategy:
		return true
	}
	return false
}

func (e Genres) String() string {
	return string(e)
}

func (e *Genres) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Genres(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Genres", str)
	}
	return nil
}

func (e Genres) MarshalGQL(w io.Writer) {
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

type TagType string

const (
	TagTypeGenre  TagType = "genre"
	TagTypeCommon TagType = "common"
)

var AllTagType = []TagType{
	TagTypeGenre,
	TagTypeCommon,
}

func (e TagType) IsValid() bool {
	switch e {
	case TagTypeGenre, TagTypeCommon:
		return true
	}
	return false
}

func (e TagType) String() string {
	return string(e)
}

func (e *TagType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TagType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TagType", str)
	}
	return nil
}

func (e TagType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
