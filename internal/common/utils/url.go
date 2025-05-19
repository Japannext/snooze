package utils

import (
	"net/url"
	"reflect"

	"github.com/spf13/viper"
)

func stringToURL(f, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}
	if t != reflect.TypeOf(url.URL{}) {
		return data, nil
	}
	// Convert it by parsing
	return url.Parse(data.(string))
}

func URLDecodeHook() viper.DecoderConfigOption {
	return viper.DecodeHook(stringToURL)
}
