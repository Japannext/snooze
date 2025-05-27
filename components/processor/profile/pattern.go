package profile

import (
	"context"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/components/processor/transform"
	"github.com/japannext/snooze/internal/lang"
	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/pipeline"
	"github.com/japannext/snooze/internal/utils"
)

type Pattern struct {
	pipeline.Pattern

	regexp  *regexp.Regexp
	actions []transform.Action
	groupBy map[string]*lang.Template
}

func (p *Pattern) String() string {
	if p.Name != "" {
		return p.Name
	}
	return fmt.Sprintf("/%s/", p.Regex)
}

// Initialize internal values at startup.
func NewPattern(cfg pipeline.Pattern) (Pattern, error) {
	var err error

	pattern := Pattern{}

	pattern.regexp, err = regexp.Compile(cfg.Regex)
	if err != nil {
		return Pattern{}, fmt.Errorf("invalid regex `%s`: %w", cfg.Regex, err)
	}

	for index, actionConfig := range cfg.Actions {
		action, err := transform.NewAction(actionConfig)
		if err != nil {
			return Pattern{}, fmt.Errorf("in action #%d: %w", index+1, err)
		}
		pattern.actions = append(pattern.actions, action)
	}

	pattern.groupBy, err = lang.NewTemplateMap(cfg.GroupBy)
	if err != nil {
		return Pattern{}, fmt.Errorf("in groupBy: %w", err)
	}

	return pattern, nil
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

	for _, action := range p.actions {
		var err error
		ctx, err = action.Process(ctx, item)
		if err != nil {
			return
		}
	}

	// Group By
	if len(p.groupBy) > 0 {
		gr := &models.Group{Name: p.Name, Labels: make(map[string]string)}
		for key, tpl := range p.groupBy {
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

	if p.regexp == nil {
		return true, capture
	}

	match = p.regexp.MatchString(item.Message)
	if !match {
		return false, capture
	}

	keys := p.regexp.SubexpNames()
	if len(keys) > 1 {
		keys = keys[1:]
		values := p.regexp.FindStringSubmatch(item.Message)

		for _, key := range keys {
			i := p.regexp.SubexpIndex(key)
			if i < 0 {
				continue
			}
			capture[key] = values[i]
		}
	}
	return true, capture
}
