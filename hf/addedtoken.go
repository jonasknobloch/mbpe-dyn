package hf

type AddedToken struct {
	ID         int    `json:"id"`
	Content    string `json:"content"`
	SingleWord bool   `json:"single_word"`
	Lstrip     bool   `json:"lstrip"`
	Rstrip     bool   `json:"rstrip"`
	Normalized bool   `json:"normalized"`
	Special    bool   `json:"special"`
}
