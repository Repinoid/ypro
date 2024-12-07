package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_addGauge(t *testing.T) {

	tests := []struct {
		name, metricName string
		metricValue      gauge
		ms               *MemStorage
		wantErr          error
	}{

		{
			name:        "gaaga Ok tst",
			ms:          &memStor,
			metricName:  "Alloc",
			metricValue: gauge(77.77),
			wantErr:     nil,
		},
		{
			name:        "gaaga bad name tst",
			ms:          &memStor,
			metricName:  "Alloca",
			metricValue: gauge(77.77),
			wantErr:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ms.addGauge(tt.metricName, tt.metricValue); err != tt.wantErr {
				t.Errorf("MemStorage.addGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_getGaugeValue(t *testing.T) {

	tests := []struct {
		name, metricName, metricValue string
		ms                            *MemStorage
	}{
		{
			name:        "not existing name",
			ms:          &memStor,
			metricName:  "wtf",
			metricValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ms.getGaugeValue(tt.metricName, &tt.metricValue)
			assert.Error(t, got)
		})
	}
}
