package workflow

/*func TestCondition_MatchConditions(t *testing.T) {
	type fields struct {
		Name      string
		Type      string
		Params    interface{}
		Condition *Condition
	}

	tests := []struct {
		name       string
		fields     fields
		event      map[string]any
		isMatching bool
	}{
		{
			name: "should return true when no conditions are set",
			fields: fields{
				Condition: nil,
			},
			event: map[string]any{
				"status": "completed",
			},
			isMatching: true,
		},
		{
			name: "should return true when condition matched",
			fields: fields{
				Condition: &Condition{
					field:    "status",
					operator: "eq",
					value:    "completed",
				},
			},
			event: map[string]any{
				"status": "completed",
			},
			isMatching: true,
		},
		{
			name: "should return true when conditions matched",
			fields: fields{
				Condition: &Condition{
					field:    "boolean_field",
					operator: "eq",
					value:    true,
				},
			},
			event: map[string]any{
				"boolean_field": true,
			},
			isMatching: true,
		},
		{
			name: "should return false when condition not matched",
			fields: fields{
				Condition: &Condition{
					field:    "status",
					operator: "eq",
					value:    "completed",
				},
			},
			event: map[string]any{
				"status": "failed",
			},
			isMatching: false,
		},
		{
			name: "should return true when 'and' nested condition matched",
			fields: fields{
				Condition: &Condition{
					operator: "and",
					conditions: []Condition{
						{
							field:    "status",
							operator: "eq",
							value:    "completed",
						},
						{
							field:    "paid",
							operator: "eq",
							value:    true,
						},
					},
				},
			},
			event: map[string]any{
				"status": "completed",
				"paid":   true,
			},
			isMatching: true,
		},
		{
			name: "should return true when 'or' nested condition matched",
			fields: fields{
				Condition: &Condition{
					operator: "or",
					conditions: []Condition{
						{
							field:    "status",
							operator: "eq",
							value:    "completed",
						},
						{
							field:    "status",
							operator: "eq",
							value:    "failed",
						},
					},
				},
			},
			event: map[string]any{
				"status": "failed",
			},
			isMatching: true,
		},
		{
			name: "should return true when nested condition not match",
			fields: fields{
				Condition: &Condition{
					operator: "and",
					conditions: []Condition{
						{
							field:    "status",
							operator: "eq",
							value:    "completed",
						},
						{
							field:    "status",
							operator: "eq",
							value:    "failed",
						},
					},
				},
			},
			event: map[string]any{
				"status": "failed",
			},
			isMatching: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Output{
				Name:        tt.fields.Name,
				Destination: tt.fields.Type,
				Params:      tt.fields.Params,
				Condition:   tt.fields.Condition,
			}

			match := o.Condition.Match(tt.event)
			require.Equal(t, tt.isMatching, match)
		})
	}
}*/
