package main

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestCreateClientWithRetries(t *testing.T) {
	type args struct {
		client        interface{}
		err           error
		retries       int
		retryInterval time.Duration
	}

	tests := []struct {
		name            string
		args            args
		expectedRetries int
		want            interface{}
		wantErr         bool
	}{
		{
			name: "getClient returns a valid client",
			args: args{
				client: "client",
				err:    nil,
			},
			expectedRetries: 0,
			want:            "client",
			wantErr:         false,
		},
		{
			name: "getClient returns an error after no retries",
			args: args{
				client: nil,
				err:    errors.New("error"),
			},
			expectedRetries: 0,
			want:            nil,
			wantErr:         true,
		},
		{
			name: "getClient returns an error after retries",
			args: args{
				client:        nil,
				err:           errors.New("error"),
				retries:       3,
				retryInterval: time.Millisecond * 1,
			},
			expectedRetries: 3,
			want:            nil,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invocations := 0
			getClient := func() (interface{}, error) {
				invocations++
				return tt.args.client, tt.args.err
			}

			got, err := createClientWithRetries(getClient, tt.args.retries, tt.args.retryInterval)

			actualRetries := invocations - 1

			if actualRetries != tt.expectedRetries {
				t.Errorf("createClientWithRetries() got %v retries, expected %v", actualRetries, tt.expectedRetries)
				return
			} else if (err != nil) != tt.wantErr {
				t.Errorf("createClientWithRetries() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && tt.wantErr {
				return
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createClientWithRetries() = %v, want %v", got, tt.want)
			}
		})
	}
}
