package neweb

import (
	"context"
	"nepse-backend/nepse"
	"net/http"
)

type test struct {
	Id string `json:"id"`
}

func (n *NewebAPI) GetFloorsheet(stockId, businessDate, randomId string, page, size int, isBulkRequest bool) (*nepse.FloorsheetResponse, error) {
	url := n.buildFloorsheetSlug(stockId, businessDate, page, size)
	ok := test{Id: randomId}

	token, err := n.RecursiveGetToken()
	if err != nil {
		return nil, err
	}

	// n.RecursiveGetToken()

	req, err := n.client.NewRequest(http.MethodPost, url, ok, token)
	if err != nil {
		return nil, err
	}

	res := &nepse.FloorsheetResponse{}

	if _, err := n.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NewebAPI) RecursiveGetToken() (string, error) {
	var isError = false
	headers, err := n.Prove()
	if err != nil {
		return "", err
	}
	token, err := n.client.Wasm(*headers)
	if err != nil {
		return "", err
	}
	if isError {
		n.RecursiveGetToken()
		isError = false
	}

	return token.Accesstoken, nil

}
