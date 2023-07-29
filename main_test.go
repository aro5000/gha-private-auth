package main

import (
	"crypto/rsa"
	"testing"
	"time"
)

var pemStr = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

func Test_parsePem(t *testing.T) {
	tests := []struct {
		name    string
		pemStr  string // create with: openssl genrsa -traditional -out test.pem 4096
		wantErr bool
	}{
		{
			name:    "Working PEM",
			pemStr:  pemStr,
			wantErr: false,
		},
		{
			name:    "Invalid PEM",
			pemStr:  "fake PEM",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parsePem(tt.pemStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_createJwt(t *testing.T) {
	key, _ := parsePem(pemStr)

	type args struct {
		key             *rsa.PrivateKey
		appId           string
		currentUnixDate int64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Functional Key",
			args: args{
				key:             key,
				appId:           "123",
				currentUnixDate: time.Date(2023, 7, 29, 0, 0, 0, 0, time.UTC).Unix(),
			},
			want:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTA1ODk0MDAsImlhdCI6MTY5MDU4ODgwMCwiaXNzIjoxMjN9.CaGW2p-zmfTsQmledlpQhRipxF4nVVa1TB--uTr2JbEyo8HC2oMBDYbNjKktrsJSFu0LsC3yK_4iS-kT9ojkyEjQTrvpnK30JL-Aaf7BqEAcmk46OmEmLXnf9GpmclwxA2m9OYJdXvHgxr5VPEUKMfKByQM9uFOaTxt3pxwFGU9hjWISKhWyIPM0Mbc75dNiJnM-5tqUSb85VJplMWFbpx9wUK70DKZT45Ypy_b5BBK7uFYsX5YdQKDDEoWy22slGkAyj5KsY8hRrNHewFd_d9179W5l9eQmK-CUIasN-Faz5IKw_xD9zLSKF3Lbjd4-YVAJUXlTl9_HxOBNBH3tuM1HwvNh5hp4Me7aCHANE1tmaivIjl_Pz6wJHs25G3EvR7Z0L3EuEj6nyZgPH9TIhh8NG3ncyTeHSZnZX18Q4R36Q0dkDUej69qc1TQNtgfwnbaKUkofrmnsIP2n9Hbx-tHeG5rl2XWb_mhnfMxOCMBSwIANuDQtxUhzN-3fnRWKRf-YhvaeJ52Lss_7nmy6B2t9DsMxiN0kg3tUB1zl3tMM_M-avmtw5ebfFw4UPL84_Rg98feoGGYzxfHFSAx4mEL1vbChUUZZsqJxGYLQmzKx6Q8FsCJyq4f_wHnnmC2FWLiXUsIzPSuXBtAqFwxnJAVn3c9xXbhv8NF5eNTv_Gs",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createJwt(tt.args.key, tt.args.appId, tt.args.currentUnixDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("createJwt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createJwt() = %v, want %v", got, tt.want)
			}
		})
	}
}
