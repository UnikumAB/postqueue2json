package filter

import (
	"regexp"
	"strings"
)

type QueueItem struct {
	QueueId   string
	Sender    string
	Queue     string
	Recipient string
	Message   string
}

var firstLineExpr = regexp.MustCompile(`^(\w+)([*!])?\s+(\d+)\s+(\w+ \w+ \d+ \d+:\d+:\d+)\s+(\S+@\S+)$`)
var msgLineExpr = regexp.MustCompile(`^\s*\((.*)\)$`)
var recipientExpr = regexp.MustCompile(`^\s*(\S+@\S+)$`)

func ConvertPostqueueToQueueItem(lines <-chan string) (<-chan QueueItem, error) {
	items := make(chan QueueItem)
	go func() {
		defer close(items)
		var item QueueItem
		for line := range lines {
			if strings.HasPrefix(line, "-") {
				continue
			}
			switch {
			case strings.TrimSpace(line) == "":
				if item.QueueId != "" {
					items <- item
				}
				item = QueueItem{}
			case firstLineExpr.MatchString(line):
				firstLine := firstLineExpr.FindStringSubmatch(line)
				item.QueueId = firstLine[1]
				switch firstLine[2] {
				case "!":
					item.Queue = "hold"
				case "*":
					item.Queue = "active"
				default:
					item.Queue = "other"
				}
				item.Sender = firstLine[5]
			case msgLineExpr.MatchString(line):
				item.Message = msgLineExpr.FindStringSubmatch(line)[1]
			case recipientExpr.MatchString(line):
				item.Recipient = recipientExpr.FindStringSubmatch(line)[1]
			}
		}
	}()
	return items, nil
}
