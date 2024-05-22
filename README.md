# LogicNG as a Service

A web service based on LogicNG.

It exposes 50 LogicNG functions as REST Web Services, among them
- SAT Solver
- MaxSAT Solver
- BDD Compilation
- DNNF Compilation
- Normal Forms
- Formula Simplifications
- ...

Many of the algorithms can be parametrized via query parameters.  Input/Output can be JSON or Protocol Buffer 
(configured via `accept` and `Content-Type` headers of the HTTP request).

## Usage

Run the service with a simple

```bash
go run main.go
```

This fires up the server on port 8080 with a default computation timeout of 5 seconds.  You can modify these 
parameters with

```bash
go run main.go -host "hostname" -port 9090 -timeout "20s"
```
to start the server on host "hostname", port 9090, and a timeout of 20 seconds.

A swagger documentation of all endpoints is available at `$host:$port/swagger` and should illustrate all available 
algorithms and configuration parameters.

## Chaining

The API is designed in a way, that the output of many of the endpoints can be used as input to many other endpoints.
So e.g. you can simply pass a set of formulas, substitute some variables, anonymize the formula, compute a normal form,
and generate a Mermaid.js visualization of the resulting formula without ever manipulating the input/output.

## Disclaimer

This is a funny little side project for playing around with Go Web Services (Standard Library only, no frameworks).  
It is by no means a production-ready piece of software with bells and whistles.  But if you just need some logical 
computations and don't want to integrate LogicNG in Go/Rust/Java in you project - perhaps it's for you :)

