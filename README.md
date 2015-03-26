# KUISP - Kubernetes UI & Service Proxy

KUISP (pronounced "kwisp") is a utility to simplify the
serving of static HTML/JavaScript UIs, backed by possibly
multiple RESTful API services. KUISP does the following:

* serve static content from the speicified directory
* use [Golang templates](http://golang.org/pkg/text/template/)
to create configuration files from specified templates &
environment variables
* sets up reverse proxies to the specified services
* http/2 server if TLS certificate/key is specified
* configurable CA certificates for remote TLS services

The benefit this brings is that browser clients only need network
access to the UI server, not the actual API servers.

## Running

```
Usage of kuisp:
  -l, --access-logging=false: Enable access logging
      --ca-cert=[]: CA certs used to verify proxied server certificates
      --compress=false: Enable gzip/deflate response compression
  -c, --config-file=[]: The configuration files to create in the form "<template>=<output>"
  -d, --default-page="": Default page to send if page not found
      --max-age=0: Set the Cache-Control header for static content with the max-age set to this value, e.g. 24h. Must confirm to http://golang.org/pkg/time/#ParseDuration
  -p, --port=80: The port to listen on
  -s, --service=[]: The Kubernetes services to proxy to in the form "<prefix>=<serviceUrl>"
      --skip-cert-validation=false: Skip remote certificate validation - dangerous!
      --tls-cert="": Certificate file to use to serve using TLS
      --tls-key="": Certificate file to use to serve using TLS
  -w, --www=".": Directory to serve static files from
      --www-prefix="/": Prefix to serve static files on
```

### Service proxying

For each service that you want to reverse proxy to, you can specify the `-s` or `--service`
flag. The input to the flag must be in the format `prefix=serviceUrl`, e.g.

    -s /db/=http://influxdb:8086/db/ -s /api/=http://apiserver:8080/api/v2/

As you can see you can specify the flag multiple times. You can also use environment
variable expansion such as:

    -s '/db/=http://${INFLUXDB_SERVICE_HOST}:${INFLUXDB_SERVICE_PORT}/db/'

Note the use of single quotes to ensure the environment variables don't get expanded
in your shell before being passed to KUISP.

### Configuration file templates

KUISP can process [Golang templates](http://golang.org/pkg/text/template/) into
configuration files. KUISP can process multiple templates into multiple configuration
files using the `-c` or `--config-file` flags. The configuration files & templates are specified in a similar fashion to the
service definitions:

    -c config.json.tmpl=config.json -c myconfig.ini.tmpl=/etc/myconfig.ini

Note that you can use relative or absolute paths for both templates & output files. Again,
you can use environment variable expansion for any values:

    -c '${TEMPLATE}=${OUTPUT}'

The only variables that are available in the template for
processing are environment variables & are accessed like this:

    {{ .Env.ENV_VAR_1 }}

As an example, the following template:

```
{
  "value": "{{ .Env.MY_VALUE }}"
}
```

When run with the environment variable `MY_VALUE` set to `3`, would output:

```
{
  "value": "3"
}
```

## Building

Just run `make`.
