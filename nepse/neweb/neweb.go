package neweb

import (
	"context"
	"fmt"
	"nepse-backend/utils"
	"net/http"
	"os"
)

const (
	Health       = "web/menu"
	Prove        = "authenticate/prove"
	All          = "nots/securityDailyTradeStat/58"
	PriceHistory = "nots/market/security/price"
	MutualFund   = "nots/securityDailyTradeStat/66"
	Floorsheet   = "nots/security/floorsheet"
)

type NewebAPI struct {
	client *utils.Client
}

func Neweb() (*NewebAPI, error) {
	var client *utils.Client
	client = utils.NewClient(nil, os.Getenv("NEPSE"), "")
	req, err := client.NewRequest(http.MethodGet, Prove, nil)

	if err != nil {
		return nil, err
	}

	res := &utils.ProveResponse{}

	if _, err := client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}

	prove := client.Wasm(*res)

	client = utils.NewClient(nil, os.Getenv("NEPSE"), prove.Accesstoken)

	newWeb := &NewebAPI{
		client: client,
	}
	return newWeb, nil
}

func (n *NewebAPI) Prove() (*utils.ProveResponse, error) {
	req, err := n.client.NewRequest(http.MethodGet, Prove, nil)

	if err != nil {
		return nil, err
	}

	res := &utils.ProveResponse{}

	if _, err := n.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (b *NewebAPI) buildHistorySlug(urlPath, ticker string) string {
	return fmt.Sprintf("%s/%v?size=65", urlPath, ticker)
}

func (b *NewebAPI) buildFloorsheetSlug(id, date string, page int, size int) string {
	if size != 0 {
		return fmt.Sprintf("%s/%s?page=%d&size=%d&businessDate=%s&sort=contractId,asc", Floorsheet, id, page, size, date)
	}
	return fmt.Sprintf("%s/%s?&businessDate=%s&sort=contractId,asc", Floorsheet, id, date)
}
