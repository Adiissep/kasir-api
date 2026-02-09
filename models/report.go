package models

type BestsellingProduct struct {
	Name    string `json:"name"`
	QtySold int    `json:"qty_sold"`
}

type TodayReport struct {
	TotalRevenue        int                  `json:"total_revenue"`
	TotalTransactions   int                  `json:"total_transactions"`
	BestsellingProducts []BestsellingProduct `json:"bestselling_products"`
}
