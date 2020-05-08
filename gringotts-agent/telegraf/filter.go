package telegraf

import (
	"fmt"
	"strings"

	"github.com/gobwas/glob"
	"github.com/jinlingan/gringotts/gringotts-agent/model"
)

// TagFilter is the name of a tag, and the values on which to filter
type TagFilter struct {
	Name   string
	Filter []string
	filter TelegrafFilterInf
}

// Filter containing drop/pass and tagdrop/tagpass rules
type Filter struct {
	NameDrop []string
	nameDrop TelegrafFilterInf
	NamePass []string
	namePass TelegrafFilterInf

	FieldDrop []string
	fieldDrop TelegrafFilterInf
	FieldPass []string
	fieldPass TelegrafFilterInf

	TagDrop []TagFilter
	TagPass []TagFilter

	TagExclude []string
	tagExclude TelegrafFilterInf
	TagInclude []string
	tagInclude TelegrafFilterInf

	isActive bool
}

// Compile all Filter lists into TelegrafFilter objects.
func (f *Filter) Compile() error {
	if len(f.NameDrop) == 0 &&
		len(f.NamePass) == 0 &&
		len(f.FieldDrop) == 0 &&
		len(f.FieldPass) == 0 &&
		len(f.TagInclude) == 0 &&
		len(f.TagExclude) == 0 &&
		len(f.TagPass) == 0 &&
		len(f.TagDrop) == 0 {
		return nil
	}

	f.isActive = true
	var err error
	f.nameDrop, err = Compile(f.NameDrop)
	if err != nil {
		return fmt.Errorf("Error compiling 'namedrop', %s", err)
	}
	f.namePass, err = Compile(f.NamePass)
	if err != nil {
		return fmt.Errorf("Error compiling 'namepass', %s", err)
	}

	f.fieldDrop, err = Compile(f.FieldDrop)
	if err != nil {
		return fmt.Errorf("Error compiling 'fielddrop', %s", err)
	}
	f.fieldPass, err = Compile(f.FieldPass)
	if err != nil {
		return fmt.Errorf("Error compiling 'fieldpass', %s", err)
	}

	f.tagExclude, err = Compile(f.TagExclude)
	if err != nil {
		return fmt.Errorf("Error compiling 'tagexclude', %s", err)
	}
	f.tagInclude, err = Compile(f.TagInclude)
	if err != nil {
		return fmt.Errorf("Error compiling 'taginclude', %s", err)
	}

	for i := range f.TagDrop {
		f.TagDrop[i].filter, err = Compile(f.TagDrop[i].Filter)
		if err != nil {
			return fmt.Errorf("Error compiling 'tagdrop', %s", err)
		}
	}
	for i := range f.TagPass {
		f.TagPass[i].filter, err = Compile(f.TagPass[i].Filter)
		if err != nil {
			return fmt.Errorf("Error compiling 'tagpass', %s", err)
		}
	}
	return nil
}

// Select returns true if the metric matches according to the
// namepass/namedrop and tagpass/tagdrop filters.  The metric is not modified.
func (f *Filter) Select(metric model.Metric) bool {
	if !f.isActive {
		return true
	}

	if !f.shouldNamePass(metric.Name()) {
		return false
	}

	if !f.shouldTagsPass(metric.TagList()) {
		return false
	}

	return true
}

// Modify removes any tags and fields from the metric according to the
// fieldpass/fielddrop and taginclude/tagexclude filters.
func (f *Filter) Modify(metric model.Metric) {
	if !f.isActive {
		return
	}

	f.filterFields(metric)
	f.filterTags(metric)
}

// IsActive checking if filter is active
func (f *Filter) IsActive() bool {
	return f.isActive
}

// shouldNamePass returns true if the metric should pass, false if should drop
// based on the drop/pass filter parameters
func (f *Filter) shouldNamePass(key string) bool {
	pass := func(f *Filter) bool {
		if f.namePass.Match(key) {
			return true
		}
		return false
	}

	drop := func(f *Filter) bool {
		if f.nameDrop.Match(key) {
			return false
		}
		return true
	}

	if f.namePass != nil && f.nameDrop != nil {
		return pass(f) && drop(f)
	} else if f.namePass != nil {
		return pass(f)
	} else if f.nameDrop != nil {
		return drop(f)
	}

	return true
}

// shouldFieldPass returns true if the metric should pass, false if should drop
// based on the drop/pass filter parameters
func (f *Filter) shouldFieldPass(key string) bool {
	if f.fieldPass != nil && f.fieldDrop != nil {
		return f.fieldPass.Match(key) && !f.fieldDrop.Match(key)
	} else if f.fieldPass != nil {
		return f.fieldPass.Match(key)
	} else if f.fieldDrop != nil {
		return !f.fieldDrop.Match(key)
	}
	return true
}

