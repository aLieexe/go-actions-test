package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Envelope map[string]any

// func (app *Application) readIdParam(r *http.Request) (int64, error) {
// 	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
// 	if err != nil || id < 1 {
// 		return 0, errors.New("invalid id parameter")
// 	}

// 	return id, nil
// }

func (app *Application) WriteJSON(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	jsonData = append(jsonData, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(jsonData)
	return err
}

// dst should be an address btw
func (app *Application) ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at characted: %d)", syntaxError.Offset)

		// occur when the JSON value is the wrong type for the target destination
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field == "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// occur when we pass something that is not a non nil pointer to decode.
		// i believe the reason we panic is because its more of a dev error
		// rather than user err
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		// This should already be handled within the syntaxError check. However there is a chance
		// Decode return this error instead. so just for precaution
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formed JSON")

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown field: %s", fieldName)

		default:
			return err
		}
	}

	// remember dumass, struct{} is only the struct. U need another {} to initialize it
	// we essentially try to decode again, but here we expect io.EOF instead
	dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// func (app *Application) readString(queryString url.Values, key string, defaultVal string) string {
// 	s := queryString.Get(key)
// 	if s == "" {
// 		return defaultVal
// 	}

// 	return s
// }

// func (app *Application) readCSV(queryString url.Values, key string, defaultVal []string) []string {
// 	csv := queryString.Get(key)
// 	if csv == "" {
// 		return defaultVal
// 	}

// 	return strings.Split(csv, ",")
// }

func GetEnv(key string, defaultVal ...string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return ""
}
