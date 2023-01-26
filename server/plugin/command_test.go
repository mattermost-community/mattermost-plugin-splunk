package plugin

import "testing"

func Test_parseServerURL(t *testing.T) {
	tests := []struct {
		name    string
		u       string
		want    string
		wantErr bool
	}{
		{name: "default", u: "https://gobyexample.com:8080/url-parsing", want: "https://gobyexample.com:8080", wantErr: false},
		{name: "no port", u: "http://gobyexample.com/url-parsing:8080", want: "http://gobyexample.com", wantErr: false},
		{name: "bad url", u: "mail://gg.com/url-parsing:8080", want: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseServerURL(tt.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseServerURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseServerURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
