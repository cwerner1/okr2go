package okr2go

type Objective struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	KeyResults  []KeyResult `json:"keyResults"`
}

type KeyResult struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Current      int64    `json:"current"`
	Target       int64    `json:"target"`
	Contributors []string `json:"contributors"`
	Comments     []string `json:"comments"`
}
