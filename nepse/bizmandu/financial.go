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
			Ticker                                 string  `json:"Ticker"`
			Year                                   string  `json:"Year"`
			Quarter                                float64 `json:"Quarter"`
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
			Totalpremium                           int     `json:"TotalPremium"`              // Non-life ko chij ho yo
			Totalnoofpolicies                      int     `json:"TotalNoOfPolicies"`         // Non-life ko chij ho yo
			Totalrenewedpolicies                   int     `json:"TotalRenewedPolicies"`      // Non-life ko chij ho yo
			Totalclaimamount                       int     `json:"TotalClaimAmount"`          // Non-life ko chij ho yo
			Totalnoofclaims                        int     `json:"TotalNoOfClaims"`           // Non-life ko chij ho yo
			Provisionforunexpiredrisk              int     `json:"ProvisionForUnexpiredRisk"` // Non-life ko chij ho yo
			Totalclaimunexpiredrisk                int     `json:"TotalClaimUnexpiredRisk"`   // Non-life ko chij ho yo
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
			Ticker                                                 string  `json:"Ticker"`
			Ppe                                                    float64 `json:"Ppe"`
			Noncurrentassets                                       int     `json:"NonCurrentAssets"`
			Year                                                   string  `json:"Year"`
			Quarter                                                float64 `json:"Quarter"`
			Totalassets                                            float64 `json:"TotalAssets"`
			Totalliabilities                                       float64 `json:"TotalLiabilities"`
			Totalequity                                            float64 `json:"TotalEquity"`
			Totalcapital                                           float64 `json:"TotalCapital"`
			Insurancefund                                          int     `json:"InsuranceFund"`
			Callinadvance                                          int     `json:"CallinAdvance"`
			Premium                                                int     `json:"Premium"`
			Ltliabilities                                          float64 `json:"LtLiabilities"`
			Longtermliabilities                                    int     `json:"LongTermLiabilities"`
			Funddeposits                                           int     `json:"FundDeposits"`
			Loans                                                  int     `json:"Loans"`
			Totalsourcesoffunds                                    int     `json:"TotalSourcesofFunds"`
			Accumulateddepreciation                                int     `json:"AccumulatedDepreciation"`
			Totallongtermliabilities                               int     `json:"TotalLongTermLiabilities"`
			Curliabiandprovision                                   int     `json:"CurLiabiAndProvision"`
			Tradeandotherpayables                                  int     `json:"TradeandOtherPayables"`
			Dividendpayable                                        int     `json:"DividendPayable"`
			Otherliabilitiesandprovisions                          int     `json:"OtherLiabilitiesandProvisions"`
			Corporateincometaxliabilities                          int     `json:"CorporateIncomeTaxLiabilities"`
			Provisionsforpossiblelosses                            int     `json:"ProvisionsforPossibleLosses"`
			Sourcesoffunds                                         float64 `json:"SourcesOfFunds"`
			Lifeinsurancefund                                      float64 `json:"LifeInsuranceFund"`
			Catastrophereserve                                     float64 `json:"CatastropheReserve"`
			Longtermloanborrowings                                 int     `json:"LongTermLoanBorrowings"`
			Usesoffunds                                            float64 `json:"UsesOfFunds"`
			Fixedassetsnet                                         float64 `json:"FixedAssetsNet"`
			Longterminvestmentsloan                                float64 `json:"LongTermInvestmentsLoan"`
			Expenditureextentnotwrittenoff                         int     `json:"ExpenditureExtentNotWrittenoff"`
			Miscellaneousexpensestotheextentnotwrittenoff          float64 `json:"MiscellaneousExpensesToTheExtentNotWrittenOff"`
			Losstransferredfromprofitlossaccount                   int     `json:"LossTransferredFromProfitLossAccount"`
			Detailsofnetcurrentassets                              int     `json:"DetailsOfNetCurrentAssets"`
			Currentassetsloansandadvances                          float64 `json:"CurrentAssetsLoansAndAdvances"`
			Shortterminvestmentsandloans                           float64 `json:"ShortTermInvestmentsAndLoans"`
			Currentliabilitiesandprovisions                        float64 `json:"CurrentLiabilitiesAndProvisions"`
			Currentliabilities                                     float64 `json:"CurrentLiabilities"`
			Provisionforunexpiredrisk                              float64 `json:"ProvisionForUnexpiredRisk"`
			Provisionforoutstandingclaims                          float64 `json:"ProvisionForOutstandingClaims"`
			Otherprovision                                         float64 `json:"OtherProvision"`
			Totalfunds                                             float64 `json:"TotalFunds"`
			Fixedassets                                            int     `json:"FixedAssets"`
			Depreciation                                           int     `json:"Depreciation"`
			Workinprogress                                         float64 `json:"WorkInProgress"`
			Cash                                                   float64 `json:"Cash"`
			Receivables                                            int     `json:"Receivables"`
			Advancesprepaymentsloansdeposits                       float64 `json:"AdvancesPrepaymentsLoansDeposits"`
			Inventory                                              float64 `json:"Inventory"`
			Totalcurrentassets                                     float64 `json:"TotalCurrentAssets"`
			Stliabilities                                          float64 `json:"StLiabilities"`
			Deferredliabilities                                    int     `json:"DeferredLiabilities"`
			Totalstliabilities                                     float64 `json:"TotalStLiabilities"`
			Applicationoffunds                                     float64 `json:"ApplicationOfFunds"`
			Capitalliabilities                                     int     `json:"CapitalLiabilities"`
			Capitalreserves                                        int     `json:"CapitalReserves"`
			Sharecapital                                           int     `json:"ShareCapital"`
			Reservesurplus                                         float64 `json:"ReserveSurplus"`
			Mediumlongtermloans                                    int     `json:"MediumLongTermLoans"`
			Securedloan                                            int     `json:"SecuredLoan"`
			Unsecuredloan                                          int     `json:"UnsecuredLoan"`
			Deferredtaxliability                                   float64 `json:"DeferredTaxLiability"`
			Grandtotal                                             float64 `json:"GrandTotal"`
			Assets                                                 int     `json:"Assets"`
			Netfixedassets                                         float64 `json:"NetFixedAssets"`
			Capitalworkinprogress                                  float64 `json:"CapitalWorkInProgress"`
			Noncoreassets                                          int     `json:"NonCoreAssets"`
			Cashandbankbalance                                     float64 `json:"CashAndBankBalance"`
			Duefromnrb                                             float64 `json:"DueFromNrb"`
			Placementwithbfis                                      float64 `json:"PlacementwithBFIs"`
			Loanandadvances                                        float64 `json:"LoanAndAdvances"`
			Loanandadvancestobfis                                  float64 `json:"LoanandadvancestoBFIs"`
			Loansandadvancestocustomers                            float64 `json:"Loansandadvancestocustomers"`
			Investments                                            float64 `json:"Investments"`
			Investmentsecurities                                   float64 `json:"InvestmentSecurities"`
			Investmentinsubsidiaries                               float64 `json:"Investmentinsubsidiaries"`
			Investmentinassociates                                 float64 `json:"Investmentinassociates"`
			Investmentproperty                                     float64 `json:"Investmentproperty"`
			Otherassetstot                                         float64 `json:"OtherAssetsTot"`
			Derivativefinancialinstruments                         float64 `json:"DerivativeFinancialInstruments"`
			Othertradingassets                                     float64 `json:"OtherTradingAssets"`
			Currenttaxassets                                       float64 `json:"CurrentTaxassets"`
			Otherassets                                            float64 `json:"OtherAssets"`
			Goodwill                                               float64 `json:"Goodwill"`
			Propequip                                              float64 `json:"PropEquip"`
			Deferredtaxassets                                      int     `json:"DeferredTaxAssets"`
			Totalnoncurrentassets                                  float64 `json:"TotalNonCurrentAssets"`
			Currentassetsloansadvances                             int     `json:"CurrentAssetsLoansAdvances"`
			Investmentsc                                           float64 `json:"InvestmentsC"`
			Inventories                                            float64 `json:"Inventories"`
			Sundrydebtors                                          float64 `json:"SundryDebtors"`
			Cashandbank                                            float64 `json:"CashandBank"`
			Liabilities                                            float64 `json:"Liabilities"`
			Borrowingstot                                          float64 `json:"BorrowingsTot"`
			Duetobankandfinancialinstitutions                      float64 `json:"DuetoBankandFinancialInstitutions"`
			Duetonrb                                               float64 `json:"DuetoNRB"`
			Borrowings                                             float64 `json:"Borrowings"`
			Debtsecuritiesissued                                   float64 `json:"DebtSecuritiesIssued"`
			Deposits                                               float64 `json:"Deposits"`
			Othliabilitiesprov                                     float64 `json:"OthLiabilitiesProv"`
			Derfinancialinstruments                                float64 `json:"DerFinancialInstruments"`
			Provisions                                             float64 `json:"Provisions"`
			Currenttaxliabilities                                  float64 `json:"CurrentTaxLiabilities"`
			Deferredtaxliabilities                                 float64 `json:"DeferredTaxLiabilities"`
			Subordinatedliabilites                                 float64 `json:"SubordinatedLiabilites"`
			Otherliabilities                                       float64 `json:"OtherLiabilities"`
			Equity                                                 float64 `json:"Equity"`
			Paidupcapital                                          float64 `json:"PaidUpCapital"`
			Reservesandsurplus                                     float64 `json:"ReservesAndSurplus"`
			Sharepremium                                           float64 `json:"SharePremium"`
			Retainedearnings                                       float64 `json:"RetainedEarnings"`
			Reserves                                               float64 `json:"Reserves"`
			Noncontrollingfloat64erest                             float64 `json:"NonControllingfloat64erest"`
			Totalcapitalandliabilities                             float64 `json:"TotalCapitalAndLiabilities"`
			Ordinaryshares                                         float64 `json:"OrdinaryShares"`
			Prefshares                                             float64 `json:"PrefShares"`
			Callaccounttimedeposits                                float64 `json:"CallAccountTimeDeposits"`
			Loansadvances                                          float64 `json:"LoansAdvances"`
			Lesscurrentliabilitiesprovisions                       int     `json:"LessCurrentLiabilitiesProvisions"`
			Tradeotherpayables                                     float64 `json:"TradeOtherPayables"`
			Proposeddividends                                      int     `json:"ProposedDividends"`
			Provisionsa                                            float64 `json:"ProvisionsA"`
			Shorttermloans                                         float64 `json:"ShortTermLoans"`
			Totalcurrentliabilities                                float64 `json:"TotalCurrentLiabilities"`
			Netcurrentassets                                       float64 `json:"NetCurrentAssets"`
			Deferredrevenueexpenditure                             int     `json:"DeferredRevenueExpenditure"`
			Lossasprofitlossaccount                                int     `json:"LossasProfitLossAccount"`
			Grandtotala                                            float64 `json:"GrandTotalA"`
			Loansextended                                          int     `json:"LoansExtended"`
			Cashandcashequivalents                                 int     `json:"CashandCashEquivalents"`
			Termdeposits                                           int     `json:"TermDeposits"`
			Othercurrentassets                                     int     `json:"OtherCurrentAssets"`
			Miscassets                                             int     `json:"MiscAssets"`
			Totalusesoffunds                                       int     `json:"TotalUsesofFunds"`
			Preferredshares                                        int     `json:"PreferredShares"`
			Intangibleassets                                       float64 `json:"IntangibleAssets"`
			Financialassets                                        float64 `json:"FinancialAssets"`
			Defferedtaxassets                                      int     `json:"DefferedTaxAssets"`
			Currentassets                                          int     `json:"CurrentAssets"`
			Capitalworkprogress                                    int     `json:"CapitalWorkProgress"`
			Tradereceivables                                       float64 `json:"TradeReceivables"`
			Prepaidreceivables                                     float64 `json:"PrepaidReceivables"`
			Cashequivalents                                        float64 `json:"CashEquivalents"`
			Investmentfd                                           float64 `json:"InvestmentFd"`
			Equityliabilities                                      int     `json:"EquityLiabilities"`
			Noncurrentliabilities                                  int     `json:"NonCurrentLiabilities"`
			Defferedtaxliabilities                                 float64 `json:"DefferedTaxLiabilities"`
			Ncprovisions                                           int     `json:"NcProvisions"`
			Totalnoncurrentliabilities                             float64 `json:"TotalNonCurrentLiabilities"`
			Tradepaybales                                          float64 `json:"TradePaybales"`
			Incometaxliability                                     float64 `json:"IncomeTaxLiability"`
			Miscellaneousexpenditure                               float64 `json:"MiscellaneousExpenditure"`
			Totaleqli                                              float64 `json:"TotalEqLi"`
			Debenturesandbonds                                     int     `json:"DebenturesAndBonds"`
			Domesticcurrency                                       int     `json:"DomesticCurrency"`
			Foreigncurrency                                        int     `json:"ForeignCurrency"`
			Incometaxliabilities                                   int     `json:"IncomeTaxLiabilities"`
			Realestateloan                                         int     `json:"RealEstateLoan,omitempty"`
			Residentialrealestateloanabove10Million                int     `json:"ResidentialRealEstateLoanAbove10Million,omitempty"`
			Businesscomplexandresidentialapartmentconstructionloan int     `json:"BusinessComplexAndResidentialApartmentConstructionLoan,omitempty"`
			Incomegeneratingcommercialcomplexloan                  int     `json:"IncomeGeneratingCommercialComplexLoan,omitempty"`
			Otherrealestateloan                                    int     `json:"OtherRealEstateLoan,omitempty"`
			Personalhomeloanupto10Millionorless                    int     `json:"PersonalHomeLoanUpto10MillionOrLess,omitempty"`
			Margintypeloan                                         int     `json:"MarginTypeLoan,omitempty"`
			Termloan                                               int     `json:"TermLoan,omitempty"`
			Overdraftloanandtrloanandwcloan                        int     `json:"OverdraftLoanAndTrLoanAndWcLoan,omitempty"`
			Otherloan                                              int     `json:"OtherLoan"`
			Nonbankingassets                                       int     `json:"NonBankingAssets"`
			Moneyatcallandshortnotice                              float64 `json:"MoneyAtCallAndShortNotice"`
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
