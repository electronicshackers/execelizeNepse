package neweb

import (
	"fmt"
	"nepse-backend/utils"
	"net/http"
	"os"
)

const (
	Health       = "web/menu"
	All          = "nots/securityDailyTradeStat/58"
	PriceHistory = "nots/market/security/price"
	MutualFund   = "nots/securityDailyTradeStat/66"
	Floorsheet   = "nots/security/floorsheet"
)

type NewebAPI struct {
	client *utils.Client
}

func Neweb() (*NewebAPI, error) {
	client := utils.NewClient(nil, os.Getenv("NEPSE"))

	_, err := client.NewRequest(http.MethodGet, Health, nil)

	if err != nil {
		return nil, err
	}

	biz := &NewebAPI{
		client: client,
	}
	return biz, nil
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
