package profile

import (
	"context"
	"fmt"
	"regexp"

	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/processor/transform"
)

type Pattern struct {
	// Name of the pattern
	Name string `json:"name" yaml:"name"`
	// Description of the pattern
	Description string `json:"description" yaml:"description"`
	// If present, the pattern will match a given regex
	Regex string `json:"regex" yaml:"regex"`

	Actions []transform.Action `json:"actions,omitempty" yaml:"actions"`

	// List of labels/fields used to group the logs
	// for notification purposes
	GroupBy map[string]string `json:"groupBy" yaml:"group_by"`

	// Drop labels from the log when the pattern match
	DroppedLabels []string `json:"droppedLabels" yaml:"dropped_labels"`

	// Drop the log (dropped from database)
	Drop bool `json:"drop" yaml:"drop"`

	// Silence the log (won't notify)
	Silence bool `json:"silence" yaml:"silence"`

	// Internal values initialized after startup
	internal struct {
		regexp  *regexp.Regexp
		actions []transform.Transformation
		groupBy map[string]*lang.Template
	}
}

func (p *Pattern) String() string {
	if p.Name != "" {
		return p.Name
	}
	return fmt.Sprintf("/%s/", p.Regex)
}

// Initialize internal values at startup.
func (p *Pattern) Load() error {
	var err error

	p.internal.regexp, err = regexp.Compile(p.Regex)
	if err != nil {
		return err
	}

	p.internal.actions = transform.LoadActions(p.Actions)
	p.internal.groupBy, err = lang.NewTemplateMap(p.GroupBy)
	if err != nil {
		return err
	}

	log.Debugf("[Startup] %+v", p.internal)
	return nil
}

func (p *Pattern) Process(ctx context.Context, item *models.Log) (match, reject bool) {
	// Matching pattern
	match, capture := p.match(item)
	if !match {
		log.Debugf("Didn't match %s", p.Name)
		return
	}
	ctx = context.WithValue(ctx, "capture", capture)
	log.Debugf("Matched pattern %s", p.String())

	// Drop
	if p.Drop {
		if ok := item.Status.Change(models.LogDropped); ok {
			item.Status.Reason = fmt.Sprintf("dropped by pattern '%s'", p.Name)
			item.Status.SkipNotification = true
			item.Status.SkipStorage = true
		}
		return
	}

	// Silence
	if p.Silence {
		if ok := item.Status.Change(models.LogSilenced); ok {
			item.Status.Reason = fmt.Sprintf("silenced by pattern '%s'", p.Name)
			item.Status.SkipNotification = true
		}
	}

	for _, action := range p.internal.actions {
		var err error
		ctx, err = action.Process(ctx, item)
		if err != nil {
			return
		}
	}

	// Group By
	if len(p.internal.groupBy) > 0 {
		gr := &models.Group{Name: p.Name, Labels: make(map[string]string)}
		for key, tpl := range p.internal.groupBy {
			value, err := tpl.Execute(ctx, item)
			if err != nil {
				log.Warnf("failed to execute template `%s`", p.GroupBy[key])
				return
			}
			log.Debugf("Adding groupBy: %s=%s", key, value)
			gr.Labels[key] = value
		}
		gr.Hash = utils.ComputeHash(gr.Labels)
		item.Groups = append(item.Groups, gr)
	}
	// Dropped labels
	for _, label := range p.DroppedLabels {
		log.Debugf("Dropping label `%s`", label)
		delete(item.Labels, label)
	}
	return
}

// Match the regex of the pattern, return the capture groups if any.
func (p *Pattern) match(item *models.Log) (bool, map[string]string) {
	var (
		match   bool
		capture = make(map[string]string)
	)
	if p.internal.regexp == nil {
		return true, capture
	}
	match = p.internal.regexp.MatchString(item.Message)
	if !match {
		return false, capture
	}
	keys := p.internal.regexp.SubexpNames()
	if len(keys) > 1 {
		keys = keys[1:]
		values := p.internal.regexp.FindStringSubmatch(item.Message)
		for _, key := range keys {
			i := p.internal.regexp.SubexpIndex(key)
			if i < 0 {
				continue
			}
			capture[key] = values[i]
		}
	}
	return true, capture
}
