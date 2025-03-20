package tests

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"gorono/internal/handlera"
	"gorono/internal/models"
	"gorono/internal/privacy"
	"io"
	"net/http"
	"net/http/httptest"
)

func (suite *TstHandlers) Test_cryptas() {
	// type want struct {
	// 	code     int
	// 	response string
	// 	//		err      error
	// }

	tests := []struct {
		name        string
		key         string
		inputString []byte
	}{
		{
			name:        "crypto Right",
			key:         "keykey",
			inputString: []byte("whtatToSend"),
		},
		{
			name:        "crypto Right2",
			key:         "key\"key\"dfgdfgdfg___6567567#$%$#",
			inputString: []byte("whtatToSenddfgdfgdfg#$%#$%#$%dfgdfgdfgdfg\"dfgdfgdfgdfg"),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {

			keyB := md5.Sum([]byte(tt.key))

			coded, err := privacy.EncryptB2B([]byte(tt.inputString), keyB[:])
			suite.Assert().NoError(err)
			ha := privacy.MakeHash(nil, coded, keyB[:])
			haHex := hex.EncodeToString(ha)

			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(coded)) // post crypted tt.inputString
			request.Header.Add("HashSHA256", haHex)

			w := httptest.NewRecorder()

			models.Key = tt.key // for CryptoHandleDecoder
			fu := thecap
			hfunc := http.HandlerFunc(fu)             // make handler from function
			hh := handlera.CryptoHandleDecoder(hfunc) // оборачиваем в мидлварь который расшифрует
			hh.ServeHTTP(w, request)

			res := w.Body
			telo, err := io.ReadAll(res)
			suite.Assert().NoError(err)
			suite.Assert().Equal(tt.inputString, telo)

		})
	}
}
