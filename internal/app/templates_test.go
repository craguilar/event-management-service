package app

import (
	"testing"
	"time"
)

func TestPendingTasksTemplate(t *testing.T) {
	tasks := []*Task{
		{
			Name:          "Pick up the kids at 3 pm",
			Status:        "PENDING",
			TimeCreatedOn: time.Now(),
		},
		{
			Name:          "Do grocery shoping",
			Status:        "PENDING",
			TimeCreatedOn: time.Now(),
		},
	}
	buf := TemplatePendingTasksNotifications("Demo", tasks)
	if buf == nil {
		t.Fail()
	}
}
