package config

import (
	"encoding/json"
	"testing"
)

var example = `
apps:
  - name: user
    exposes:
     - protocol: grpc
       port: 8080
       paths:
         - "/User.XXXX/"
  - name: rec
    exposes:
     - protocol: grpc
       port: 8081
       paths:
        - "/User.XXXX/"
  - name: api
    exposes:
     - protocol: http
       port: 8082
       paths:
         - "/api/v1/"
         - "/api2/v1/"
environments:
  - name: default
    apps:
      user: 
        grpc: user.prod.nutlet.com:8081
        http: user.prod.nutlet.com:8081
      api: 
        http: api.prod.nutlet.com:8081
      rec: 
        grpc: rec.prod.nutlet.com:8081
`

var a = `
environments:
  - name: default
    apps:
    - name: user
      address: user.prod.nutlet.com:8081
    - name: api
      address: api.prod.nutlet.com:8081
    - name: rec
      address: rec.prod.nutlet.com:8081
`

func Test_unmarshalYaml(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"demo", args{data: []byte(example)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := unmarshalYaml(tt.args.data)
			jd, e := json.Marshal(d)
			t.Log(">>>>", string(jd), "err:", e)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}