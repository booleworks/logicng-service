package sio

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Unmarshal[T ServiceInput[T]](r *http.Request) (object T, sErr ServiceError) {
	switch ct := r.Header.Get("Content-Type"); ct {
	case "", "application/json":
		err := json.NewDecoder(r.Body).Decode(&object)
		if err != nil {
			sErr = ErrIllegalInput(err)
		}
	case "application/protobuf":
		data, err := io.ReadAll(r.Body)
		if err != nil {
			sErr = ErrIllegalInput(err)
			return
		}
		object, err = object.DeserProtoBuf(data)
		if err != nil {
			sErr = ErrIllegalInput(err)
		}
	default:
		sErr = ErrUnsupportedContentType(ct)
	}
	if valErrs := object.Validate(); valErrs != nil {
		sErr = ErrIllegalInput(fmt.Errorf(formatValidationErrors(valErrs)))
	}
	return
}

func WriteResult[T ServiceOutput[T]](w http.ResponseWriter, r *http.Request, object T) {
	var err error
	switch acc := r.Header.Get("accept"); acc {
	case "", "*/*", "application/json":
		w.Header().Add("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(object)
	case "application/protobuf":
		w.Header().Add("Content-Type", "application/protobuf")
		var data []byte
		data, err = object.ProtoBuf()
		if err == nil {
			w.Write(data)
		}
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("Unsupported accept type"))
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}
}

func formatValidationErrors(errors map[string]string) string {
	var sb strings.Builder
	sb.WriteString("validation errors ")
	size := len(errors)
	count := 0
	for k, v := range errors {
		count++
		sb.WriteString(fmt.Sprintf("[field '%s': %s]", k, v))
		if count < size {
			sb.WriteString(" ")
		}
	}
	return sb.String()
}
