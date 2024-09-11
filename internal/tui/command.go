package tui

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elseano/dollop/internal/config"
	"github.com/elseano/dollop/internal/templating"
)

type scanMsg struct {
	lines  []list.Item
	status string
}

type disconnectedMsg struct{}

var scanMutex = sync.Mutex{}

func (m Model) scanLogs() tea.Cmd {
	return func() tea.Msg {
		scanMutex.Lock()
		status, err := processLog(m.config)

		if err != nil {
			return disconnectedMsg{}
		}

		disp := []list.Item{}

		for _, v := range itemsCache {
			disp = append(disp, v)
		}

		sort.Slice(disp, func(i, j int) bool {
			diff := disp[i].(*logGroup).timestamp.Sub(disp[j].(*logGroup).timestamp)

			if diff == 0 {
				return strings.Compare(disp[i].(*logGroup).title, disp[j].(*logGroup).title) > 0
			} else {
				return diff > 0
			}
		})

		scanMutex.Unlock()

		return scanMsg{lines: disp, status: status}
	}
}

var itemsCache = map[string]list.Item{
	"errors": &logGroup{
		title:       "Parse Failures",
		description: "Errors",
		timestamp:   time.Now(),
	},
	"not-json": &logGroup{
		title:       "Text",
		description: "Not JSON",
		timestamp:   time.Now(),
	},
}

var scanner *bufio.Reader

func processLog(config config.Config) (status string, err error) {
	if scanner == nil {
		scanner = bufio.NewReaderSize(os.Stdin, 1048576)
	}

	line, err := scanner.ReadString('\n')

	if err != nil {
		return "", err
	}

	jsonIdx := strings.Index(line, "{")

	if jsonIdx == -1 {
		tcache := itemsCache["not-json"].(*logGroup)
		tcache.lines = append(tcache.lines, logLine{
			message: line,
			data:    nil,
		})

		return
	}

	var res map[string]interface{}
	err = json.Unmarshal([]byte(line[jsonIdx:]), &res)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			tcache := itemsCache["errors"].(*logGroup)
			tcache.lines = append(tcache.lines, logLine{
				message: fmt.Sprintf("Error loading '%s': %s", line, err.Error()),
				data:    res,
			})
		}

		return
	}

	groupValue, groupTitle, groupSpec := getGroupAndTitle(config, res)
	timestamp := getTimestamp(config, res)
	status = getStatus(config, res)

	var specName string

	if groupSpec == nil {
		groupValue = "nogroup"
		groupTitle = "No Group"
		specName = "Ungrouped"
	} else {
		specName = groupSpec.Name
	}

	cache, exists := itemsCache[groupValue]
	if !exists {
		cache = &logGroup{
			title:       groupTitle,
			description: specName,
			groupValue:  groupValue,
		}
		itemsCache[groupValue] = cache
	}

	message, err := templating.ApplyTemplate(config.MessageTmpl, res)
	if err != nil {
		message = fmt.Sprintf("Field %s not found in data: %s", config.MessageField, err.Error())
	}

	tcache := cache.(*logGroup)
	tcache.timestamp = timestamp

	logLine := logLine{
		message:   message,
		data:      res,
		level:     getLevel(config, res),
		timestamp: timestamp,
	}

	if groupSpec == nil {
		logLine.tags = getTags(config.Tags, res)
	} else {
		logLine.tags = getTags(append(config.Tags, groupSpec.Tags...), res)
	}

	tcache.lines = append(tcache.lines, logLine)

	return
}

func getGroupAndTitle(config config.Config, line map[string]interface{}) (value string, title string, spec *config.GroupSpec) {
	for _, spec := range config.Groups {
		currentValue, valueErr := templating.ApplyTemplate(spec.ValueTmpl, line)
		currentTitle, titleErr := templating.ApplyTemplate(spec.TitleTmpl, line)

		if valueErr == nil && titleErr == nil && currentValue != "" && currentTitle != "" {
			return currentValue, currentTitle, spec
		}
	}

	return
}

func getStatus(config config.Config, line map[string]interface{}) (status string) {
	for _, spec := range config.Statuses {
		currentDisplay, displayErr := templating.ApplyTemplate(spec.DisplayTmpl, line)

		if displayErr == nil && currentDisplay != "" {
			return currentDisplay
		}
	}

	return
}

func getTimestamp(config config.Config, line map[string]interface{}) time.Time {
	timestampStr, err := templating.ApplyTemplate(config.TimestampTmpl, line)

	if err == nil && timestampStr != "" {
		t, err := time.Parse(time.RFC3339, timestampStr)
		if err == nil {
			return t
		}
	}

	return time.Now()
}

func getLevel(config config.Config, line map[string]interface{}) string {
	level, err := templating.ApplyTemplate(config.LevelTmpl, line)
	if err == nil {
		return level
	} else {
		return "unknown"
	}
}

func getTags(tags []*config.TagSpec, line map[string]interface{}) []logTag {
	result := []logTag{}
	tagsAlready := map[string]struct{}{}

	for _, tagSpec := range tags {
		key, err := templating.ApplyTemplate(tagSpec.KeyTmpl, line)
		if err != nil || key == "" {
			continue
		}

		if _, ok := tagsAlready[key]; ok {
			continue
		}

		if tagSpec.ValueTmpl != nil {
			value, err := templating.ApplyTemplate(tagSpec.ValueTmpl, line)
			if err == nil && value != "" {
				result = append(result, logTag{name: key, value: value})
				tagsAlready[key] = struct{}{}
			}
		} else {
			result = append(result, logTag{name: key})
			tagsAlready[key] = struct{}{}
		}
	}

	return result
}