// shouldTagsPass returns true if the metric should pass, false if should drop
// based on the tagdrop/tagpass filter parameters
func (f *Filter) shouldTagsPass(tags []*model.Tag) bool {
	pass := func(f *Filter) bool {
		for _, pat := range f.TagPass {
			if pat.filter == nil {
				continue
			}
			for _, tag := range tags {
				if tag.Key == pat.Name {
					if pat.filter.Match(tag.Value) {
						return true
					}
				}
			}
		}
		return false
	}

	drop := func(f *Filter) bool {
		for _, pat := range f.TagDrop {
			if pat.filter == nil {
				continue
			}
			for _, tag := range tags {
				if tag.Key == pat.Name {
					if pat.filter.Match(tag.Value) {
						return false
					}
				}
			}
		}
		return true
	}

	// Add additional logic in case where both parameters are set.
	// see: https://github.com/influxdata/telegraf/issues/2860
	if f.TagPass != nil && f.TagDrop != nil {
		// return true only in case when tag pass and won't be dropped (true, true).
		// in case when the same tag should be passed and dropped it will be dropped (true, false).
		return pass(f) && drop(f)
	} else if f.TagPass != nil {
		return pass(f)
	} else if f.TagDrop != nil {
		return drop(f)
	}

	return true
}

// filterFields removes fields according to fieldpass/fielddrop.
func (f *Filter) filterFields(metric model.Metric) {
	filterKeys := []string{}
	for _, field := range metric.FieldList() {
		if !f.shouldFieldPass(field.Key) {
			filterKeys = append(filterKeys, field.Key)
		}
	}

	for _, key := range filterKeys {
		metric.RemoveField(key)
	}
}

// filterTags removes tags according to taginclude/tagexclude.
func (f *Filter) filterTags(metric model.Metric) {
	filterKeys := []string{}
	if f.tagInclude != nil {
		for _, tag := range metric.TagList() {
			if !f.tagInclude.Match(tag.Key) {
				filterKeys = append(filterKeys, tag.Key)
			}
		}
	}
	for _, key := range filterKeys {
		metric.RemoveTag(key)
	}

	if f.tagExclude != nil {
		for _, tag := range metric.TagList() {
			if f.tagExclude.Match(tag.Key) {
				filterKeys = append(filterKeys, tag.Key)
			}
		}
	}
	for _, key := range filterKeys {
		metric.RemoveTag(key)
	}
}

type TelegrafFilterInf interface {
	Match(string) bool
}

// Compile takes a list of string filters and returns a Filter interface
// for matching a given string against the filter list. The filter list
// supports glob matching too, ie:
//
//   f, _ := Compile([]string{"cpu", "mem", "net*"})
//   f.Match("cpu")     // true
//   f.Match("network") // true
//   f.Match("memory")  // false
//
func Compile(filters []string) (TelegrafFilterInf, error) {
	// return if there is nothing to compile
	if len(filters) == 0 {
		return nil, nil
	}

	// check if we can compile a non-glob filter
	noGlob := true
	for _, filter := range filters {
		if hasMeta(filter) {
			noGlob = false
			break
		}
	}

	switch {
	case noGlob:
		// return non-globbing filter if not needed.
		return compileFilterNoGlob(filters), nil
	case len(filters) == 1:
		return glob.Compile(filters[0])
	default:
		return glob.Compile("{" + strings.Join(filters, ",") + "}")
	}
}

// hasMeta reports whether path contains any magic glob characters.
func hasMeta(s string) bool {
	return strings.IndexAny(s, "*?[") >= 0
}

type filter struct {
	m map[string]struct{}
}

func (f *filter) Match(s string) bool {
	_, ok := f.m[s]
	return ok
}

type filtersingle struct {
	s string
}

func (f *filtersingle) Match(s string) bool {
	return f.s == s
}

func compileFilterNoGlob(filters []string) TelegrafFilterInf {
	if len(filters) == 1 {
		return &filtersingle{s: filters[0]}
	}
	out := filter{m: make(map[string]struct{})}
	for _, filter := range filters {
		out.m[filter] = struct{}{}
	}
	return &out
}

type IncludeExcludeFilter struct {
	include TelegrafFilterInf
	exclude TelegrafFilterInf
}

func NewIncludeExcludeFilter(
	include []string,
	exclude []string,
) (TelegrafFilterInf, error) {
	in, err := Compile(include)
	if err != nil {
		return nil, err
	}

	ex, err := Compile(exclude)
	if err != nil {
		return nil, err
	}

	return &IncludeExcludeFilter{in, ex}, nil
}

func (f *IncludeExcludeFilter) Match(s string) bool {
	if f.include != nil {
		if !f.include.Match(s) {
			return false
		}
	}

	if f.exclude != nil {
		if f.exclude.Match(s) {
			return false
		}
	}
	return true
}
