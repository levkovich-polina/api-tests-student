package favorites

type FavoriteListResponse struct {
	Items      []FavoriteItem `json:"items"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}

type FavoriteItem struct {
	FavoriteID    string             `json:"favorite_id"`
	Advertisement AdvertisementBrief `json:"advertisement"`
	AddedAt       string             `json:"added_at"`
}

type AdvertisementBrief struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Price  float64 `json:"price"`
	UserID string  `json:"user_id"`
}

type FavoriteResponse struct {
	FavoriteID    string             `json:"favorite_id"`
	Advertisement AdvertisementBrief `json:"advertisement"`
	AddedAt       string             `json:"added_at"`
}