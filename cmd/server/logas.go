package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"
)

type responseData struct {
	status int
	size   int
}
type loggingResponseWriter struct {
	http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
	responseData        *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func WithLogging(origFunc func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	loggedFunc := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		origFunc(&lw, r) // обслуживание оригинального запроса

		duration := time.Since(start)

		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"size", responseData.size, // получаем перехваченный размер ответа
		)
	}
	return loggedFunc
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	//	return w.Writer.Write([]byte("zalup"))
	return w.Writer.Write(b)
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respon http.ResponseWriter, claim *http.Request) {
		rwr := respon
		req := claim
		if strings.Contains(claim.Header.Get("Accept-Encoding"), "gzip") &&
			(strings.Contains(claim.Header.Get("Content-Type"), "application/json") ||
				strings.Contains(claim.Header.Get("Content-Type"), "text/html")) {
			respon.Header().Set("Content-Encoding", "gzip") //
			//claim.Header.Set("Content-Encoding", "gzip")    // без этого в тестах -
			// iteration8_test.go:326:
			//     Error Trace:    y:\GO\ypro\iteration8_test.go:326
			//                                             y:\GO\ypro\suite.go:91
			//     Error:          "" does not contain "gzip"
			//     Test:           TestIteration8/TestGetGzipHandlers/get_info_page
			//     Messages:       Заголовок ответа Content-Encoding содержит несоответствующее значение

			respon.Header().Set("Content-Type", "application/octet-stream") //
			//	claim.Header.Set("Content-Type", "application/octet-stream")    //
			//									req.Header.Set("Content-Type", "application/octet-stream")      //
			//									respon.Header().Set("Content-Type", "application/json") //
			//									req.Header.Set("Content-Type", "application/json")      //
			//		gz, err := gzip.NewWriterLevel(respon, gzip.BestSpeed) // compressing
			gz := gzip.NewWriter(respon) // compressing
			// if err != nil {
			// 	io.WriteString(respon, err.Error())
			// 	return
			// }
			defer gz.Close()
			//respon.Header().Set("Content-Encoding", "gzip")
			rwr = gzipWriter{ResponseWriter: respon, Writer: gz}
		}
		if strings.Contains(claim.Header.Get("Content-Encoding"), "gzip") {
			respon.Header().Set("Content-Type", "application/json") // без этого в тестах -
			// Error Trace:    y:\GO\ypro\iteration8_test.go:183
			// 						y:\GO\ypro\suite.go:91
			// Error:          "application/octet-stream" does not contain "application/json"
			// Test:           TestIteration8/TestCounterGzipHandlers/update
			// Messages:       Заголовок ответа Content-Type содержит несоответствующее значение

			//			req.Header.Set("Content-Type", "application/json")      //
			//			rwr.Header().Set("Content-Encoding", "")  --------------
			//			req.Header.Set("Content-Encoding", "")   --------------

			gzipReader, err := gzip.NewReader(claim.Body) // decompressing
			if err != nil {
				io.WriteString(respon, err.Error())
				return
			}
			newReq, err := http.NewRequest(claim.Method, claim.URL.String(), gzipReader)
			if err != nil {
				io.WriteString(respon, err.Error())
				return
			}
			//rwr.Header().Set("Content-Type", "application/json")  //
			//newReq.Header.Set("Content-Type", "application/json") //+++++++++++++++++++
			//newReq.Header.Set("Accept", "application/json")       //

			req = newReq
		}

		next.ServeHTTP(rwr, req)
	})
}

/*
curl localhost:8087/update/ -H "Content-Type":"application/json" -d "{\"type\":\"gauge\",\"id\":\"nam\",\"value\":77}"


*/
