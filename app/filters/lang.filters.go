package filters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/labstack/echo/v4"
)

var translations map[string]map[string]string

func LoadTranslations() error {
	translations = make(map[string]map[string]string)

	langs := []string{"en", "id"}
	for _, lang := range langs {
		data, err := ioutil.ReadFile(fmt.Sprintf("app/views/js/src/locale/%s.json", lang))
		if err != nil{
			log.Fatal("Error load file lang", err.Error())
			return err
		}

		var m map[string]string
		if err := json.Unmarshal(data, &m); err != nil{
			log.Fatal("Error parse json", err.Error())
			return err
		}

		translations[lang] = m
	}

	return nil
}

func LangFilters(next echo.HandlerFunc)echo.HandlerFunc{
	return func(c echo.Context) error {
		lang := c.Request().Header.Get("Accept-Lang")
		if lang == "" || translations[lang] == nil{
			lang = "id"
		}
		c.Set("lang", lang)
		return next(c)
	}
}

func Translate(c echo.Context, key string)string{
	lang := c.Get("lang").(string)
	if val, ok := translations[lang][key]; ok {
		return val
	}
	return translations["id"][key]
}