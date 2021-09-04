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

	if isBulkRequest {
		headers, err := n.Prove()
		if err != nil {
			return nil, err
		}
		token := n.client.Wasm(*headers)
		n.client.Headers = token.Accesstoken
	}

	req, err := n.client.NewRequest(http.MethodPost, url, ok)
	if err != nil {
		return nil, err
	}

	res := &nepse.FloorsheetResponse{}

	if _, err := n.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}
	return res, nil
}
