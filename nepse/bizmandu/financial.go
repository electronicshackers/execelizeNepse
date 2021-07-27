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
			Totalpremium                           float64 `json:"TotalPremium"`              // Non-life ko chij ho yo
			Totalnoofpolicies                      float64 `json:"TotalNoOfPolicies"`         // Non-life ko chij ho yo
			Totalrenewedpolicies                   float64 `json:"TotalRenewedPolicies"`      // Non-life ko chij ho yo
			Totalclaimamount                       float64 `json:"TotalClaimAmount"`          // Non-life ko chij ho yo
			Totalnoofclaims                        float64 `json:"TotalNoOfClaims"`           // Non-life ko chij ho yo
			Provisionforunexpiredrisk              float64 `json:"ProvisionForUnexpiredRisk"` // Non-life ko chij ho yo
			Totalclaimunexpiredrisk                float64 `json:"TotalClaimUnexpiredRisk"`   // Non-life ko chij ho yo
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
			Noncurrentassets                                       float64 `json:"NonCurrentAssets"`
			Year                                                   string  `json:"Year"`
			Quarter                                                float64 `json:"Quarter"`
			Totalassets                                            float64 `json:"TotalAssets"`
			Totalliabilities                                       float64 `json:"TotalLiabilities"`
			Totalequity                                            float64 `json:"TotalEquity"`
			Totalcapital                                           float64 `json:"TotalCapital"`
			Insurancefund                                          float64 `json:"InsuranceFund"`
			Callinadvance                                          float64 `json:"CallinAdvance"`
			Premium                                                float64 `json:"Premium"`
			Ltliabilities                                          float64 `json:"LtLiabilities"`
			Longtermliabilities                                    float64 `json:"LongTermLiabilities"`
			Funddeposits                                           float64 `json:"FundDeposits"`
			Loans                                                  float64 `json:"Loans"`
			Totalsourcesoffunds                                    float64 `json:"TotalSourcesofFunds"`
			Accumulateddepreciation                                float64 `json:"AccumulatedDepreciation"`
			Totallongtermliabilities                               float64 `json:"TotalLongTermLiabilities"`
			Curliabiandprovision                                   float64 `json:"CurLiabiAndProvision"`
			Tradeandotherpayables                                  float64 `json:"TradeandOtherPayables"`
			Dividendpayable                                        float64 `json:"DividendPayable"`
			Otherliabilitiesandprovisions                          float64 `json:"OtherLiabilitiesandProvisions"`
			Corporateincometaxliabilities                          float64 `json:"CorporateIncomeTaxLiabilities"`
			Provisionsforpossiblelosses                            float64 `json:"ProvisionsforPossibleLosses"`
			Sourcesoffunds                                         float64 `json:"SourcesOfFunds"`
			Lifeinsurancefund                                      float64 `json:"LifeInsuranceFund"`
			Catastrophereserve                                     float64 `json:"CatastropheReserve"`
			Longtermloanborrowings                                 float64 `json:"LongTermLoanBorrowings"`
			Usesoffunds                                            float64 `json:"UsesOfFunds"`
			Fixedassetsnet                                         float64 `json:"FixedAssetsNet"`
			Longterminvestmentsloan                                float64 `json:"LongTermInvestmentsLoan"`
			Expenditureextentnotwrittenoff                         float64 `json:"ExpenditureExtentNotWrittenoff"`
			Miscellaneousexpensestotheextentnotwrittenoff          float64 `json:"MiscellaneousExpensesToTheExtentNotWrittenOff"`
			Losstransferredfromprofitlossaccount                   float64 `json:"LossTransferredFromProfitLossAccount"`
			Detailsofnetcurrentassets                              float64 `json:"DetailsOfNetCurrentAssets"`
			Currentassetsloansandadvances                          float64 `json:"CurrentAssetsLoansAndAdvances"`
			Shortterminvestmentsandloans                           float64 `json:"ShortTermInvestmentsAndLoans"`
			Currentliabilitiesandprovisions                        float64 `json:"CurrentLiabilitiesAndProvisions"`
			Currentliabilities                                     float64 `json:"CurrentLiabilities"`
			Provisionforunexpiredrisk                              float64 `json:"ProvisionForUnexpiredRisk"`
			Provisionforoutstandingclaims                          float64 `json:"ProvisionForOutstandingClaims"`
			Otherprovision                                         float64 `json:"OtherProvision"`
			Totalfunds                                             float64 `json:"TotalFunds"`
			Fixedassets                                            float64 `json:"FixedAssets"`
			Depreciation                                           float64 `json:"Depreciation"`
			Workinprogress                                         float64 `json:"WorkInProgress"`
			Cash                                                   float64 `json:"Cash"`
			Receivables                                            float64 `json:"Receivables"`
			Advancesprepaymentsloansdeposits                       float64 `json:"AdvancesPrepaymentsLoansDeposits"`
			Inventory                                              float64 `json:"Inventory"`
			Totalcurrentassets                                     float64 `json:"TotalCurrentAssets"`
			Stliabilities                                          float64 `json:"StLiabilities"`
			Deferredliabilities                                    float64 `json:"DeferredLiabilities"`
			Totalstliabilities                                     float64 `json:"TotalStLiabilities"`
			Applicationoffunds                                     float64 `json:"ApplicationOfFunds"`
			Capitalliabilities                                     float64 `json:"CapitalLiabilities"`
			Capitalreserves                                        float64 `json:"CapitalReserves"`
			Sharecapital                                           float64 `json:"ShareCapital"`
			Reservesurplus                                         float64 `json:"ReserveSurplus"`
			Mediumlongtermloans                                    float64 `json:"MediumLongTermLoans"`
			Securedloan                                            float64 `json:"SecuredLoan"`
			Unsecuredloan                                          float64 `json:"UnsecuredLoan"`
			Deferredtaxliability                                   float64 `json:"DeferredTaxLiability"`
			Grandtotal                                             float64 `json:"GrandTotal"`
			Assets                                                 float64 `json:"Assets"`
			Netfixedassets                                         float64 `json:"NetFixedAssets"`
			Capitalworkinprogress                                  float64 `json:"CapitalWorkInProgress"`
			Noncoreassets                                          float64 `json:"NonCoreAssets"`
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
			Deferredtaxassets                                      float64 `json:"DeferredTaxAssets"`
			Totalnoncurrentassets                                  float64 `json:"TotalNonCurrentAssets"`
			Currentassetsloansadvances                             float64 `json:"CurrentAssetsLoansAdvances"`
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
			Lesscurrentliabilitiesprovisions                       float64 `json:"LessCurrentLiabilitiesProvisions"`
			Tradeotherpayables                                     float64 `json:"TradeOtherPayables"`
			Proposeddividends                                      float64 `json:"ProposedDividends"`
			Provisionsa                                            float64 `json:"ProvisionsA"`
			Shorttermloans                                         float64 `json:"ShortTermLoans"`
			Totalcurrentliabilities                                float64 `json:"TotalCurrentLiabilities"`
			Netcurrentassets                                       float64 `json:"NetCurrentAssets"`
			Deferredrevenueexpenditure                             float64 `json:"DeferredRevenueExpenditure"`
			Lossasprofitlossaccount                                float64 `json:"LossasProfitLossAccount"`
			Grandtotala                                            float64 `json:"GrandTotalA"`
			Loansextended                                          float64 `json:"LoansExtended"`
			Cashandcashequivalents                                 float64 `json:"CashandCashEquivalents"`
			Termdeposits                                           float64 `json:"TermDeposits"`
			Othercurrentassets                                     float64 `json:"OtherCurrentAssets"`
			Miscassets                                             float64 `json:"MiscAssets"`
			Totalusesoffunds                                       float64 `json:"TotalUsesofFunds"`
			Preferredshares                                        float64 `json:"PreferredShares"`
			Intangibleassets                                       float64 `json:"IntangibleAssets"`
			Financialassets                                        float64 `json:"FinancialAssets"`
			Defferedtaxassets                                      float64 `json:"DefferedTaxAssets"`
			Currentassets                                          float64 `json:"CurrentAssets"`
			Capitalworkprogress                                    float64 `json:"CapitalWorkProgress"`
			Tradereceivables                                       float64 `json:"TradeReceivables"`
			Prepaidreceivables                                     float64 `json:"PrepaidReceivables"`
			Cashequivalents                                        float64 `json:"CashEquivalents"`
			Investmentfd                                           float64 `json:"InvestmentFd"`
			Equityliabilities                                      float64 `json:"EquityLiabilities"`
			Noncurrentliabilities                                  float64 `json:"NonCurrentLiabilities"`
			Defferedtaxliabilities                                 float64 `json:"DefferedTaxLiabilities"`
			Ncprovisions                                           float64 `json:"NcProvisions"`
			Totalnoncurrentliabilities                             float64 `json:"TotalNonCurrentLiabilities"`
			Tradepaybales                                          float64 `json:"TradePaybales"`
			Incometaxliability                                     float64 `json:"IncomeTaxLiability"`
			Miscellaneousexpenditure                               float64 `json:"MiscellaneousExpenditure"`
			Totaleqli                                              float64 `json:"TotalEqLi"`
			Debenturesandbonds                                     float64 `json:"DebenturesAndBonds"`
			Domesticcurrency                                       float64 `json:"DomesticCurrency"`
			Foreigncurrency                                        float64 `json:"ForeignCurrency"`
			Incometaxliabilities                                   float64 `json:"IncomeTaxLiabilities"`
			Realestateloan                                         float64 `json:"RealEstateLoan,omitempty"`
			Residentialrealestateloanabove10Million                float64 `json:"ResidentialRealEstateLoanAbove10Million,omitempty"`
			Businesscomplexandresidentialapartmentconstructionloan float64 `json:"BusinessComplexAndResidentialApartmentConstructionLoan,omitempty"`
			Incomegeneratingcommercialcomplexloan                  float64 `json:"IncomeGeneratingCommercialComplexLoan,omitempty"`
			Otherrealestateloan                                    float64 `json:"OtherRealEstateLoan,omitempty"`
			Personalhomeloanupto10Millionorless                    float64 `json:"PersonalHomeLoanUpto10MillionOrLess,omitempty"`
			Margintypeloan                                         float64 `json:"MarginTypeLoan,omitempty"`
			Termloan                                               float64 `json:"TermLoan,omitempty"`
			Overdraftloanandtrloanandwcloan                        float64 `json:"OverdraftLoanAndTrLoanAndWcLoan,omitempty"`
			Otherloan                                              float64 `json:"OtherLoan"`
			Nonbankingassets                                       float64 `json:"NonBankingAssets"`
			Moneyatcallandshortnotice                              float64 `json:"MoneyAtCallAndShortNotice"`
			Reportedpeannualized                                   float64 `json:"ReportedPeAnnualized"`
			Bookvaluepershare                                      float64 `json:"BookValuePerShare"`
		} `json:"data"`
	} `json:"message"`
}

