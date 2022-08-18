package tui

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elseano/sl/internal/config"
)

type scanMsg struct {
	lines []list.Item
}

func (m Model) scanLogs() tea.Cmd {
	return func() tea.Msg {
		processLog(m.config)
		disp := []list.Item{}

		for _, v := range itemsCache {
			disp = append(disp, v)
		}

		sort.Slice(disp, func(i, j int) bool {
			return disp[i].(*item).timestamp.Sub(disp[j].(*item).timestamp) > 0
		})

		return scanMsg{lines: disp}
	}
}

var scanner = bufio.NewReader(os.Stdin)
var itemsCache = map[string]list.Item{
	"errors": &item{
		title:     "Errors",
		timestamp: time.Now(),
	},
}

func processLog(config config.Config) {

	line, err := scanner.ReadString('\n')

	if err == nil {
		var res map[string]interface{}
		if err := json.Unmarshal([]byte(line), &res); err == nil {
			groupValue, groupTitle, specName := getGroupAndTitle(config, res)
			timestamp := getTimestamp(config, res)

			if groupValue == "" {
				groupValue = "nogroup"
				groupTitle = "No Group"
			}

			cache, exists := itemsCache[groupValue]
			if !exists {
				cache = &item{
					title:       groupTitle,
					description: specName,
					timestamp:   timestamp,
				}
				itemsCache[groupValue] = cache
			}

			message, ok := res[config.MessageField].(string)
			if !ok {
				message = fmt.Sprintf("Field %s not found in data", config.MessageField)
			}

			tcache := cache.(*item)
			tcache.lines = append(tcache.lines, logLine{
				message: message,
				data:    res,
			})
		} else {
			tcache := itemsCache["errors"].(*item)
			tcache.lines = append(tcache.lines, logLine{
				message: fmt.Sprintf("Error loading '%s': %s", line, err.Error()),
				data:    res,
			})

		}
	}
}

func getGroupAndTitle(config config.Config, line map[string]interface{}) (value string, title string, name string) {
	for _, spec := range config.Groups {
		currentValue, valueOk := line[spec.ValueField].(string)
		currentTitle, titleOk := line[spec.TitleField].(string)

		if valueOk && titleOk {
			return currentValue, currentTitle, spec.Name
		}
	}

	return
}

func getTimestamp(config config.Config, line map[string]interface{}) time.Time {
	timestampStr, ok := line[config.TimestampField].(string)
	if ok {
		t, err := time.Parse(time.RFC3339, timestampStr)
		if err == nil {
			return t
		}
	}

	return time.Now()
}
