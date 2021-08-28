package smartwealth

import (
	"fmt"
	"nepse-backend/utils"
	"net/http"
	"os"
)

const (
	Technical = "GetTechnicalScreenerData"
)

type SmartWealthApi struct {
	client *utils.Client
}

func SmartWealth() (*SmartWealthApi, error) {
	client := utils.NewClient(nil, os.Getenv("SMART_WEALTH"))

	_, err := client.NewRequest(http.MethodGet, Technical, nil)

	if err != nil {
		return nil, err
	}

	biz := &SmartWealthApi{
		client: client,
	}
	return biz, nil
}

func (b *SmartWealthApi) buildTechnicalSlug(urlPath string, sectorId int) string {
	return fmt.Sprintf("%s?sectorId=%d&period=YTD&nature=company&dataCount=-1", urlPath, sectorId)
}
