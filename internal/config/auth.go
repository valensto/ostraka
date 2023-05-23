package config

type Auth struct {
	Type    string  `yaml:"type" validate:"required"`
	Secret  string  `yaml:"secret" validate:"required"`
	Encoder Encoder `yaml:"encoder" validate:"required,dive,required"`
}

type Encoder struct {
	Type   string  `yaml:"type" validate:"required"`
	Fields []Field `yaml:"fields" validate:"required,dive,required"`
}
