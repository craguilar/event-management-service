package app

import (
	"bytes"
	"log"
	"text/template"
)

// I had to do this as a quick workaround to embed the content in the final binary :()
const TEMPLATE = `
<!DOCTYPE html>
<html>

<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Pending Tasks</title>
</head>

<body
  style="margin: 0; padding: 0; font-family: 'Helvetica Neue', Arial, sans-serif, sans-serif; background-color: #f4f4f4;">
  <table align="center" border="0" cellpadding="0" cellspacing="0" width="600"
    style="border-collapse: collapse; margin: 20px auto; background-color: #ffffff; border: 1px solid #dddddd; box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);">
    <tr>
      <td style="padding: 20px; text-align: left; background-color: #9494b8; color: white;">
        <h1 style="margin: 0;">Pending tasks for {{.EventName}}</h1>
      </td>
    </tr>
    <tr>
      <td style="padding: 20px;">
        <table border="0" cellpadding="10" cellspacing="0" width="100%">
          <thead>
            <tr>
              <th style="border-bottom: 2px solid #dddddd;">Name</th>
              <th style="border-bottom: 2px solid #dddddd;">Status</th>
              <th style="border-bottom: 2px solid #dddddd;">Created</th>
            </tr>
          </thead>
          <tbody>
            {{range .Tasks}}
            <tr>
              <td>{{.Name}}</td>
              <td style="color: #ffc107; padding: 5px 10px; border-radius: 3px;">{{.Status}}</td>
              <td>{{.TimeCreatedOn}}</td>
            </tr>
            {{end}}
          </tbody>

        </table>
      </td>
    </tr>
  </table>

</body>

</html>
`

type EventTasksTemplate struct {
	EventName string
	Tasks     []Task
}

func TemplatePendingTasksNotifications(eventName string, tasks []*Task) *bytes.Buffer {

	temp, err := template.New("pending-tasks").Parse(TEMPLATE)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	// Copy pointers into objects
	dataTasks := make([]Task, len(tasks))
	for i, ptr := range tasks {
		dataTasks[i] = *ptr
	}
	// prepare data
	data := EventTasksTemplate{
		EventName: eventName,
		Tasks:     dataTasks,
	}
	buf := new(bytes.Buffer)
	err = temp.Execute(buf, data)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return buf
}
