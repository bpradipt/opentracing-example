package main

import (
	"flag"
	"fmt"
	"log"
        "os"
	"net/http"

	"github.com/opentracing/opentracing-go"

        zipkin "github.com/openzipkin/zipkin-go-opentracing"

)

var (
	port           = flag.Int("port", 8080, "Example app port.")
)

func main() {
	flag.Parse()

	var tracer opentracing.Tracer

        collector, err := zipkin.NewHTTPCollector(
                fmt.Sprintf("http://%s:9411/api/v1/spans", os.Args[1]))
        if err != nil {
             log.Fatal(err)
             return
        }
        tracer, err = zipkin.NewTracer(
            zipkin.NewRecorder(collector, false, "127.0.0.1:0", "example"),
        )
        if err != nil {
           log.Fatal(err)
         }

	opentracing.InitGlobalTracer(tracer)

	addr := fmt.Sprintf(":%d", *port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/home", homeHandler)
	mux.HandleFunc("/async", serviceHandler)
	mux.HandleFunc("/service", serviceHandler)
	mux.HandleFunc("/db", dbHandler)
	fmt.Printf("Go to http://localhost:%d/home to start a request!\n", *port)
	log.Fatal(http.ListenAndServe(addr, mux))
}
