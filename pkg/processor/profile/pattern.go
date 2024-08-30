package profile

import (
	"context"
	"regexp"
	"fmt"

	api "github.com/japannext/snooze/pkg/common/api/v2"
	"github.com/japannext/snooze/pkg/common/lang"
)

type Pattern struct {
	// Name of the pattern
	Name string	`yaml:"name"`
	// If present, the pattern will match a given regex
	Regex string `yaml:"regex"`
	// List of labels/fields used to group the logs
	// for notification purposes
	GroupBy map[string]string `yaml:"group_by"`
	// Override the identity of the log
	IdentityOverride map[string]string `yaml:"identity_override"`
	// Drop labels from the log when the pattern match
	DroppedLabels []string `yaml:"dropped_labels"`
	// Add extra labels to the log when the pattern match
	ExtraLabels map[string]string `yaml:"extra_labels"`
	// Drop the log (dropped from database)
	Drop bool `yaml:"drop"`
	// Silence the log (won't notify)
	Silence bool `yaml:"silence"`
	// A template to override the message of the log
	Message string `yaml:"message"`

	// Internal values initialized after startup
	internal struct{
		regexp *regexp.Regexp
		identityOverride map[string]lang.Template
		groupBy map[string]lang.Template
		extraLabels map[string]lang.Template
		message *lang.Template
	}
}

func (p *Pattern) String() string {
	if p.Name != "" {
		return p.Name
	}
	return fmt.Sprintf("/%s/", p.Regex)
}

// Initialize internal values at startup
func (p *Pattern) Load() error {
	var err error
	p.internal.regexp, err = regexp.Compile(p.Regex)
	if err != nil {
		return err
	}
	p.internal.extraLabels, err = lang.NewTemplateMap(p.ExtraLabels)
	if err != nil {
		return err
	}
	p.internal.groupBy, err = lang.NewTemplateMap(p.GroupBy)
	if err != nil {
		return err
	}
	p.internal.identityOverride, err = lang.NewTemplateMap(p.IdentityOverride)
	if err != nil {
		return err
	}
	if p.Message != "" {
		p.internal.message, err = lang.NewTemplate(p.Message)
		if err != nil {
			return err
		}
	}
	log.Debugf("[Startup] %+v", p.internal)
	return nil
}

func (p *Pattern) Process(item *api.Log) (match, reject bool) {
	ctx := context.TODO()
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
		log.Debugf("Dropping the log")
		item.Mute.Drop(fmt.Sprintf("Silenced by pattern `%s`", p.Name))
		return
	}

	// Silence
	if p.Silence {
		log.Debugf("Silencing log")
		item.Mute.Silence(fmt.Sprintf("Silenced by pattern `%s`", p.Name))
	}

	// Extra labels
	log.Debugf(".ExtraLabels: %s", p.ExtraLabels)
	log.Debugf(".internal.extraLabels: %v", p.internal.extraLabels)
	for label, tpl := range p.internal.extraLabels {
		value, err := tpl.Execute(ctx, item)
		if err != nil {
			log.Warnf("failed to execute template `%s`", p.ExtraLabels[label])
			return
		}
		log.Debugf("Adding extra label: %s=%s", label, value)
		item.Labels[label] = value
	}

	// Group By
	var groupByLabels = map[string]string{}
	for label, tpl := range p.internal.groupBy {
		value, err := tpl.Execute(ctx, item)
		if err != nil {
			log.Warnf("failed to execute template `%s`", p.GroupBy[label])
			return
		}
		log.Debugf("Adding groupBy: %s=%s", label, value)
		groupByLabels[label] = value
	}
	item.Group.Labels = groupByLabels

	// Identity override
	if len(p.internal.identityOverride) > 0 {
		identity := make(map[string]string)
		for key, tpl := range p.internal.identityOverride {
			value, err := tpl.Execute(ctx, item)
			if err != nil {
				log.Warnf("failed to execute template `%s`", p.IdentityOverride[key])
				return
			}
			identity[key] = value
		}
		log.Debugf("Overriding identity: %s", identity)
		item.Identity = identity
	}

	// Dropped labels
	for _, label := range p.DroppedLabels {
		log.Debugf("Dropping label `%s`", label)
		delete(item.Labels, label)
	}

	// Message override
	if p.internal.message != nil {
		var err error
		if err != nil {
			log.Warnf("failed to execute template `%s`", p.Message)
			return
		}
		msg, err := p.internal.message.Execute(ctx, item)
		if err != nil {
			return
		}
		log.Debugf("Changing message to `%s`", msg)
		item.Message = msg
	}

	return
}

// Match the regex of the pattern, return the capture groups if any
func (p *Pattern) match(item *api.Log) (bool, map[string]string) {
	var (
		match bool
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
