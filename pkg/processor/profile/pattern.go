package profile

import (
	"regexp"
	"fmt"

	api "github.com/japannext/snooze/pkg/common/api/v2"
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
		identityOverride map[string]Template
		groupBy map[string]Template
		extraLabels map[string]Template
		message *Template
	}
}

// Initialize internal values at startup
func (p *Pattern) Startup() (err error) {
	p.internal.regexp, err = regexp.Compile(p.Regex)
	if err != nil {
		return
	}
	p.internal.extraLabels, err = NewTemplateMap(p.ExtraLabels)
	if err != nil {
		return
	}
	p.internal.groupBy, err = NewTemplateMap(p.GroupBy)
	if err != nil {
		return
	}
	p.internal.identityOverride, err = NewTemplateMap(p.IdentityOverride)
	if err != nil {
		return
	}
	if p.Message != "" {
		p.internal.message, err = NewTemplate(p.Message)
		if err != nil {
			return
		}
	}
	return
}

func (p *Pattern) Process(item *api.Log) (match, reject bool) {
	match, capture := p.match(item)
	if !match {
		log.Debugf("Didn't match %s", p.Name)
		return
	}
	log.Debugf("Matched pattern %s", p.Name)
	if p.Drop {
		log.Debugf("Dropping the log")
		item.Mute.Drop(fmt.Sprintf("Silenced by pattern `%s`", p.Name))
		return
	}
	if p.Silence {
		log.Debugf("Silencing log")
		item.Mute.Silence(fmt.Sprintf("Silenced by pattern `%s`", p.Name))
	}
	log.Debugf(".ExtraLabels: %s", p.ExtraLabels)
	log.Debugf(".internal.extraLabels: %v", p.internal.extraLabels)
	for label, tpl := range p.internal.extraLabels {
		value, err := tpl.Execute(item, capture)
		if err != nil {
			log.Warnf("failed to execute template `%s`", tpl)
			return
		}
		log.Debugf("Adding extra label: %s=%s", label, value)
		item.Labels[label] = value
	}
	for label, tpl := range p.internal.groupBy {
		value, err := tpl.Execute(item, capture)
		if err != nil {
			log.Warnf("failed to execute template `%s`", tpl)
			return
		}
		log.Debugf("Adding groupBy: %s=%s", label, value)
		item.GroupLabels[label] = value
	}
	if len(p.internal.identityOverride) > 0 {
		identity := make(map[string]string)
		for key, tpl := range p.internal.identityOverride {
			value, err := tpl.Execute(item, capture)
			if err != nil {
				log.Warnf("failed to execute template `%s`", tpl)
				return
			}
			identity[key] = value
		}
		log.Debugf("Overriding identity: %s", identity)
		item.Identity = identity
	}
	for _, label := range p.DroppedLabels {
		log.Debugf("Dropping label `%s`", label)
		delete(item.Labels, label)
	}
	if p.internal.message != nil {
		var err error
		if err != nil {
			log.Warnf("failed to execute template `%s`", p.internal.message)
			return
		}
		msg, err := p.internal.message.Execute(item, capture)
		if err != nil {
			return
		}
		log.Debugf("Changing message to `%s`", msg)
		item.Message = msg
	}
	return
}

// Match the regex of the pattern, return the capture groups if any
func (p *Pattern) match(item *api.Log) (match bool, capture map[string]string) {
	if p.internal.regexp == nil {
		match = true
		return
	}
	match = p.internal.regexp.MatchString(item.Message)
	if !match {
		return
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
