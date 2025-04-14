package models

type Stats struct {
	ImpressionsCount int     `db:"impressions_count"`
	ClicksCount      int     `db:"clicks_count"`
	Conversion       float64 `db:"conversion"`
	SpentImpressions float64 `db:"spent_impressions"`
	SpentClicks      float64 `db:"spent_clicks"`
	SpentTotal       float64 `db:"spent_total"`
}

type StatsDaily struct {
	Stats
	Date int `db:"date"`
}
