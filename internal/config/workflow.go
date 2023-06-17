package config

type Workflow struct {
	Name string `yaml:"name" validate:"required"`

	EventType EventType `yaml:"event_type"   validate:"required,dive,required"`

	Middlewares Middlewares `yaml:"middlewares"   validate:"required"`

	Inputs  []Input  `yaml:"inputs"   validate:"required,dive,required"`
	Outputs []Output `yaml:"outputs"   validate:"required,dive,required"`
}

type EventType struct {
	Format string  `yaml:"format"   validate:"required"`
	Fields []Field `yaml:"fields"   validate:"required,dive,required"`
}

type Field struct {
	Name     string `yaml:"name"   validate:"required"`
	DataType string `yaml:"data_type"   validate:"required"`
	Required bool   `yaml:"required"`
}

type Input struct {
	Name    string  `yaml:"name"   validate:"required"`
	Source  string  `yaml:"source"   validate:"required"`
	Decoder Decoder `yaml:"decoder"   validate:"required,dive,required"`
	Params  any     `yaml:"params"   validate:"required"`
}

type Decoder struct {
	Format  string   `yaml:"format"   validate:"required"`
	Mappers []Mapper `yaml:"mappers"   validate:"required,dive,required"`
}

type Mapper struct {
	Source string `yaml:"source"   validate:"required"`
	Target string `yaml:"target"   validate:"required"`
}

type Output struct {
	Name        string     `yaml:"name"   validate:"required"`
	Destination string     `yaml:"destination"   validate:"required"`
	Params      any        `yaml:"params"   validate:"required"`
	Condition   *Condition `yaml:"condition,omitempty" `
}

type Condition struct {
	Field      string      `yaml:"field,omitempty"`
	Operator   string      `yaml:"operator"`
	Value      any         `yaml:"value,omitempty"`
	Conditions []Condition `yaml:"conditions,omitempty"`
}

type Middlewares struct {
	CORS []CORS `yaml:"cors"`
	Auth []Auth `yaml:"auth"`
}

type CORS struct {
	Name             string   `yaml:"name"`
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

type Auth struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Params any    `yaml:"params"`
}
