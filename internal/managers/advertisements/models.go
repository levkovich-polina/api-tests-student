package advertisements

type AdvertisementListResponse struct {
	HasNext    bool            `json:"has_next"`
	Items      []Advertisement `json:"items"`
	NextCursor string          `json:"next_cursor"`
	Total      int             `json:"total"`
}

type Advertisement struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CreatedAt   string  `json:"created_at"`
	UserID      string  `json:"user_id"`
}
