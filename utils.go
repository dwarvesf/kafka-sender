package main

import (
	"fmt"
	"strconv"
	"strings"
)

func parseMap(aMap map[string]interface{}, prefix string, b *strings.Builder) {
	var sol string
	for key, val := range aMap {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			parseMap(concreteVal, prefix+fmt.Sprintf(`"%v"`, key)+`:{`, b)
		case string:
			sol = fmt.Sprintf(`"%v"`, key) + `:"` + concreteVal + `"`
			b.WriteString(prefix)
			b.WriteString(sol)
			b.WriteRune(',')
		case float64:
			sol = fmt.Sprintf(`"%v"`, key) + ":" + strconv.FormatFloat(concreteVal, 'f', -1, 64)
			b.WriteString(prefix)
			b.WriteString(sol)
			b.WriteRune(',')
		case bool:
			sol = fmt.Sprintf(`"%v"`, key) + ":" + strconv.FormatBool(concreteVal)
			b.WriteString(prefix)
			b.WriteString(sol)
			b.WriteRune(',')
		}
	}
}

func parsecustom(aMap map[string]interface{}) string {
	b := &strings.Builder{}
	parseMap(aMap, "", b)

	return strings.TrimSuffix(b.String(), ",")
}
