package model

// Team represents a team in a sports event
type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Event represents a sports event from Pinnacle API
type Event struct {
	ID      int      `json:"id"`
	Home    Team     `json:"home"`
	Away    Team     `json:"away"`
	Periods []Period `json:"periods"`
}
