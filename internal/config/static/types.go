package static

// ContentFile is a map of file name and content.
type ContentFile map[string][]byte

type workflowModel struct {
	Name        string           `yaml:"name"   validate:"required"`
	Middlewares middlewaresModel `yaml:"middlewares"   validate:"required"`
	EventType   eventTypeModel   `yaml:"event_type"   validate:"required,dive,required"`
	Inputs      []inputModel     `yaml:"inputs"   validate:"required,dive,required"`
	Outputs     []outputModel    `yaml:"outputs"   validate:"required,dive,required"`
}

type eventTypeModel struct {
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

type middlewaresModel struct {
	CORS []corsModel `yaml:"cors"`
	Auth []authModel `yaml:"auth"`
}

type corsModel struct {
	Name             string   `yaml:"name"`
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

type authModel struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Params any    `yaml:"params"`
}
