package workflow

type Auth struct {
	Type    string  `yaml:"type"`
	Secret  string  `yaml:"secret"`
	Encoder Encoder `yaml:"encoder"`
}

type Encoder struct {
	Type   string  `yaml:"type"`
	Fields []Field `yaml:"fields"`
}
