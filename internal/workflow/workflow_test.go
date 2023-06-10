package workflow

/*func Test_parseFile(t *testing.T) {
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
					format: "json",
					fields: []Field{
						{
							name:     "customerId",
							dataType: "string",
							required: true,
						},
						{
							name:     "orderNumber",
							dataType: "int",
							required: true,
						},
						{
							name:     "orderStatus",
							dataType: "string",
							required: true,
						},
						{
							name:     "nonRequiredField",
							dataType: "string",
							required: false,
						},
					},
				},
				Subscribers: []Input{
					{
						Name:   "webhook-orders",
						Source: "webhook",
						params: WebhookParams{
							Endpoint: "/webhook/orders",
						},
						Decoder: Decoder{
							Format: "json",
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
						Name:   "mqtt-orders",
						Source: "mqtt",
						params: MQTTParams{
							Broker:        "mqtt.example.com",
							User:          "my-user",
							Password:      "my-password",
							Topic:         "my-topic",
							AutoReconnect: true,
							KeepAlive:     true,
						},
						Decoder: Decoder{
							Format: "json",
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
				Publishers: []Output{
					{
						Name:        "sse-orders",
						Destination: "sse",
						Params: SSEParams{
							Endpoint: "/sse/orders",
							Auth: Auth{
								Type:   "jwt",
								Secret: "my-secret-key",
								Encoder: Encoder{
									Type: "json",
									Fields: []Field{
										{
											name:     "customer_id",
											dataType: "string",
										},
										{
											name:     "customer_email",
											dataType: "string",
										},
									},
								},
							},
						},
						Condition: &Condition{
							field:    "orderStatus",
							operator: "eq",
							value:    "completed",
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

			for i := range tt.want.Subscribers {
				require.Equal(t, tt.want.Subscribers[i].Name, got.Subscribers[i].Name)
				require.Equal(t, tt.want.Subscribers[i].Source, got.Subscribers[i].Type)
				require.Equal(t, tt.want.Subscribers[i].params, got.Subscribers[i].Params)
				require.Equal(t, tt.want.Subscribers[i].Decoder.Mappers, got.Subscribers[i].Decoder.Mappers)
			}

			for i := range tt.want.Publishers {
				require.Equal(t, tt.want.Publishers[i].Name, got.Publishers[i].Name)
				require.Equal(t, tt.want.Publishers[i].Destination, got.Publishers[i].Type)
				require.Equal(t, tt.want.Publishers[i].Params, got.Publishers[i].Params)
				require.Equal(t, tt.want.Publishers[i].Condition, got.Publishers[i].Condition)
			}
		})
	}
}*/
