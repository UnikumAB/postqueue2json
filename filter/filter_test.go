package filter

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestConvertPostqueueToQueueItem(t *testing.T) {
	type args struct {
		filename string
	}
	type wants struct {
		count int
	}
	tests := []struct {
		name    string
		args    args
		want    wants
		wantErr bool
	}{
		{
			name: "basic test",
			args: args{
				filename: "testdata/postqueue.txt",
			},
			want:    wants{count: 654},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := make(chan string)
			go func() {
				defer close(lines)
				file, err := os.Open(tt.args.filename)
				if err != nil {
					t.Errorf("failed to open file %v: %v", tt.args.filename, err)
				}
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					lines <- scanner.Text()
				}
			}()

			queueItems, err := ConvertPostqueueToQueueItem(lines)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertPostqueueToQueueItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			counter := 0
			for queueItem := range queueItems {
				counter++
				if queueItem.QueueId == "" {
					t.Errorf("No QueueId found for item %v", queueItem)
				}
				if strings.HasSuffix(queueItem.QueueId, "*") {
					t.Errorf("Not split QueueID correct: %v", queueItem.QueueId)
				}
				if queueItem.Sender == "" {
					t.Errorf("No sender found for item %v", queueItem)
				}
				if queueItem.Recipient == "" {
					t.Errorf("No recipient found for item %v", queueItem)
				}
			}
			if counter != tt.want.count {
				t.Errorf("Expected %v items but queueItems %v", tt.want.count, counter)
			}
		})
	}
}
