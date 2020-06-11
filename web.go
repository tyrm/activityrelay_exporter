package main

import (
	"net/http"
	"text/template"
)

const templateMetrics = `# HELP activitypub_requests Total number of requests received 
# TYPE activitypub_requests counter
{{- range $relay, $stats := . }}
{{- range $instance, $requests := $stats.Requests }}
activitypub_requests{relay="{{ $relay }}",instance="{{ $instance }}"} {{ $requests }}
{{- end }}
{{- end }}

# HELP activitypub_response_codes Total number of requests received by status code 
# TYPE activitypub_response_codes counter
{{- range $relay, $stats := . }}
{{- range $status, $requests := $stats.ResponseCodes }}
activitypub_response_codes{relay="{{ $relay }}",status="{{ $status }}"} {{ $requests }}
{{- end }}
{{- end }}

# HELP activitypub_response_codes_per_domain Total number of requests received by status code 
# TYPE activitypub_response_codes_per_domain counter
{{- range $relay, $stats := . }}
{{- range $instance, $requestsList := $stats.ResponseCodesPerDomain }}
{{- range $status, $requests := $requestsList }}
activitypub_response_codes_per_domain{relay="{{ $relay }}",instance="{{ $instance }}",status="{{ $status }}"} {{ $requests }}
{{- end }}
{{- end }}
{{- end }}

# HELP activitypub_delivery_codes Total number of requests received by status code 
# TYPE activitypub_delivery_codes counter
{{- range $relay, $stats := . }}
{{- range $status, $deliveries := $stats.DeliveryCodes }}
activitypub_delivery_codes{relay="{{ $relay }}",status="{{ $status }}"} {{ $deliveries }}
{{- end }}
{{- end }}

# HELP activitypub_delivery_codes_per_domain Total number of requests received by status code 
# TYPE activitypub_delivery_codes_per_domain counter
{{- range $relay, $stats := . }}
{{- range $instance, $requestsList := $stats.ResponseCodesPerDomain }}
{{- range $status, $deliveries := $requestsList }}
activitypub_delivery_codes_per_domain{relay="{{ $relay }}",instance="{{ $instance }}",status="{{ $status }}"} {{ $deliveries }}
{{- end }}
{{- end }}
{{- end }}

# HELP activitypub_exceptions Total number of requests received by status code 
# TYPE activitypub_exceptions counter
{{- range $relay, $stats := . }}
{{- range $error, $exceptions := $stats.Exceptions }}
activitypub_exceptions{relay="{{ $relay }}",error="{{ $error }}"} {{ $exceptions }}
{{- end }}
{{- end }}

# HELP activitypub_exceptions_per_domain Total number of requests received by status code 
# TYPE activitypub_exceptions_per_domain counter
{{- range $relay, $stats := . }}
{{- range $instance, $exceptionList := $stats.ExceptionsPerDomain }}
{{- range $error, $exceptions := $exceptionList }}
activitypub_exceptions_per_domain{relay="{{ $relay }}",instance="{{ $instance }}",error="{{ $error }}"} {{ $exceptions }}
{{- end }}
{{- end }}
{{- end }}

# HELP activitypub_delivery_exceptions Total number of requests received by status code 
# TYPE activitypub_delivery_exceptions counter
{{- range $relay, $stats := . }}
{{- range $error, $exceptions := $stats.DeliveryExceptions }}
activitypub_delivery_exceptions{relay="{{ $relay }}",error="{{ $error }}"} {{ $exceptions }}
{{- end }}
{{- end }}

# HELP activitypub_delivery_exceptions_per_domain Total number of requests received by status code 
# TYPE activitypub_delivery_exceptions_per_domain counter
{{- range $relay, $stats := . }}
{{- range $instance, $exceptionList := $stats.DeliveryExceptionsPerDomain }}
{{- range $error, $exceptions := $exceptionList }}
activitypub_delivery_exceptions_per_domain{relay="{{ $relay }}",instance="{{ $instance }}",error="{{ $error }}"} {{ $exceptions }}
{{- end }}
{{- end }}
{{- end }}
`

func handler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("metrics").Parse(templateMetrics))

	err := t.Execute(w, relays.Get())
	if err !=nil {
		logger.Errorf("Error processing template: %s", err.Error())
	}
}
