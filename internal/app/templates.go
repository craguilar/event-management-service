package app

import (
	"bytes"
	"log"
	"text/template"
)

// I had to do this as a quick workaround to embed the content in the final binary :()
const TEMPLATE = `<!DOCTYPE html>
<html>

<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Bootstrap Look and Feel</title>
  <!-- Add Bootstrap CSS link -->
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
</head>

<body>
  <div class="container">
    <h1 class="mt-4">Pending tasks for {{.EventName}}</h1>
    <hr class="mb-4">
    <div class="table-responsive">
      <table class="table table-bordered">
        <thead class="thead-dark">
          <tr>
            <th>Name</th>
            <th>Status</th>
            <th>Created</th>
          </tr>
        </thead>
        <tbody>
          {{range .Tasks}}
          <tr>
            <td>{{.Name}}</td>
            <td><span class="badge badge-warning">{{.Status}}</span></td>
            <td>{{.TimeCreatedOn}}</td>
          </tr>
          {{end}}
        </tbody>
      </table>
    </div>
  </div>

  <!-- Add Bootstrap JS and Popper.js scripts (optional) -->
  <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.16.0/umd/popper.min.js"></script>
  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.5.2/js/bootstrap.min.js"></script>
</body>

</html>`

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
