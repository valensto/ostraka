package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_parseFile(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		want          *File
		expectedError error
	}{
		{
			name: "should parse file without error",
			path: "__test__/valid.yaml",
			want: &File{
				Inputs: []Input{
					{
						Name: "webhook-orders",
						Type: "webhook",
						Params: WebhookParams{
							Endpoint: "/webhook/orders",
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
					},
					{
						Name: "mqtt-orders",
						Type: "mqtt",
						Params: MQTTParams{
							Broker:   "mqtt.example.com",
							User:     "my-user",
							Password: "my-password",
							Topic:    "my-topic",
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
				},
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
								Operator: "and",
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
				},
			},
			expectedError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFile(tt.path)
			if tt.expectedError != nil {
				require.ErrorIs(t, err, tt.expectedError)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.Event, got.Event)
			require.Equal(t, tt.want.Inputs, got.Inputs)

			for i, _ := range tt.want.Outputs {
				require.Equal(t, tt.want.Outputs[i].Name, got.Outputs[i].Name)
				require.Equal(t, tt.want.Outputs[i].Type, got.Outputs[i].Type)
				require.Equal(t, tt.want.Outputs[i].Params, got.Outputs[i].Params)
				//require.Equal(t, tt.want.Outputs[i].Conditions, got.Outputs[i].Conditions)
			}
		})
	}
}
