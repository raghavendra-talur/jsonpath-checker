package jsonpath

import (
	"encoding/json"
	"strconv"
	"testing"
)

type testCases struct {
	jsonPathExprs string
	result        bool
}

var jsonText1 = []byte(`{
    "kind": "Deployment",
    "spec": {
        "progressDeadlineSeconds": 600,
        "replicas": 1,
        "revisionHistoryLimit": 10
    },
    "status": {
        "replicas": 1,
        "conditions": [
            {
                "status": "True",
                "type": "Progressing"
            },
            {
                "status": "True",
                "type": "Available"
            }
        ]
    }
    }`)

var testCasesData = []testCases{
	{
		jsonPathExprs: "{$.status.conditions[0].status} == True",
		result:        true,
	},
	{
		jsonPathExprs: "{$.spec.replicas} == 1",
		result:        true,
	},
	{
		jsonPathExprs: "{$.status.conditions[0].status} == {True}",
		result:        false,
	},
}

func TestXYZ(t *testing.T) {
	for i, tt := range testCasesData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var jsonData map[string]interface{}
			err := json.Unmarshal(jsonText1, &jsonData)
			if err != nil {
				t.Error(err)
			}
			_, err = EvaluateCheckHook(jsonData, tt.jsonPathExprs)
			if (err == nil) != tt.result {
				t.Errorf("EvaluateCheckHook() = %v, want %v", err, tt.result)
			}
		})
	}
}

func Test_isValidJsonPathExpression(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Simple expression",
			args: args{
				expr: "$.spec.replicas",
			},
			want: true,
		},
		{
			name: "no $ at the start",
			args: args{
				expr: "{.spec.replicas}",
			},
			want: true,
		},
		{
			name: "element in array",
			args: args{
				expr: "{$.status.conditions[0].status}",
			},
			want: true,
		},
		{
			name: "spec 1",
			args: args{
				expr: "{$.status.readyReplicas}",
			},
			want: true,
		},
		{
			name: "spec 2",
			args: args{
				expr: "{$.status.containerStatuses[0].ready}",
			},
			want: true,
		},
		{
			name: "spec 3a",
			args: args{
				expr: "{True}",
			},
			want: false,
		},
		{
			name: "spec 3b",
			args: args{
				expr: "{False}",
			},
			want: false,
		},
		{
			name: "spec 3c",
			args: args{
				expr: "{true}",
			},
			want: true,
		},
		{
			name: "spec 3d",
			args: args{
				expr: "{false}",
			},
			want: true,
		},
		{
			name: "Spec 4",
			args: args{
				expr: "{$.spec.replicas}",
			},
			want: true,
		},
		{
			name: "expression with == operator",
			args: args{
				expr: "$.store.book[?(@.price > 10)].title==$.store.book[0].title",
			},
			want: true,
		},
		{
			name: "expression with > operator",
			args: args{
				expr: "$.store.book[?(@.author CONTAINS 'Smith')].price>20",
			},
			want: true,
		},
		{
			name: "expression with >= operator",
			args: args{
				expr: "$.user.age>=$.minimum.age",
			},
			want: true,
		},
		{
			name: "expression with < operator",
			args: args{
				expr: "$.user.age<$.maximum.age",
			},
			want: true,
		},
		{
			name: "expression with <= operator",
			args: args{
				expr: "$.user.age<=$.maximum.age",
			},
			want: true,
		},
		{
			name: "expression with != operator",
			args: args{
				expr: "$.user.age!=$.maximum.age",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidJsonPathExpression(tt.args.expr); got != tt.want {
				t.Errorf("isValidJsonPathExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}
