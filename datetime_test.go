package goprsc

import (
	"encoding/json"
	"testing"
	"time"
)

const (
	referenceDateTimeStr = `"2017-01-02T15:47:59+0100"`
)

var (
	referenceDateTime = time.Date(2017, 01, 02, 14, 47, 59, 0, time.UTC)
)

func TestDateTime_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		desc    string
		data    string
		want    DateTime
		wantErr bool
		equal   bool
	}{
		{"Reference", referenceDateTimeStr, DateTime{referenceDateTime}, false, true},
		{"Mismatch", referenceDateTimeStr, DateTime{}, false, false},
	}
	for _, tc := range testCases {
		var got DateTime
		err := json.Unmarshal([]byte(tc.data), &got)
		if gotErr := err != nil; gotErr != tc.wantErr {
			t.Errorf("%s: gotErr=%v, wantErr=%v, err=%v", tc.desc, gotErr, tc.wantErr, err)
			continue
		}
		equal := got.Equal(tc.want)
		if equal != tc.equal {
			t.Errorf("%s: got=%#v, want=%#v, equal=%v, want=%v", tc.desc, got, tc.want, equal, tc.equal)
		}
	}
}
