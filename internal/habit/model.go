package habit

type Habit struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Frequency   string `json:"frequency"`
	IsActive    bool   `json:"is_active"`
}
