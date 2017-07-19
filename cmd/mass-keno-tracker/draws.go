package main

type DrawResponse struct {
	Bonus        string `json:"bonus"`
	Date         string `json:"date"`
	Id           string `json:"id"`
	NormalizedId int    `json:"normalized_id"`
	Value        string `json:"value"`
}

type StateDrawResponse struct {
	Id    string `json:"draw_id"`
	Date  string `json:"draw_date_value"`
	Value string `json:"winning_num"`
	Bonus string `json:"bonus"`
}