type IncomeStatementData struct {
	Response float64 `json:"response"`
	Error    string  `json:"error"`
	Message  struct {
		Ticker string `json:"ticker"`
		Data   []struct {
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
			Dividendincome                                   float64 `json:"DividendIncome"`
			Feesandcommissionincome                          float64 `json:"FeesandCommissionIncome"`
			Otherincome                                      float64 `json:"OtherIncome"`
			Totalincome                                      float64 `json:"TotalIncome"`
			Operatingexpenses                                float64 `json:"OperatingExpenses"`
			Operationalexpenses                              float64 `json:"OperationalExpenses"`
			Provisionsforinvestmentrisksandotherlosses       float64 `json:"ProvisionsforInvestmentRisksandOtherLosses"`
			Miscexpenses                                     float64 `json:"MiscExpenses"`
			Totaloperatingexpenses                           float64 `json:"TotaloperatingExpenses"`
			Ebitda                                           float64 `json:"Ebitda"`
			Ebit                                             float64 `json:"Ebit"`
			Financialexpenses                                float64 `json:"FinancialExpenses"`
			Profitbeforebonusandtaxes                        float64 `json:"ProfitbeforeBonusandTaxes"`
			Provisionforbonus                                float64 `json:"ProvisionforBonus"`
			Profitbeforetaxes                                float64 `json:"ProfitbeforeTaxes"`
			Provisionfortaxes                                float64 `json:"ProvisionforTaxes"`
			Energysales                                      float64 `json:"EnergySales"`
			Costofproduction                                 float64 `json:"CostOfProduction"`
			Grossprofit                                      float64 `json:"GrossProfit"`
			Forexgainloss                                    float64 `json:"ForexGainLoss"`
			Adminexpenses                                    float64 `json:"AdminExpenses"`
			Interestincomeexpense                            float64 `json:"InterestIncomeExpense"`
			Depreciation                                     float64 `json:"Depreciation"`
			Provisions                                       float64 `json:"Provisions"`
			Taxes                                            float64 `json:"Taxes"`
			Bonus                                            float64 `json:"Bonus"`
			Profitaftertax                                   float64 `json:"ProfitAfterTax"`
			Roomsrestaurantsotherincome                      float64 `json:"RoomsRestaurantsOtherIncome"`
			Otherbusinessincome                              float64 `json:"OtherBusinessIncome"`
			Expenditure                                      float64 `json:"Expenditure"`
			Foodbeverageexpense                              float64 `json:"FoodBeverageExpense"`
			Salariesemployeebenefit                          float64 `json:"SalariesEmployeeBenefit"`
			Operatingadminexpenses                           float64 `json:"OperatingAdminExpenses"`
			Totaloperatingexpenditure                        float64 `json:"TotalOperatingExpenditure"`
			Grossoperatingprofit                             float64 `json:"GrossOperatingProfit"`
			Interest                                         float64 `json:"Interest"`
			Depreciationamortizationexpenses                 float64 `json:"DepreciationAmortizationExpenses"`
			Netprofitbeforebonustax                          float64 `json:"NetProfitBeforeBonusTax"`
			Deferredexpenses                                 float64 `json:"DeferredExpenses"`
			Exchangefluctuation                              float64 `json:"ExchangeFluctuation"`
			Salefixedasset                                   float64 `json:"SaleFixedAsset"`
			Provisionhousing                                 float64 `json:"ProvisionHousing"`
			Provisionstaffbonus                              float64 `json:"ProvisionStaffBonus"`
			Netprofitbeforetax                               float64 `json:"NetProfitBeforeTax"`
			Provisionincometax                               float64 `json:"ProvisionIncomeTax"`
			Netprofitaftertax                                float64 `json:"NetProfitAfterTax"`
			Income                                           float64 `json:"Income"`
			Transferredfromrevenueaccount                    float64 `json:"TransferredFromRevenueAccount"`
			Transferredfromlifefund                          float64 `json:"TransferredFromLifeFund"`
			Incomefrominvestmentloansothers                  float64 `json:"IncomeFromInvestmentLoansOthers"`
			Provisionwrittenback                             float64 `json:"ProvisionWrittenBack"`
			Expenses                                         float64 `json:"Expenses"`
			Managementexpenses                               float64 `json:"ManagementExpenses"`
			Amortizationexpenses                             float64 `json:"AmortizationExpenses"`
			Sharerelatedexpenses                             float64 `json:"ShareRelatedExpenses"`
			Otherexpenses                                    float64 `json:"OtherExpenses"`
			Provisionforlosses                               float64 `json:"ProvisionForLosses"`
			Provisionforstaffhousing                         float64 `json:"ProvisionForStaffHousing"`
			Provisionforstaffbonus                           float64 `json:"ProvisionForStaffBonus"`
			Adjustedincometax                                float64 `json:"AdjustedIncomeTax"`
			Transferredtolifefund                            float64 `json:"TransferredToLifeFund"`
			Netprofit                                        float64 `json:"NetProfit"`
			Variance                                         float64 `json:"Variance"`
			Totalrevenue                                     float64 `json:"TotalRevenue"`
			Cogs                                             float64 `json:"Cogs"`
			Sga                                              float64 `json:"Sga"`
			Totalrevenuemrq                                  float64 `json:"TotalRevenueMrq"`
			Cogsmrq                                          float64 `json:"CogsMrq"`
			Sgamrq                                           float64 `json:"SgaMrq"`
			Ebitdamrq                                        float64 `json:"EbitdaMrq"`
			Ebitmrq                                          float64 `json:"EbitMrq"`
			Sales                                            float64 `json:"Sales"`
			Costofmaterials                                  float64 `json:"CostofMaterials"`
			Manufacturingexpenses                            float64 `json:"ManufacturingExpenses"`
			Administrativeexpenses                           float64 `json:"AdministrativeExpenses"`
			Sellingdistributionexpenses                      float64 `json:"SellingDistributionExpenses"`
			Nonrecurring                                     float64 `json:"NonRecurring"`
			Totalexpenditure                                 float64 `json:"TotalExpenditure"`
			Provisionobsolescence                            float64 `json:"ProvisionObsolescence"`
			Provisionwriteoff                                float64 `json:"ProvisionWriteOff"`
			Housingfundallocation                            float64 `json:"HousingFundAllocation"`
			Profitbeforeprovisionbonus                       float64 `json:"ProfitBeforeProvisionBonus"`
			Provisiontaxation                                float64 `json:"ProvisionTaxation"`
			Current                                          float64 `json:"Current"`
			Deferred                                         float64 `json:"Deferred"`
			Dividend                                         float64 `json:"Dividend"`
			Feesandcommissionanddiscount                     float64 `json:"FeesAndCommissionAndDiscount"`
			Foreignexchangegainandloss                       float64 `json:"ForeignExchangeGainAndLoss"`
			Operatingprofitbeforeprovision                   float64 `json:"OperatingProfitBeforeProvision"`
			Provisionforpossiblelosses                       float64 `json:"ProvisionForPossibleLosses"`
			Nonoperatingincomeandexpenses                    float64 `json:"NonOperatingIncomeAndExpenses"`
			Writebackofprovisionforpossibleloss              float64 `json:"WriteBackOfProvisionForPossibleLoss"`
			Profitfromregularactivities                      float64 `json:"ProfitFromRegularActivities"`
			Extraordinaryincomeandexpenses                   float64 `json:"ExtraordinaryIncomeAndExpenses"`
			Provisionfortax                                  float64 `json:"ProvisionForTax"`
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
