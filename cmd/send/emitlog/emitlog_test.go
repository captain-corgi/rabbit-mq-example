package emitlog

import (
	"os"
	"testing"
)

func Test_severityFrom(t *testing.T) {
	tests := map[string]struct {
		args []string
		want string
	}{
		"1. Array empty":                {[]string{}, "info"},
		"2. Array has 1 empty element":  {[]string{""}, "info"},
		"3. First element is empty":     {[]string{"", "something"}, "info"},
		"4. First element is info":      {[]string{"info", "something"}, "info"},
		"5. First element is warning":   {[]string{"warning", "something"}, "warning"},
		"6. First element is error":     {[]string{"error", "something"}, "error"},
		"7. First element is something": {[]string{"something", "something"}, "info"},
		"8. First element is ` `":       {[]string{" ", "something"}, "info"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			os.Args = append([]string{"Test_severityFrom"}, tt.args...)
			if got := severityFrom(tt.args); got != tt.want {
				t.Errorf("severityFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}
