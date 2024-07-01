package pinterest

type BoardInfo struct {
	Id   string
	Name string
}

type BoardData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
}
