package command

import "testing"

func Test_parseServerURL(t *testing.T) {
	tests := []struct {
		name     string
		u        string
		withPort bool
		want     string
		wantErr  bool
	}{
		{name: "default", u: "https://gobyexample.com/url-parsing:8080", withPort: true, want: "https://gobyexample.com:8089", wantErr: false},
		{name: "no port", u: "http://gobyexample.com/url-parsing:8080", withPort: false, want: "https://gobyexample.com", wantErr: false},
		{name: "bad url", u: "mail://gg.com/url-parsing:8080", withPort: true, want: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseServerURL(tt.u, tt.withPort)
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
