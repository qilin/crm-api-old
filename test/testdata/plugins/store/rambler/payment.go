package rambler

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/qilin/crm-api/internal/sdk/common"
	"github.com/qilin/crm-api/pkg/qilin"
)

type PaymentConfig struct {
	Secret string // p.config.Billing.Secret
	Qilin  string // p.config.URL.Qilin

	ShopID     int // p.config.Billing.ShopID
	ShowcaseID int // p.config.Billing.ScID
}

// PaymentHandler responsible for integration with rambler payment system
type PaymentHandler struct {
	cfg *PaymentConfig

	orders map[string]CreateOrderRequest // TODO must be in DB
}

func NewPaymentHandler(cfg *PaymentConfig) *PaymentHandler {
	return &PaymentHandler{
		cfg:    cfg,
		orders: make(map[string]CreateOrderRequest),
	}
}

type CreateOrderRequest struct {
	GameId string `json:"game_id" query:"game_id" form:"game_id"`
	ItemId string `json:"item_id" query:"item_id" form:"item_id"`
}

func (h *PaymentHandler) CreateOrder(ctx echo.Context) error {
	var req CreateOrderRequest
	if err := ctx.Bind(&req); err != nil {
		fmt.Println(err)
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	item, err := h.queryItem(h.cfg.Qilin, req.GameId, req.ItemId)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	price, err := strconv.ParseFloat(item.Price, 64)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	orderId := fmt.Sprintf("%d", time.Now().Unix()-1574170567) // TODO
	h.orders[orderId] = req

	order := map[string]interface{}{
		"products": []map[string]interface{}{
			{
				// оффер в формате: https://yandex.ru/support/partnermarket/yml/
				"offer": map[string]interface{}{
					"name": item.Title,
					// Описание товара
					"description": "Item 1",
					// Ссылка на товар
					"picture": item.Photo_url,
				},
				// количество покупаемого товара
				"count": 1,
				// итоговая цена
				"total": map[string]interface{}{
					"price":      price,
					"currencyId": "RUB",
				},
			},
		},
		"order": map[string]interface{}{

			"shopId":      h.cfg.ShopID,
			"scId":        h.cfg.ShowcaseID,
			"orderNumber": orderId,
			"orderAmount": price,

			"customerNumber": "1683086",                         //TODO
			"cpsEmail":       "aleksandr.barsukov@protocol.one", // TODO
			"cpsPhone":       "",
			"productName":    item.Title,
			"productImage":   item.Photo_url,
			// "gameTitle":   req.GameId, //TODO

			"requestSign": h.signPaymentRequest(orderId, item.Price),
			"addData": map[string]interface{}{
				"game_id":      req.GameId,
				"game_name":    req.GameId, // TODO
				"product_id":   req.ItemId,
				"product_name": item.Title,
			},
			"orderParams": map[string]interface{}{
				"positions": []map[string]interface{}{
					map[string]interface{}{
						"name":     item.Title,
						"quantity": 1,
						"price":    price,
						"tax":      3,
						"taxValue": price * 20 / 120,
					},
				},
			},
		},
	}
	return ctx.JSON(http.StatusOK, order)
}

func (h *PaymentHandler) CheckOrder(ctx echo.Context) error {
	var req = make(map[string]interface{})
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{})
	}

	fmt.Println("checkOrder:", req)
	fmt.Println("signature header:", ctx.Request().Header.Get("x-rps-signature"))

	sign := paymentSign(
		fmt.Sprint(req["requestAction"]),
		fmt.Sprint(req["operationUid"]),
		fmt.Sprint(req["operationDatetime"]),
		fmt.Sprint(req["paymentType"]),
		fmt.Sprint(req["orderNumber"]),
		fmt.Sprint(req["orderAmount"]),
		fmt.Sprint(req["currencyCode"]),
		h.cfg.Secret,
	)
	fmt.Println("check sign:", sign)

	if req["requestSign"] != sign {
		fmt.Println("invalid sign")
		return ctx.JSON(http.StatusBadRequest, "invalid sign")
	}

	var v = map[string]interface{}{
		"checkResponse": map[string]interface{}{
			"operationUid":      req["operationUid"],
			"operationDatetime": req["operationDatetime"],
			"orderAmount":       req["orderAmount"],
			"code":              0,
		},
	}

	return ctx.JSON(http.StatusOK, v)
}

func (h *PaymentHandler) PaymentAviso(ctx echo.Context) error {
	var req = make(map[string]interface{})
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{})
	}

	fmt.Println("paymentAviso:", req)
	fmt.Println("signature header:", ctx.Request().Header.Get("x-rps-signature"))

	sign := paymentSign(
		fmt.Sprint(req["requestAction"]),
		fmt.Sprint(req["operationUid"]),
		fmt.Sprint(req["operationDatetime"]),
		fmt.Sprint(req["paymentType"]),
		fmt.Sprint(req["orderNumber"]),
		fmt.Sprint(req["orderAmount"]),
		fmt.Sprint(req["currencyCode"]),
		h.cfg.Secret,
	)
	fmt.Println("aviso sign:", sign)
	if req["requestSign"] != sign {
		fmt.Println("invalid sign")
		return ctx.JSON(http.StatusBadRequest, "invalid sign")
	}

	var v = map[string]interface{}{
		"avisoResponse": map[string]interface{}{
			"operationUid":      req["operationUid"],
			"operationDatetime": req["operationDatetime"],
			"orderAmount":       req["orderAmount"],
			"code":              0,
		},
	}
	data := h.orders[req["orderNumber"].(string)]

	orderReq := common.OrderRequest{
		GameID: data.GameId,
		UserID: "123",
		ItemID: data.ItemId,
	}
	if err := h.confirmBuy(h.cfg.Qilin, orderReq); err != nil {
		fmt.Println("payment aviso failed:", err.Error())
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, v)
}

func (h *PaymentHandler) PaymentNotification(ctx echo.Context) error {
	var v = map[string]interface{}{}
	ctx.Bind(&v)
	fmt.Printf("payment notification: %#v\n", v)
	return ctx.JSON(http.StatusOK, map[string]interface{}{})
}

func (h *PaymentHandler) signPaymentRequest(orderId, amount string) string {
	return paymentSign(
		strconv.Itoa(h.cfg.ShopID),
		strconv.Itoa(h.cfg.ShowcaseID),
		orderId,
		amount,
		h.cfg.Secret,
	)
}

func paymentSign(params ...string) string {
	str := strings.Join(params, ";")
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}

type Item struct {
	Title     string `json:"title"`
	Photo_url string `json:"photo_url"`
	Price     string `json:"price"`
}

func (h *PaymentHandler) queryItem(entry, gameId, itemId string) (*Item, error) {
	u := fmt.Sprintf("%s?item_id=%s&game_id=%s", qilin.ItemsURL(entry), itemId, gameId)

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("provider error")
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var item Item
	err = json.Unmarshal(d, &item)
	return &item, err
}

func (h *PaymentHandler) confirmBuy(entry string, req common.OrderRequest) error {
	data, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	resp, err := http.Post(qilin.OrderURL(entry), "application/json;charset=utf-8", bytes.NewReader(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("payment proceeding failed")
	}

	return nil
}
