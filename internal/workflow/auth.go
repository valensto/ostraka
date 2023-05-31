package workflow

type Auth struct {
	Type    string  `json:"type"`
	Secret  string  `json:"secret"`
	Encoder Encoder `json:"encoder"`
}

type Encoder struct {
	Type   string  `json:"type"`
	Fields []Field `json:"fields"`
}
