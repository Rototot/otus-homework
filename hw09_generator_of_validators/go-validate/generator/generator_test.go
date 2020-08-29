package generator

import (
	"bytes"
	"testing"
)

func TestStructureValidatorGenerator_Generate(t *testing.T) {
	type fields struct {
		outPkgName string
	}
	type args struct {
		validateTargets []validationStruct
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantDestination string
		wantErr         bool
	}{
		{
			name: "simple",
			fields: fields{
				"test",
			},
			args: args{
				validateTargets: []validationStruct{
					{
						Name:  "User",
						Alias: "u",
					},
				},
			},
			wantDestination: "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &structureValidatorGenerator{
				outPkgName: tt.fields.outPkgName,
			}
			destination := &bytes.Buffer{}
			err := g.Generate(destination, tt.args.validateTargets)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDestination := destination.String(); gotDestination != tt.wantDestination {
				t.Errorf("Generate() gotDestination = %v, want %v", gotDestination, tt.wantDestination)
			}
		})
	}
}
