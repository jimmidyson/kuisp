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
	"fmt"
	"net/url"
	"os"
	"strings"
)

type service struct {
	prefix string
	url    *url.URL
}
type services []service

func (s *services) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *services) Set(value string) error {
	splitServiceDef := strings.Split(value, "=")
	if len(splitServiceDef) != 2 {
		return fmt.Errorf("Invalid service definition: ", value)
	}
	if serviceUrl, err := url.Parse(os.ExpandEnv(splitServiceDef[1])); err != nil {
		return fmt.Errorf("Invalid service URL: ", splitServiceDef[1])
	} else {
		serviceDef := service{
			prefix: os.ExpandEnv(splitServiceDef[0]),
			url:    serviceUrl,
		}
		options.Services = append(options.Services, serviceDef)
	}
	return nil
}

func (s *services) Type() string {
	return "services"
}

type config struct {
	template string
	output   string
}
type configs []config

func (s *configs) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *configs) Set(value string) error {
	splitConfigDef := strings.Split(value, "=")
	if len(splitConfigDef) != 2 {
		return fmt.Errorf("Invalid config definition: ", value)
	}
	configDef := config{
		template: os.ExpandEnv(splitConfigDef[0]),
		output:   os.ExpandEnv(splitConfigDef[1]),
	}
	options.Configs = append(options.Configs, configDef)
	return nil
}

func (s *configs) Type() string {
	return "configs"
}
