/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
)

func getKubeEnv() (map[string]string, error) {
	environS := os.Environ()
	environ := make(map[string]string)
	for _, val := range environS {
		split := strings.Split(val, "=")
		if len(split) != 2 {
			return environ, fmt.Errorf("Some weird env vars")
		}
		environ[split[0]] = split[1]
	}
	for key := range environ {
		if !(strings.HasSuffix(key, "_SERVICE_HOST") ||
			strings.HasSuffix(key, "_SERVICE_PORT")) {
			delete(environ, key)
		}
	}
	return environ, nil
}

func printInfo(resp http.ResponseWriter, req *http.Request) {
	kubeVars, err := getKubeEnv()
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	backendHost := os.Getenv("BACKEND_SRV_SERVICE_HOST")
	backendPort := os.Getenv("BACKEND_SRV_SERVICE_PORT")
	backendRsp, backendErr := http.Get(fmt.Sprintf(
		"http://%v:%v/",
		backendHost,
		backendPort))
	if backendErr == nil {
		defer backendRsp.Body.Close()
	}

	name := os.Getenv("POD_NAME")
	namespace := os.Getenv("POD_NAMESPACE")
        podip := os.Getenv("POD_IP")
	fmt.Fprintf(resp, "Application version: v7 \n")
	fmt.Fprintf(resp, "Pod Name: %v \n", name)
	fmt.Fprintf(resp, "Pod Namespace: %v \n", namespace)
	fmt.Fprintf(resp, "Pod IP: %v \n", podip)

	/* Get and display secrets */	
	uname := os.Getenv("SECRET_USERNAME")
	pwd := os.Getenv("SECRET_PASSWORD")
	fmt.Fprintf(resp, "SECRET_USERNAME: %v \n", uname)
	fmt.Fprintf(resp, "SECRET_PASSWORD: %v \n", pwd)

	envvar := os.Getenv("USER_VAR")
	fmt.Fprintf(resp, "USER_VAR: %v \n", envvar)

	fmt.Fprintf(resp, "\nKubernetes environment variables\n")
	var keys []string
	for key := range kubeVars {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(resp, "%v = %v \n", key, kubeVars[key])
	}

	fmt.Fprintf(resp, "\nFound backend ip: %v port: %v\n", backendHost, backendPort)
	if backendErr == nil {
		fmt.Fprintf(resp, "Response from backend\n")
		io.Copy(resp, backendRsp.Body)
	} else {
		fmt.Fprintf(resp, "Error from backend: %v", backendErr.Error())
	}
}

func main() {
	started := time.Now()
	http.HandleFunc("/", printInfo)

        http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
                duration := time.Now().Sub(started)
                w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("ok. Alive for duration(s): %v", duration.Seconds())))
        })
        http.HandleFunc("/poststart", func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(200)
                w.Write([]byte("Post-start invoked"))
        })
        http.HandleFunc("/prestop", func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(200)
                w.Write([]byte("Cleaning up before stopping"))
        })
	http.HandleFunc("/redisincr", func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(200)
                w.Write([]byte("redisincr invoked"))


		// Now we connect to the Redis server
		conn, err := redisurl.ConnectToURL("redis://aml-rg-redis-001.8yizaq.0001.usw2.cache.amazonaws.com")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
                defer conn.Close()

                w.Write([]byte("redisincr - connected to Redis"))
		n, err := conn.Do("INCR", "aml-counter")
		w.Write([]byte(fmt.Sprintf("redisincr after INCR: %v", n)))
	})

        log.Fatal(http.ListenAndServe(":8080", nil))
}
