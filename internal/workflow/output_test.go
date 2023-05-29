package workflow

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOutput_MatchConditions(t *testing.T) {
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
					Field:    "status",
					Operator: "eq",
					Value:    "completed",
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
					Field:    "boolean_field",
					Operator: "eq",
					Value:    true,
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
					Field:    "status",
					Operator: "eq",
					Value:    "completed",
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
					Operator: "and",
					Conditions: []Condition{
						{
							Field:    "status",
							Operator: "eq",
							Value:    "completed",
						},
						{
							Field:    "paid",
							Operator: "eq",
							Value:    true,
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
					Operator: "or",
					Conditions: []Condition{
						{
							Field:    "status",
							Operator: "eq",
							Value:    "completed",
						},
						{
							Field:    "status",
							Operator: "eq",
							Value:    "failed",
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
					Operator: "and",
					Conditions: []Condition{
						{
							Field:    "status",
							Operator: "eq",
							Value:    "completed",
						},
						{
							Field:    "status",
							Operator: "eq",
							Value:    "failed",
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
				Name:      tt.fields.Name,
				Type:      tt.fields.Type,
				Params:    tt.fields.Params,
				Condition: tt.fields.Condition,
			}

			match := o.Condition.IsMatching(tt.event)
			require.Equal(t, tt.isMatching, match)
		})
	}
}
