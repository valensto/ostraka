package static

type staticWorkflow struct {
	Event   staticEvent    `yaml:"event"   validate:"required,dive,required"`
	Inputs  []staticInput  `yaml:"inputs"   validate:"required,dive,required"`
	Outputs []staticOutput `yaml:"outputs"   validate:"required,dive,required"`
}

type staticEvent struct {
	Format string        `yaml:"format"   validate:"required"`
	Fields []staticField `yaml:"fields"   validate:"required,dive,required"`
}

type staticField struct {
	Name     string `yaml:"name"   validate:"required"`
	DataType string `yaml:"data_type"   validate:"required"`
	Required bool   `yaml:"required"`
}

type staticInput struct {
	Name    string        `yaml:"name"   validate:"required"`
	Source  string        `yaml:"source"   validate:"required"`
	Decoder staticDecoder `yaml:"decoder"   validate:"required,dive,required"`
	Params  any           `yaml:"params"   validate:"required"`
}

type staticDecoder struct {
	Format  string         `yaml:"format"   validate:"required"`
	Mappers []staticMapper `yaml:"mappers"   validate:"required,dive,required"`
}

type staticMapper struct {
	Source string `yaml:"source"   validate:"required"`
	Target string `yaml:"target"   validate:"required"`
}

type staticOutput struct {
	Name        string           `yaml:"name"   validate:"required"`
	Destination string           `yaml:"destination"   validate:"required"`
	Params      any              `yaml:"params"   validate:"required"`
	Condition   *staticCondition `yaml:"condition,omitempty" `
}

type staticCondition struct {
	Field      string            `yaml:"field,omitempty"`
	Operator   string            `yaml:"operator"`
	Value      any               `yaml:"value,omitempty"`
	Conditions []staticCondition `yaml:"conditions,omitempty"`
}
