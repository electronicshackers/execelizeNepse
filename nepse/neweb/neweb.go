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
