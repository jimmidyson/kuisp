// KUISP - A utility to serve static content & reverse proxy to RESTful services
//
// Copyright 2015 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/bradfitz/http2"
	"github.com/gorilla/handlers"
	flag "github.com/spf13/pflag"
)

type Options struct {
	Port               int
	StaticDir          string
	StaticPrefix       string
	Services           services
	Configs            configs
	CACerts            caCerts
	SkipCertValidation bool
	TlsCertFile        string
	TlsKeyFile         string
}

var options = &Options{}

func initFlags() {
	flag.IntVarP(&options.Port, "port", "p", 80, "The port to listen on")
	flag.StringVarP(&options.StaticDir, "www", "w", ".", "Directory to serve static files from")
	flag.StringVar(&options.StaticPrefix, "www-prefix", "/", "Prefix to serve static files on")
	flag.VarP(&options.Services, "service", "s", "The Kubernetes services to proxy to in the form \"<prefix>=<serviceUrl>\"")
	flag.VarP(&options.Configs, "config-file", "c", "The configuration files to create in the form \"<template>=<output>\"")
	flag.Var(&options.CACerts, "ca-cert", "CA certs used to verify proxied server certificates")
	flag.StringVar(&options.TlsCertFile, "tls-cert", "", "Certificate file to use to serve using TLS")
	flag.StringVar(&options.TlsKeyFile, "tls-key", "", "Certificate file to use to serve using TLS")
	flag.BoolVar(&options.SkipCertValidation, "skip-cert-validation", false, "Skip remote certificate validation - dangerous!")
	flag.Parse()
}

func main() {
	initFlags()

	if len(options.Configs) > 0 {
		for _, configDef := range options.Configs {
			fmt.Printf("Creating config file:  %v => %v\n", configDef.template, configDef.output)
			createConfig(configDef.template, configDef.output)
		}
		fmt.Println()
	}

	if len(options.Services) > 0 {
		tlsConfig := &tls.Config{
			RootCAs:            x509.NewCertPool(),
			InsecureSkipVerify: options.SkipCertValidation,
		}
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		if len(options.CACerts) > 0 {
			for _, caFile := range options.CACerts {
				// Load our trusted certificate path
				pemData, err := ioutil.ReadFile(caFile)
				if err != nil {
					log.Fatal("Couldn't read CA file, ", caFile, ": ", err)
				}
				if ok := tlsConfig.RootCAs.AppendCertsFromPEM(pemData); !ok {
					log.Fatal("Couldn't load PEM data from CA file, ", caFile)
				}
			}
		}
		for _, serviceDef := range options.Services {
			fmt.Printf("Creating service proxy: %v => %v\n", serviceDef.prefix, serviceDef.url.String())
			rp := httputil.NewSingleHostReverseProxy(serviceDef.url)
			rp.Transport = transport
			http.Handle(serviceDef.prefix, handlers.CombinedLoggingHandler(os.Stdout, http.StripPrefix(serviceDef.prefix, rp)))
		}
		fmt.Println()
	}

	fs := http.FileServer(http.Dir(options.StaticDir))
	http.Handle(options.StaticPrefix, handlers.CombinedLoggingHandler(os.Stdout, fs))

	fmt.Printf("Listening on :%d\n", options.Port)
	fmt.Println()

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", options.Port),
	}
	http2.ConfigureServer(srv, &http2.Server{})
	srv.Handler = http.DefaultServeMux

	if len(options.TlsCertFile) > 0 && len(options.TlsKeyFile) > 0 {
		log.Fatal(srv.ListenAndServeTLS(options.TlsCertFile, options.TlsKeyFile))
	} else {
		log.Fatal(srv.ListenAndServe())
	}
}
