package snype

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/aptos-monkey-monitor/pkg/request"
)

type MonkeyEntityBo struct {
	Name  string `json:"name"`
	Owner struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	} `json:"owner"`
	Image   string  `json:"image"`
	Rank    int     `json:"rank"`
	Royalty float64 `json:"royalty"`
	Listing struct {
		Marketplace struct {
			Name string  `json:"name"`
			Fee  float64 `json:"fee"`
		} `json:"marketplace"`
		Price  int `json:"price"`
		Seller struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"seller"`
	} `json:"listing"`
	LastPrice struct {
		Timestamp   int64 `json:"timestamp"`
		Marketplace struct {
			Name string  `json:"name"`
			Fee  float64 `json:"fee"`
		} `json:"marketplace"`
		Price struct {
			Value int     `json:"value"`
			Usd   float64 `json:"usd"`
		} `json:"price"`
		Buyer struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"buyer"`
	} `json:"lastPrice"`
}

type MonkeyCollectionBo struct {
	Collection string `json:"collection"`
	Creator    struct {
		Address string `json:"address"`
	} `json:"creator"`
	Description string  `json:"description"`
	Items       int     `json:"items"`
	Owners      int     `json:"owners"`
	Listed      int     `json:"listed"`
	Staked      int     `json:"staked"`
	Floor       int     `json:"floor"`
	Volume      []any   `json:"volume"`
	Sales       []int   `json:"sales"`
	Royalty     float64 `json:"royalty"`
	Ranked      bool    `json:"ranked"`
	Image       string  `json:"image"`
	Website     string  `json:"website"`
	Twitter     string  `json:"twitter"`
	Followers   int     `json:"followers"`
	Discord     string  `json:"discord"`
}

type MonkeyBo struct {
	Collection MonkeyCollectionBo `json:"collection"`
	Items      []MonkeyEntityBo   `json:"items"`
	Length     int64              `json:"length"`
	Price      float64            `json:"price"`
}

type MonkeyLowPriceBo struct {
	Name     string
	Price    float64
	UsdPrice float64
	ImageUrl string
	ListUrl  string
}

const monkeyAPI = "https://api.spookslabs.com/v1/collection/aptos-monkeys-mdbl3e/items?page=0&sort=price&collection"
const monkeyImageAPI = "https://ipfs.io/ipfs/bafybeig6bepf5ci5fyysxlfefpjzwkfp7sarj6ed2f5a34kowgc6qenjfa"
const monkeyTopazAPI = "https://www.topaz.so/assets/Aptos-Monkeys-f932dcb983"

func GetMonkeyData() (*MonkeyLowPriceBo, error) {
	sReq := request.New(monkeyAPI, http.MethodGet, nil)
	body, err := sReq.Get(nil)
	if err != nil {
		return nil, err
	}

	monkeyData, err := format(body)
	if err != nil {
		return nil, err
	}

	monkey := monkeyData.Items[0]
	monkeyNum := strings.Split(monkey.Name, "#")[1]

	price := float64(monkey.Listing.Price) / math.Pow(10, 8)
	usdPrice, err := strconv.ParseFloat(fmt.Sprintf("%.2f", price*monkeyData.Price), 64)
	if err != nil {
		return nil, err
	}

	listurl := ""
	if monkey.Listing.Marketplace.Name == "topaz" {
		listurl, err = url.JoinPath(monkeyTopazAPI, monkey.Name, "/0")
		if err != nil {
			return nil, err
		}
	}

	tmpData := &MonkeyLowPriceBo{
		Name:     monkey.Name,
		Price:    price,
		UsdPrice: usdPrice,
		ImageUrl: fmt.Sprintf("%s/%s.png", monkeyImageAPI, monkeyNum),
		ListUrl:  listurl,
	}

	return tmpData, nil
}

func format(body []byte) (*MonkeyBo, error) {
	data := &MonkeyBo{}
	if err := json.Unmarshal(body, data); err != nil {
		return nil, err
	}
	return data, nil
}
