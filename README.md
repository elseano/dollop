# Dollop

Bring a dollop of pizzaz to your JSON log output

Modern web development can generate lots of logs. Even if you have a nice console log formatter, you'll still run into hard to read log output (especially if you're logging database queries). Instead why not just use the JSON logging output, and pipe it through Dollop!

Dollop provides a configurable way to:

* Group log messages by requests, background workers, and other tasks.
* See at a glance where errors occurred.
* Dive deeper into the metadata of a particular log entry.


## Screenshots

### Configuration

Place a `.dollop.yml` file in your folder to configure how Dollop parses you json logs. Go templating is supported.

``` yaml
---
# All log messages are in the msg key
messageField: msg

# Timestamp comes from the time key
timestampField: time

# Log level is the level key
levelField: level

tags:
  # Adds an "error" tag onto the log entry if there's an error key in the metadata.
  - key: "{{ if .error }}error{{ end }}"
  # Adds a "stacktrace" tag onto the log entry if there's a stacktrace key in the metadata.
  - key: "{{ if .stacktrace }}stacktrace{{ end }}"
  # Adds duration onto the log entry if either duration or duration_ms is present.
  - key: duration
    value: "{{ if .elapsed_ms }}{{ FormatSeconds (div .elapsed_ms 1000) }}{{ else if .duration }}{{ FormatSeconds .duration }}{{end}}"

groups:
  # Group request logs by their request_id, and label them as HTTP in the sidebar.
  - valueField: request_id
    titleField: "{{ .method }} {{ .path }}"
    name: HTTP
  # Group background job logs by their messageId, and label them as Worker in the sidebar.
  - valueField: messageId
    titleField: "{{ .messageId }} - {{ .messageType }}"
    name: Worker
  # Otherwise, categorise by category if present.
  - valueField: category
    titleField: category
    name: Category
```

