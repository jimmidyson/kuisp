package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

type Options struct {
	Port               int
	StaticDir          string
	StaticPrefix       string
	ServiceName        string
	ServiceScheme      string
	ServiceProxyPrefix string
	ServiceProxyDir    string
	ConfigFileTemplate string
	OutputConfigFile   string
}

var options = &Options{}

func initFlags() {
	flag.IntVarP(&options.Port, "port", "p", 80, "The port to listen on")
	flag.StringVarP(&options.StaticDir, "www", "w", ".", "Directory to serve static files from")
	flag.StringVar(&options.StaticPrefix, "www-prefix", "/", "Prefix to serve static files on")
	flag.StringVarP(&options.ServiceName, "service", "s", "", "The Kubernetes service to proxy")
	flag.StringVarP(&options.ServiceScheme, "service-scheme", "S", "http", "The protocol for the Kubernetes service")
	flag.StringVar(&options.ServiceProxyPrefix, "service-proxy-prefix", "service", "The prefix to serve the service proxy under")
	flag.StringVar(&options.ServiceProxyDir, "service-proxy-dir", "", "The optional directory in the service to proxy to")
	flag.StringVarP(&options.ConfigFileTemplate, "config-template", "t", "", "The template file to create the specified config file from")
	flag.StringVarP(&options.OutputConfigFile, "config-file", "c", "", "The output config file")
	flag.Parse()
}

func main() {
	initFlags()

	if len(options.ServiceName) > 0 {
		serviceHost := os.ExpandEnv(strings.Replace(strings.ToUpper(fmt.Sprintf("${%v_SERVICE_HOST}", options.ServiceName)), "-", "_", -1))
		if len(serviceHost) == 0 {
			log.Fatal("Service host for ", options.ServiceName, " does not exist")
		}
		servicePort := os.ExpandEnv(strings.Replace(strings.ToUpper(fmt.Sprintf("${%v_SERVICE_PORT}", options.ServiceName)), "-", "_", -1))
		if len(servicePort) == 0 {
			log.Fatal("Service port for ", options.ServiceName, " does not exist")
		}
		serviceUrl := &url.URL{
			Scheme: options.ServiceScheme,
			Host:   fmt.Sprintf("%v:%v", serviceHost, servicePort),
			Path:   options.ServiceProxyDir,
		}
		http.Handle(options.ServiceProxyPrefix, http.StripPrefix(options.ServiceProxyPrefix, httputil.NewSingleHostReverseProxy(serviceUrl)))
	}

	if len(options.ConfigFileTemplate) > 0 && len(options.OutputConfigFile) > 0 {
		createConfig(options.ConfigFileTemplate, options.OutputConfigFile)
	}

	fs := http.FileServer(http.Dir(options.StaticDir))
	http.Handle(options.StaticPrefix, fs)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(options.Port), nil))
}
