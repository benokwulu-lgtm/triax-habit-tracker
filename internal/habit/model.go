package habit

type Habit struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Frequency   string `json:"frequency"`
	IsActive    bool   `json:"is_active"`
}
