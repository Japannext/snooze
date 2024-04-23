package field

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"

	api "github.com/japannext/snooze/common/api/v2"
)

// An object representing an alert's field (or a nested
// sub-field in case of maps).
// Useful for interacting with Alerts programmatically
// from a Pipeline (conditions/transforms/snooze),
// while considering the values to be only strings.
type AlertField struct {
	Name   string
	SubKey string
}

var allowedFields = []string{"severity_number", "severity_text"}
var allowedSubFields = []string{"group_labels", "labels", "attributes", "body"}

func (field *AlertField) Validate() error {
	if slices.Contains(allowedFields, field.Name) {
		if field.SubKey == "" {
			return nil
		}
		return fmt.Errorf("subkey not authorized with field '%s'", field.Name)
	}
	if slices.Contains(allowedSubFields, field.Name) {
		if field.SubKey == "" {
			return fmt.Errorf("subkey required with field '%s'", field.Name)
		}
		return nil
	}
	return fmt.Errorf("unknown/unauthorized field '%s'", field.Name)
}

func (field *AlertField) Get(alert *api.Alert) (string, bool) {
	var (
		v     string
		found bool
	)
	switch field.Name {
	case "severity_number":
		v = strconv.Itoa(int(alert.SeverityNumber))
		found = (alert.SeverityNumber != 0)
	case "severity_text":
		v = alert.SeverityText
		found = (alert.SeverityText != "")
	case "labels":
		v, found = alert.Labels[field.SubKey]
	case "attributes":
		v, found = alert.Attributes[field.SubKey]
	case "group_labels":
		v, found = alert.GroupLabels[field.SubKey]
	case "body":
		v, found = alert.Body[field.SubKey]
	default:
		v, found = "", false
	}
	return v, found
}

func (field *AlertField) Set(alert *api.Alert, v string) error {
	switch field.Name {
	case "severity_number":
		i, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		alert.SeverityNumber = int32(i)
	case "severity_text":
		alert.SeverityText = v
	case "labels":
		alert.Labels[field.SubKey] = v
	case "attributes":
		alert.Attributes[field.SubKey] = v
	case "group_labels":
		alert.GroupLabels[field.SubKey] = v
	case "body":
		alert.Body[field.SubKey] = v
	default:
		return fmt.Errorf("field '%s' not found", field.Name)
	}
	return nil
}

func (field *AlertField) Reset(alert *api.Alert) {
	switch field.Name {
	case "severity_number":
		alert.SeverityNumber = 0
	case "severity_text":
		alert.SeverityText = ""
	case "labels":
		delete(alert.Labels, field.SubKey)
	case "attributes":
		delete(alert.Attributes, field.SubKey)
	case "group_labels":
		delete(alert.GroupLabels, field.SubKey)
	case "body":
		delete(alert.Body, field.SubKey)
	default:
		// do nothing
	}
}

func (field *AlertField) String() string {
	var s strings.Builder
	s.WriteString(field.Name)
	if field.SubKey != "" {
		s.WriteString("[")
		s.WriteString(field.SubKey)
		s.WriteString("]")
	}
	return s.String()
}
