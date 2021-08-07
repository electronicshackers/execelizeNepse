package nepse

type Ticker struct {
	Ticker          string  `json:"ticker"`
	Id              string  `json:"id"`
	Companyname     string  `json:"companyName"`
	Sector          string  `json:"sector"`
	Lasttradedprice float64 `json:"lastTradedPrice"`
}

type LastTradingDayStats struct {
	Ticker              string      `json:"ticker"`
	Openprice           float64     `json:"openPrice"`
	Highprice           float64     `json:"highPrice"`
	Lowprice            float64     `json:"lowPrice"`
	PointChanged        float64     `json:"pointChanged"`
	Totaltradequantity  int         `json:"totalTradeQuantity"`
	Lasttradedprice     float64     `json:"lastTradedPrice"`
	Percentagechange    float64     `json:"percentageChange"`
	Lastupdateddatetime string      `json:"lastUpdatedDateTime"`
	Lasttradedvolume    interface{} `json:"lastTradedVolume"`
	Previousclose       float64     `json:"previousClose"`
}

type PriceHistoryMinified struct {
	Date         string  `json:"date"`
	Price        float64 `json:"price"`
	AveragePrice float64 `json:"averagePrice"`
	HighPrice    float64 `json:"highPrice"`
	LowPrice     float64 `json:"lowPrice"`
}

type PriceHistory struct {
	Content []struct {
		ID           int    `json:"id"`
		Businessdate string `json:"businessDate"`
		Security     struct {
			ID               int         `json:"id"`
			Symbol           string      `json:"symbol"`
			Isin             string      `json:"isin"`
			Permittedtotrade string      `json:"permittedToTrade"`
			Listingdate      string      `json:"listingDate"`
			Creditrating     interface{} `json:"creditRating"`
			Ticksize         float64     `json:"tickSize"`
			Instrumenttype   struct {
				ID           int    `json:"id"`
				Code         string `json:"code"`
				Description  string `json:"description"`
				Activestatus string `json:"activeStatus"`
			} `json:"instrumentType"`
			Capitalgainbasedate string      `json:"capitalGainBaseDate"`
			Facevalue           float64     `json:"faceValue"`
			Highrangedpr        interface{} `json:"highRangeDPR"`
			Issuername          interface{} `json:"issuerName"`
			Meinstancenumber    int         `json:"meInstanceNumber"`
			Parentid            interface{} `json:"parentId"`
			Recordtype          int         `json:"recordType"`
			Schemedescription   interface{} `json:"schemeDescription"`
			Schemename          interface{} `json:"schemeName"`
			Secured             interface{} `json:"secured"`
			Series              interface{} `json:"series"`
			Sharegroupid        struct {
				ID              int         `json:"id"`
				Name            string      `json:"name"`
				Description     string      `json:"description"`
				Capitalrangemin int         `json:"capitalRangeMin"`
				Modifiedby      interface{} `json:"modifiedBy"`
				Modifieddate    interface{} `json:"modifiedDate"`
				Activestatus    string      `json:"activeStatus"`
				Isdefault       string      `json:"isDefault"`
			} `json:"shareGroupId"`
			Activestatus       string  `json:"activeStatus"`
			Divisor            int     `json:"divisor"`
			Cdsstockrefid      int     `json:"cdsStockRefId"`
			Securityname       string  `json:"securityName"`
			Tradingstartdate   string  `json:"tradingStartDate"`
			Networthbaseprice  float64 `json:"networthBasePrice"`
			Securitytradecycle int     `json:"securityTradeCycle"`
			Ispromoter         string  `json:"isPromoter"`
			Companyid          struct {
				ID                   int         `json:"id"`
				Companyshortname     string      `json:"companyShortName"`
				Companyname          string      `json:"companyName"`
				Email                string      `json:"email"`
				Companywebsite       interface{} `json:"companyWebsite"`
				Companycontactperson string      `json:"companyContactPerson"`
				Sectormaster         struct {
					ID                int    `json:"id"`
					Sectordescription string `json:"sectorDescription"`
					Activestatus      string `json:"activeStatus"`
					Regulatorybody    string `json:"regulatoryBody"`
				} `json:"sectorMaster"`
				Companyregistrationnumber interface{} `json:"companyRegistrationNumber"`
				Activestatus              string      `json:"activeStatus"`
			} `json:"companyId"`
		} `json:"security"`
		Openprice             float64 `json:"openPrice"`
		Highprice             float64 `json:"highPrice"`
		Lowprice              float64 `json:"lowPrice"`
		Closeprice            float64 `json:"closePrice"`
		Totaltradedquantity   float64 `json:"totalTradedQuantity"`
		Totaltradedvalue      float64 `json:"totalTradedValue"`
		Previousdaycloseprice float64 `json:"previousDayClosePrice"`
		Fiftytwoweekhigh      float64 `json:"fiftyTwoWeekHigh"`
		Fiftytwoweeklow       float64 `json:"fiftyTwoWeekLow"`
		Lastupdatedtime       string  `json:"lastUpdatedTime"`
		Totaltrades           int     `json:"totalTrades"`
		Lasttradedprice       float64 `json:"lastTradedPrice"`
		Averagetradedprice    float64 `json:"averageTradedPrice"`
	} `json:"content"`
	Pageable struct {
		Sort struct {
			Sorted   bool `json:"sorted"`
			Unsorted bool `json:"unsorted"`
			Empty    bool `json:"empty"`
		} `json:"sort"`
		Pagesize   int  `json:"pageSize"`
		Pagenumber int  `json:"pageNumber"`
		Offset     int  `json:"offset"`
		Paged      bool `json:"paged"`
		Unpaged    bool `json:"unpaged"`
	} `json:"pageable"`
	Totalpages       int  `json:"totalPages"`
	Totalelements    int  `json:"totalElements"`
	Last             bool `json:"last"`
	Number           int  `json:"number"`
	Size             int  `json:"size"`
	Numberofelements int  `json:"numberOfElements"`
	Sort             struct {
		Sorted   bool `json:"sorted"`
		Unsorted bool `json:"unsorted"`
		Empty    bool `json:"empty"`
	} `json:"sort"`
	First bool `json:"first"`
	Empty bool `json:"empty"`
}

type NepseInterface interface {
	GetStocks() ([]Ticker, error)
	GetSectors()
	GetCurrentPrice(ticker string) (*LastTradingDayStats, error)
	GetIncomeStatement(ticker string)
	GetBalanceSheets(ticker string)
	GetFinancialDetails(ticker string)
	GetStockDetails(ticker string)
	GetSummary(ticker string)
	GetPriceHistory(ticker string) ([]PriceHistoryMinified, error)
}
