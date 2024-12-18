package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/mdw-smarty/calc-lib/calc"
)

func NewHTTPRouter() http.Handler {
	h := http.NewServeMux()
	h.Handle("/add", NewHTTPHandler(calc.Addition{}))
	h.Handle("/sub", NewHTTPHandler(calc.Subtraction{}))
	h.Handle("/mul", NewHTTPHandler(calc.Multiplication{}))
	h.Handle("/div", NewHTTPHandler(calc.Division{}))
	h.Handle("/bog", NewHTTPHandler(calc.Bogus{Offset: 42}))
	return h
}

type HTTPHandler struct {
	calculator Calculator
}

func NewHTTPHandler(calculator Calculator) http.Handler {
	return &HTTPHandler{calculator: calculator}
}

func (this *HTTPHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	rawA := query.Get("a")
	a, err := strconv.Atoi(rawA)
	if err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(writer, "invalid 'a' parameter: [%s]", rawA)
		return
	}

	rawB := query.Get("b")
	b, err := strconv.Atoi(rawB)
	if err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(writer, "invalid 'b' parameter: [%s]", rawB)
		return
	}

	writer.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintln(writer, this.calculator.Calculate(a, b))
	if err != nil {
		log.Println("Failed to write response:", err)
	}
}
