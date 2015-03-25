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
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"

	flag "github.com/spf13/pflag"
)

type Options struct {
	Port         int
	StaticDir    string
	StaticPrefix string
	Services     services
	Configs      configs
}

var options = &Options{}

func initFlags() {
	flag.IntVarP(&options.Port, "port", "p", 80, "The port to listen on")
	flag.StringVarP(&options.StaticDir, "www", "w", ".", "Directory to serve static files from")
	flag.StringVar(&options.StaticPrefix, "www-prefix", "/", "Prefix to serve static files on")
	flag.VarP(&options.Services, "service", "s", "The Kubernetes services to proxy to in the form \"<prefix>=<serviceUrl>\"")
	flag.VarP(&options.Configs, "config-file", "c", "The configuration files to create in the form \"<template>=<output>\"")
	flag.Parse()
}

func main() {
	initFlags()

	if len(options.Services) > 0 {
		for _, serviceDef := range options.Services {
			http.Handle(serviceDef.prefix, http.StripPrefix(serviceDef.prefix, httputil.NewSingleHostReverseProxy(serviceDef.url)))
		}
	}

	if len(options.Configs) > 0 {
		for _, configDef := range options.Configs {
			createConfig(configDef.template, configDef.output)
		}
	}

	fs := http.FileServer(http.Dir(options.StaticDir))
	http.Handle(options.StaticPrefix, fs)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(options.Port), nil))
}
