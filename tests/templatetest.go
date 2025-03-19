package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorono/internal/handlera"
	"gorono/internal/middlas"
	"gorono/internal/models"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

func (suite *TstHandlers) Test_cryptas() {
	type want struct {
		code     int
		response string
		//		err      error
	}

	tests := []struct {
		name            string
	
	}{


	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
		})
	}
}