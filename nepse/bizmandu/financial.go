package bizmandu

import (
	"context"
	"net/http"
)

type FinancialData struct {
	Response float64 `json:"response"`
	Error    string  `json:"error"`
	Message  struct {
		Ticker string `json:"ticker"`
		Data   []struct {
			Idcbkeystats                           float64 `json:"idcbKeyStats"`
			Ticker                                 string  `json:"Ticker"`
			Year                                   string  `json:"Year"`
			Quarter                                float64 `json:"Quarter"`
			Datasource                             string  `json:"DataSource"`
			Totalrevenue                           float64 `json:"TotalRevenue"`
			Growthoverpriorperiod                  float64 `json:"GrowthOverPriorPeriod"`
			Grossprofit                            float64 `json:"GrossProfit"`
			Grossprofitmargin                      float64 `json:"GrossProfitMargin"`
			Netincome                              float64 `json:"NetIncome"`
			Netincomemargin                        float64 `json:"NetIncomeMargin"`
			Returnonasset                          float64 `json:"ReturnOnAsset"`
			Returnonequity                         float64 `json:"ReturnOnEquity"`
			Epsannualized                          float64 `json:"EpsAnnualized"`
			Reportedpeannualized                   float64 `json:"ReportedPeAnnualized"`
			Bookvaluepershare                      float64 `json:"BookValuePerShare"`
			Dividendpershare                       float64 `json:"DividendPerShare"`
			Capitalfundtorwa                       float64 `json:"CapitalFundToRwa"`
			Nonperformingloannpltototalloan        float64 `json:"NonPerformingLoanNplToTotalLoan"`
			Totalloanlossprovisiontototalnpl       float64 `json:"TotalLoanLossProvisionToTotalNpl"`
			Costoffunds                            float64 `json:"CostOfFunds"`
			Creditdepositratioaspernrbcalculations float64 `json:"CreditDepositRatioAsPerNrbCalculations"`
			Baserate                               float64 `json:"BaseRate"`
			NetInterestspread                      float64 `json:"NetInterestSpread"`
			Averageyield                           float64 `json:"AverageYield"`
			Outstandingshares                      float64 `json:"OutstandingShares,omitempty"`
			Tier1Capital                           float64 `json:"Tier1Capital,omitempty"`
			Tier2Capital                           float64 `json:"Tier2Capital,omitempty"`
			Totalcapital                           float64 `json:"TotalCapital,omitempty"`
			Shareprice                             float64 `json:"SharePrice,omitempty"`
			Marketcapitalization                   float64 `json:"MarketCapitalization,omitempty"`
			Mps                                    float64 `json:"Mps"`
			Netliquidasset                         float64 `json:"NetLiquidAsset"`
		} `json:"data"`
	} `json:"message"`
}

type BalanceSheetData struct {
	Response float64 `json:"response"`
	Error    string  `json:"error"`
	Message  struct {
		Ticker string `json:"ticker"`
		Data   []struct {
			Idcbbalancesheet2                 float64 `json:"idcbbalancesheet2"`
			Ticker                            string  `json:"Ticker"`
			Year                              string  `json:"Year"`
			Quarter                           float64 `json:"Quarter"`
			Datasource                        string  `json:"DataSource"`
			Statement                         string  `json:"Statement"`
			Totalassets                       float64 `json:"TotalAssets"`
			Cashandbankbalance                float64 `json:"CashAndBankBalance"`
			Moneyatcallandshortnotice         float64 `json:"MoneyAtCallAndShortNotice"`
			Duefromnrb                        float64 `json:"DueFromNrb"`
			Placementwithbfis                 float64 `json:"PlacementwithBFIs"`
			Loanandadvances                   float64 `json:"LoanAndAdvances"`
			Loanandadvancestobfis             float64 `json:"LoanandadvancestoBFIs"`
			Loansandadvancestocustomers       float64 `json:"Loansandadvancestocustomers"`
			Investments                       float64 `json:"Investments"`
			Investmentsecurities              float64 `json:"InvestmentSecurities"`
			Investmentinsubsidiaries          float64 `json:"Investmentinsubsidiaries"`
			Investmentinassociates            float64 `json:"Investmentinassociates"`
			Investmentproperty                float64 `json:"Investmentproperty"`
			Otherassetstot                    float64 `json:"OtherAssetsTot"`
			Derivativefinancialinstruments    float64 `json:"DerivativeFinancialInstruments"`
			Othertradingassets                float64 `json:"OtherTradingAssets"`
			Currenttaxassets                  float64 `json:"CurrentTaxassets"`
			Deferredtaxassets                 float64 `json:"DeferredTaxAssets"`
			Otherassets                       float64 `json:"OtherAssets"`
			Goodwill                          float64 `json:"Goodwill"`
			Propequip                         float64 `json:"PropEquip"`
			Liabilities                       float64 `json:"Liabilities"`
			Borrowingstot                     float64 `json:"BorrowingsTot"`
			Duetobankandfinancialinstitutions float64 `json:"DuetoBankandFinancialInstitutions"`
			Duetonrb                          float64 `json:"DuetoNRB"`
			Borrowings                        float64 `json:"Borrowings"`
			Debtsecuritiesissued              float64 `json:"DebtSecuritiesIssued"`
			Deposits                          float64 `json:"Deposits"`
			Othliabilitiesprov                float64 `json:"OthLiabilitiesProv"`
			Derfinancialinstruments           float64 `json:"DerFinancialInstruments"`
			Provisions                        float64 `json:"Provisions"`
			Currenttaxliabilities             float64 `json:"CurrentTaxLiabilities"`
			Deferredtaxliabilities            float64 `json:"DeferredTaxLiabilities"`
			Subordinatedliabilites            float64 `json:"SubordinatedLiabilites"`
			Otherliabilities                  float64 `json:"OtherLiabilities"`
			Equity                            float64 `json:"Equity"`
			Paidupcapital                     float64 `json:"PaidUpCapital"`
			Reservesandsurplus                float64 `json:"ReservesAndSurplus"`
			Sharepremium                      float64 `json:"SharePremium"`
			Retainedearnings                  float64 `json:"RetainedEarnings"`
			Reserves                          float64 `json:"Reserves"`
			Noncontrollingfloat64erest        float64 `json:"NonControllingfloat64erest"`
			Totalcapitalandliabilities        float64 `json:"TotalCapitalAndLiabilities"`
			Ordinaryshares                    float64 `json:"OrdinaryShares"`
			Prefshares                        float64 `json:"PrefShares"`
		} `json:"data"`
	} `json:"message"`
}

