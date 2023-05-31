package static

type workflowModel struct {
	Event   eventModel    `yaml:"event"   validate:"required,dive,required"`
	Inputs  []inputModel  `yaml:"inputs"   validate:"required,dive,required"`
	Outputs []outputModel `yaml:"outputs"   validate:"required,dive,required"`
}

type eventModel struct {
	Format string       `yaml:"format"   validate:"required"`
	Fields []fieldModel `yaml:"fields"   validate:"required,dive,required"`
}

type fieldModel struct {
	Name     string `yaml:"name"   validate:"required"`
	DataType string `yaml:"data_type"   validate:"required"`
	Required bool   `yaml:"required"`
}

type inputModel struct {
	Name    string       `yaml:"name"   validate:"required"`
	Source  string       `yaml:"source"   validate:"required"`
	Decoder decoderModel `yaml:"decoder"   validate:"required,dive,required"`
	Params  any          `yaml:"params"   validate:"required"`
}

type decoderModel struct {
	Format  string        `yaml:"format"   validate:"required"`
	Mappers []mapperModel `yaml:"mappers"   validate:"required,dive,required"`
}

type mapperModel struct {
	Source string `yaml:"source"   validate:"required"`
	Target string `yaml:"target"   validate:"required"`
}

type outputModel struct {
	Name        string          `yaml:"name"   validate:"required"`
	Destination string          `yaml:"destination"   validate:"required"`
	Params      any             `yaml:"params"   validate:"required"`
	Condition   *conditionModel `yaml:"condition,omitempty" `
}

type conditionModel struct {
	Field      string           `yaml:"field,omitempty"`
	Operator   string           `yaml:"operator"`
	Value      any              `yaml:"value,omitempty"`
	Conditions []conditionModel `yaml:"conditions,omitempty"`
}
