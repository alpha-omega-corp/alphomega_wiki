package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/uptrace/bunrouter"
)

type GetOne struct {
	Id int64 `json:"id"`
}

func GetBody[T any](w http.ResponseWriter, req bunrouter.Request) *T {
	data := new(T)

	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		Error(w, err, http.StatusBadRequest)
	}

	return data
}

func GetParams[T any](w http.ResponseWriter, req bunrouter.Request) *T {
	paramsMap := req.Params().Map()
	paramsData := make(map[string]interface{}, len(paramsMap))

	for key, value := range paramsMap {
		paramsData[key] = value
	}

	if idStr, ok := paramsData["id"].(string); ok {
		if idInt, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			paramsData["id"] = idInt
		}
	}

	params, err := json.Marshal(paramsData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	data := new(T)
	if err := json.Unmarshal(params, data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	return data
}

func Response[T any](w http.ResponseWriter, req func() (*T, error)) error {
	res, err := req()
	return JSON(w, res, err)
}

func JSON[T any](w http.ResponseWriter, res *T, err error) error {
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
	}

	return bunrouter.JSON(w, res)
}

func Error(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)

	log.Printf("%+v\n", err)
	_, _ = w.Write([]byte(err.Error()))
}