type IncomeStatementData struct {
	Response float64 `json:"response"`
	Error    string  `json:"error"`
	Message  struct {
		Ticker string `json:"ticker"`
		Data   []struct {
			Idcbincomestatement                              float64 `json:"idcbincomestatement"`
			Ticker                                           string  `json:"Ticker"`
			Year                                             string  `json:"Year"`
			Quarter                                          float64 `json:"Quarter"`
			Datasource                                       string  `json:"DataSource"`
			Statement                                        string  `json:"Statement"`
			Interestincome                                   float64 `json:"InterestIncome"`
			Interestexpense                                  float64 `json:"InterestExpense"`
			NetInterestincome                                float64 `json:"NetInterestIncome"`
			Feesincome                                       float64 `json:"FeesIncome"`
			Feesexpense                                      float64 `json:"FeesExpense"`
			Netfeesincome                                    float64 `json:"NetFeesIncome"`
			NetInterestfeeandcommissionincome                float64 `json:"NetInterestfeeandcommissionincome"`
			Nettradingincome                                 float64 `json:"NetTradingIncome"`
			Otheroperatingincome                             float64 `json:"OtherOperatingIncome"`
			Totaloperatingincome                             float64 `json:"TotalOperatingIncome"`
			Impairment                                       float64 `json:"Impairment"`
			Netopincome                                      float64 `json:"NetOpIncome"`
			Operatingexpense                                 float64 `json:"Operatingexpense"`
			Staffexpenses                                    float64 `json:"StaffExpenses"`
			Otheroperatingexpenses                           float64 `json:"OtherOperatingExpenses"`
			Depreciationandamortization                      float64 `json:"DepreciationandAmortization"`
			Operatingprofit                                  float64 `json:"OperatingProfit"`
			Nonoperatingincome                               float64 `json:"Nonoperatingincome"`
			Nonoperatingexpense                              float64 `json:"Nonoperatingexpense"`
			Profitbeforetax                                  float64 `json:"ProfitbeforeTax"`
			Incometax                                        float64 `json:"IncomeTax"`
			Currenttax                                       float64 `json:"CurrentTax"`
			Deferredtax                                      float64 `json:"DeferredTax"`
			Netprofitorloss                                  float64 `json:"NetProfitOrLoss"`
			Compincome                                       float64 `json:"CompIncome"`
			Totalcompincome                                  float64 `json:"TotalCompIncome"`
			Netprofitlossasperprofitorloss                   float64 `json:"NetprofitlossasperProfitorLoss"`
			Profitrequiredtobeapropriatedtostatutoryreserve  float64 `json:"ProfitrequiredtobeapropriatedtoStatutoryreserve"`
			Profitrequiredtobetransferredtoregulatoryreserve float64 `json:"ProfitrequiredtobetransferredtoRegulatoryreserve"`
			Freeprofit                                       float64 `json:"FreeProfit"`
			Prefsharediv                                     float64 `json:"PrefShareDiv"`
		} `json:"data"`
	} `json:"message"`
}

func (b *BizmanduAPI) GetFinancial(ticker string) (*FinancialData, error) {
	url := b.buildTickerSlug(Financial, ticker)

	req, err := b.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res := &FinancialData{}

	if _, err := b.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (b *BizmanduAPI) GetBalanceSheet(ticker string) (*BalanceSheetData, error) {
	url := b.buildTickerSlug(BalanceSheet, ticker)

	req, err := b.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res := &BalanceSheetData{}

	if _, err := b.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (b *BizmanduAPI) GetIncomeStatement(ticker string) (*IncomeStatementData, error) {
	url := b.buildTickerSlug(IncomeStatement, ticker)

	req, err := b.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res := &IncomeStatementData{}

	if _, err := b.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}

	return res, nil
}
