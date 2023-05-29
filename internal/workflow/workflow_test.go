package workflow

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseFile(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		want  *Workflow
		fails bool
	}{
		{
			name: "should parse file without error",
			path: "__test__/valid.yaml",
			want: &Workflow{
				Event: Event{
					Type: "json",
					Fields: []Field{
						{
							Name:     "customerId",
							Type:     "string",
							Required: true,
						},
						{
							Name:     "orderNumber",
							Type:     "int",
							Required: true,
						},
						{
							Name:     "orderStatus",
							Type:     "string",
							Required: true,
						},
						{
							Name:     "nonRequiredField",
							Type:     "string",
							Required: false,
						},
					},
				},
				Inputs: []Input{
					{
						Name: "webhook-orders",
						Type: "webhook",
						Params: WebhookParams{
							Endpoint: "/webhook/orders",
						},
						Decoder: Decoder{
							Type: "json",
							Mappers: []Mapper{
								{
									Source: "o_customer_id",
									Target: "customerId",
								},
								{
									Source: "o_number",
									Target: "orderNumber",
								},
								{
									Source: "o_status",
									Target: "orderStatus",
								},
							},
						},
					},
					{
						Name: "mqtt-orders",
						Type: "mqtt",
						Params: MQTTParams{
							Broker:        "mqtt.example.com",
							User:          "my-user",
							Password:      "my-password",
							Topic:         "my-topic",
							AutoReconnect: true,
							KeepAlive:     true,
						},
						Decoder: Decoder{
							Type: "json",
							Mappers: []Mapper{
								{
									Source: "customer_id",
									Target: "customerId",
								},
								{
									Source: "number",
									Target: "orderNumber",
								},
								{
									Source: "status",
									Target: "orderStatus",
								},
							},
						},
					},
				},
				Outputs: []Output{
					{
						Name: "sse-orders",
						Type: "sse",
						Params: SSEParams{
							Endpoint: "/sse/orders",
							Auth: Auth{
								Type:   "jwt",
								Secret: "my-secret-key",
								Encoder: Encoder{
									Type: "json",
									Fields: []Field{
										{
											Name: "customer_id",
											Type: "string",
										},
										{
											Name: "customer_email",
											Type: "string",
										},
									},
								},
							},
						},
						Conditions: []Condition{
							{
								Source:   "token",
								Field:    "customerId",
								Operator: "eq",
								Value:    "$customer_id",
							},
							{
								Source:   "event",
								Field:    "orderStatus",
								Operator: "eq",
								Value:    "completed",
							},
						},
					},
				},
			},
			fails: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractWorkflow(tt.path)
			if tt.fails {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.Event, got.Event)

			for i, _ := range tt.want.Inputs {
				require.Equal(t, tt.want.Inputs[i].Name, got.Inputs[i].Name)
				require.Equal(t, tt.want.Inputs[i].Type, got.Inputs[i].Type)
				require.Equal(t, tt.want.Inputs[i].Params, got.Inputs[i].Params)
				require.Equal(t, tt.want.Inputs[i].Decoder.Mappers, got.Inputs[i].Decoder.Mappers)
			}

			for i, _ := range tt.want.Outputs {
				require.Equal(t, tt.want.Outputs[i].Name, got.Outputs[i].Name)
				require.Equal(t, tt.want.Outputs[i].Type, got.Outputs[i].Type)
				require.Equal(t, tt.want.Outputs[i].Params, got.Outputs[i].Params)
				require.Equal(t, tt.want.Outputs[i].Conditions, got.Outputs[i].Conditions)
			}
		})
	}
}
