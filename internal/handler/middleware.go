package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"

	"github.com/labstack/echo/v4"
)

func LogRequestBodyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}

		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if c.Request().Header.Get("Content-Type") == "application/json" {
			var jsonObj interface{}
			err := json.Unmarshal(bodyBytes, &jsonObj)
			if err == nil {
				prettyJSON, err := json.MarshalIndent(jsonObj, "", "  ")
				if err == nil {
					log.Println("Request JSON:\n", string(prettyJSON))
				} else {
					log.Println("Request Body:", string(bodyBytes))
				}
			} else {
				log.Println("Request Body:", string(bodyBytes))
			}
		} else {
			log.Println("Request Body:", string(bodyBytes))
		}

		return next(c)
	}
}
