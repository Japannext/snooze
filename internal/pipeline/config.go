package pipeline

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/schedule"
)

type Pipeline struct {
	Mappings            []Mapping            `json:"mappings" yaml:"mappings"`
	Groupings           []Grouping           `json:"groupings" yaml:"groupings"`
	Transforms          []Transform          `json:"transforms" yaml:"transforms"`
	Services            []Service            `json:"services" yaml:"services"`
	Profiles            []Profile            `json:"profiles" yaml:"profiles"`
	Silences            []Silence            `json:"silences" yaml:"silences"`
	Ratelimits          []Ratelimit          `json:"ratelimits" yaml:"ratelimits"`
	Notifications       []Notification       `json:"notifications" yaml:"notifications"`
	DefaultDestinations []models.Destination `json:"defaultDestinations" yaml:"default_destinations"`
}

// A key-value object used as data in rules/conditions/templates
type Mapping struct {
	// Name of the mapping
	Name string `json:"name" yaml:"name"`
	// The actual mapping
	Map map[string]string `json:"map" yaml:"map"`
}

type Field struct {
	// Name of the field. Will appear in the web interface.
	Name string `json:"name"`
	// The list of nested fields to access the field. Used in opensearch
	// queries to access the data.
	NestedFields []string `json:"nestedFields" yaml:"nested_fields"`
}

type Grouping struct {
	// Name of the grouping
	Name string `json:"name" yaml:"name"`
	// Condition for applying the grouping only to certain entries
	If string `json:"if"   yaml:"if"`
	// Mutually exclusive with `group_by_map`.
	GroupBy []Field `json:"groupBy" yaml:"group_by"`
}

type Transform struct {
	Name    string   `json:"name" yaml:"name"`
	If      string   `json:"if,omitempty" yaml:"if"`
	Actions []Action `json:"actions" yaml:"actions"`
}

type Action struct {
	Set   *Set   `json:"set,omitempty"   yaml:"set"`
	Unset *Unset `json:"unset,omitempty" yaml:"unset"`
	Regex *Regex `json:"regex,omitempty" yaml:"regex"`
}

type Set struct {
	// Set one or more identity field to a given value
	Identity map[string]string `json:"identity,omitempty" yaml:"identity"`
	// Set one or more label field to a given value
	Labels map[string]string `json:"labels,omitempty" yaml:"labels"`
	// Set a message to a given value
	Message string `json:"message,omitempty" yaml:"message"`
}

type Unset struct {
	Identity    []string `yaml:""`
	Labels      []string `json:"labels,omitempty"      yaml:"labels"`
	AllIdentity bool     `json:"allIdentity,omitempty" yaml:"all_identity"`
	AllLabels   bool     `json:"allLabels,omitempty"   yaml:"all_labels"`
}

type Regex struct {
	Field string `json:"field" yaml:"field"`
	Match string `json:"match" yaml:"match"`
}

type Profile struct {
	// Name of the profile group
	Name string `json:"name" yaml:"name"`
	// The main condition for a log to match this rule. Used to
	// reduce the amount of processing by the use of maps.
	// Examples: process=sshd, service.name=keycloak, k8s.statefulset.name=postgresql
	Selector Selector `json:"selector" yaml:"selector"`
	// Patterns and actions to apply to logs matching this pattern
	Patterns []Pattern `json:"patterns" yaml:"patterns"`
}

func (p Profile) GetSelectors() []Selector {
	return []Selector{p.Selector}
}

// A selector for a label.
type Selector struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

type Pattern struct {
	// Name of the pattern
	Name string `json:"name" yaml:"name"`
	// Description of the pattern
	Description string `json:"description" yaml:"description"`
	// If present, the pattern will match a given regex
	Regex string `json:"regex" yaml:"regex"`

	Actions []Action `json:"actions,omitempty" yaml:"actions"`

	// List of labels/fields used to group the logs
	// for notification purposes
	GroupBy map[string]string `json:"groupBy" yaml:"group_by"`

	// Drop labels from the log when the pattern match
	DroppedLabels []string `json:"droppedLabels" yaml:"dropped_labels"`

	// Drop the log (dropped from database)
	Drop bool `json:"drop" yaml:"drop"`

	// Silence the log (won't notify)
	Silence bool `json:"silence" yaml:"silence"`
}

type Silence struct {
	// Name of the silence rule
	Name string `json:"name" yaml:"name"`
	// Description of the silence rule
	Description string `json:"description,omitempty" yaml:"description"`
	// Condition to execute the silence rule
	If string `json:"if" yaml:"if"`
	// If the silence rule is on schedule
	Schedule *schedule.Schedule `json:",inline" yaml:",inline"`
	// Whether the logs should be dropped or kept if they match the
	// silence rule
	Drop bool `json:"drop" yaml:"drop"`
}

type Ratelimit struct {
	// Name of the rate limit. Must be unique (used for key)
	Name string `json:"name" yaml:"name"`
	// An optional condition to only apply the rate limit if the condition match
	If string `json:"if,omitempty" yaml:"if"`
	// The group to group by
	Group string `json:"group" yaml:"group"`
	// The amount of authorized logs during the period
	CountPerPeriod int64 `json:"countPerPeriod" yaml:"countPerPeriod"`
	// The period used to measure the rate
	Period time.Duration `json:"period" yaml:"period"`
}

type Notification struct {
	// Name of the notification rule
	Name string `json:"name" yaml:"name"`
	// Condition to use for matching the notification
	If string `json:"if,omitempty" yaml:"if"`
	// Which queue the notification will send to
	Destinations []models.Destination `json:"destinations" yaml:"destinations"`
}

// A Service represent a subsystem that is monitored and should group alerts.
// It can group specialized alerts such as SLOs, and group logs from specific servers.
type Service struct {
	// Name of the service.
	Name string `json:"name" yaml:"name"`
	// Description of the service.
	Description string `json:"description,omitempty" yaml:"description"`
	// The owner of the service. Must be a team.
	Owner string `json:"owner,omitempty" yaml:"owner"`
	// List of impacts.
	Impacts []string `json:"impacts" yaml:"impacts"`
	// A list of selectors used to know what log/alert match the service.
	Selectors []Selector `json:"selectors" yaml:"selectors"`
}

func (s Service) GetSelectors() []Selector {
	return s.Selectors
}

// An impact is a way to categorize alerts and services in groups to make
// it easier to understand.
// Example of categories: "production", "office", "devoffice", "monitoring".
type Impact struct {
	// Name of the impact
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description"`
}

// An object representing a group of people.
type Team struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description"`
}

var (
	pl     *Pipeline
	plOnce sync.Once
)

func load() (*Pipeline, error) {
	v := viper.New()
	v.SetDefault("PIPELINE_FILE", "/etc/snooze/pipeline.yaml")
	v.AutomaticEnv()

	pipelineFile := v.GetString("PIPELINE_FILE")

	var p *Pipeline
	data, err := os.ReadFile(pipelineFile)
	if err != nil {
		return &Pipeline{}, fmt.Errorf("failed to read %s: %s", pipelineFile, err)
	}

	if err := yaml.Unmarshal(data, &p); err != nil {
		return &Pipeline{}, fmt.Errorf("failed to unmarshal yaml file %s: %s", pipelineFile, err)
	}

	return p, nil
}

func Reload() (*Pipeline, error) {
	var err error
	pl, err = load()
	return pl, err
}

// Initialize the pipeline, and return it once
func Get() (*Pipeline, error) {
	var rerr error
	plOnce.Do(func() {
		pl, rerr = load()
	})

	return pl, rerr
}
