<a href="https://www.logicng.org"><img src="https://github.com/booleworks/logicng-go/blob/main/doc/logos/logicng_logo_gopher.png?raw=true" alt="logo" width="400"></a>

<a href="https://pkg.go.dev/github.com/booleworks/logicng-go"><img src="https://pkg.go.dev/badge/github.com/booleworks/logicng-go.svg" alt="Go Reference"></a>
[![license](https://img.shields.io/badge/license-MIT-purple?style=flat-square)]()


# LogicNG as a Service

A web service based on LogicNG.

It exposes over 50 LogicNG functions as REST Web Services, among them
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

## Compile it yourself

Install and initialize `swag` (required for generating the Swagger UI)
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

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

## Use the binary
You can just download a binary under [releases](https://github.com/booleworks/logicng-service/releases) and you should be ready to go.

## Docker
... or just use docker
```bash
docker run -p 8080:8080 ghcr.io/booleworks/logicng-service:0.0.2
```

## Swagger

A swagger documentation of all endpoints is available at `$host:$port/swagger` and should illustrate all available 
algorithms and configuration parameters.

<img src="https://github.com/booleworks/logicng-service/blob/main/assets/swagger.png?raw=true" alt="swagger" width="600">

<img src="https://github.com/booleworks/logicng-service/blob/main/assets/swagger_detail.png?raw=true" alt="swagger details" width="600">

## Functions

| Method   | Endpoint                         | Input               | Output            | Query Params                         |
| -------  | -------------------------------- | ------------------- | ----------------- | ------------------------------------ |
| `POST`   | `assignment/evaluation`          | `AssignmentInput`   | `BoolResult`      | -                                    |
| `POST`   | `assignment/restriction`         | `AssignmentInput`   | `FormulaResult`   | -                                    |
| `POST`   | `bdd/compilation`                | `FormulaInput`      | `GraphResult`     | Variable Ordering                    |
| `POST`   | `bdd/graphical`                  | `FormulaInput`      | `String`          | Variable Ordering, Graph Format      |
| `POST`   | `dnnf/compilation`               | `FormulaInput`      | `FormulaResult`   | -                                    |
| `POST`   | `encoding/cc`                    | `FormulaInput`      | `FormulaResult`   | Encoding Algorithm                   |
| `POST`   | `encoding/pbc`                   | `FormulaInput`      | `FormulaResult`   | Encoding Algorithm                   |
| `POST`   | `explanation/mus`                | `FormulaInput`      | `FormulaResult`   | MUS Algorithm                        |
| `POST`   | `explanation/smus`               | `FormulaInput`      | `FormulaResult`   | -                                    |
| `POST`   | `formula/atoms`                  | `FormulaInput`      | `IntResult`       | -                                    |
| `POST`   | `formula/depth`                  | `FormulaInput`      | `IntResult`       | -                                    |
| `POST`   | `formula/graphical`              | `FormulaInput`      | `String`          | Graph Type, Graph Format             |
| `POST`   | `formula/lit-profile`            | `FormulaInput`      | `ProfileResult`   | -                                    |
| `POST`   | `formula/literals`               | `FormulaInput`      | `StringSetResult` | -                                    |
| `POST`   | `formula/nodes`                  | `FormulaInput`      | `IntResult`       | -                                    |
| `POST`   | `formula/sub-formulas`           | `FormulaInput`      | `FormulaResult`   | -                                    |
| `POST`   | `formula/var-profile`            | `FormulaInput`      | `ProfileResult`   | -                                    |
| `POST`   | `formula/variables`              | `FormulaInput`      | `StringSetResult` | -                                    |
| `POST`   | `graph/components`               | `FormulaInput`      | `ComponentResult` | -                                    |
| `POST`   | `graph/constraint`               | `FormulaInput`      | `GraphResult`     | -                                    |
| `POST`   | `graph/constraint/graphical`     | `FormulaInput`      | `String`          | Graph Format                         |
| `POST`   | `model/counting`                 | `FormulaInput`      | `StringResult`    | Counting Algorithm                   |
| `POST`   | `model/counting/projection`      | `FormulaVarsInput`  | `StringResult`    | Counting Algorithm                   |
| `POST`   | `model/enumeration`              | `FormulaInput`      | `FormulaResult`   | Enumeration Algorithm                |
| `POST`   | `model/enumeration/projection`   | `FormulaVarsInput`  | `FormulaResult`   | Enumeration Algorithm                |
| `POST`   | `normalform/predicate/nnf`       | `FormulaInput`      | `BoolResult`      | -                                    |
| `POST`   | `normalform/predicate/cnf`       | `FormulaInput`      | `BoolResult`      | -                                    |
| `POST`   | `normalform/predicate/dnf`       | `FormulaInput`      | `BoolResult`      | -                                    |
| `POST`   | `normalform/predicate/aig`       | `FormulaInput`      | `BoolResult`      | -                                    |
| `POST`   | `normalform/predicate/minterm`   | `FormulaInput`      | `BoolResult`      | -                                    |
| `POST`   | `normalform/predicate/maxterm`   | `FormulaInput`      | `BoolResult`      | -                                    |
| `POST`   | `normalform/transformation/aig`  | `FormulaInput`      | `FormulaResult`   | -                                    |
| `POST`   | `normalform/transformation/cnf`  | `FormulaInput`      | `FormulaResult`   | CNF Algorithm                        |
| `POST`   | `normalform/transformation/dnf`  | `FormulaInput`      | `FormulaResult`   | DNF Algorithm                        |
| `POST`   | `normalform/transformation/nnf`  | `FormulaInput`      | `FormulaResult`   | -                                    |
| `POST`   | `prime/minimal-cover`            | `FormulaInput`      | `FormulaResult`   | Min or Max Models                    |
| `POST`   | `prime/minimal-implicant`        | `FormulaInput`      | `FormulaResult`   | -                                    |
| `GET`    | `randomizer`                     | -                   | `FormulaResult`   | Seed, Depth, Vars, Formulas          |
| `POST`   | `simplification/advanced`        | `FormulaInput`      | `FormulaResult`   | Backbone, Factor Out, Negation Flags |
| `POST`   | `simplification/backbone`        | `FormulaInput`      | `FormulaResult`   | -                                    |
| `POST`   | `simplification/distribution`    | `FormulaInput`      | `FormulaResult`   | -                                    |
| `POST`   | `simplification/factorout`       | `FormulaInput`      | `FormulaResult`   | -                                    |
| `POST`   | `simplification/negation`        | `FormulaInput`      | `FormulaResult`   | -                                    |
| `POST`   | `simplification/qmc`             | `FormulaInput`      | `FormulaResult`   | -                                    |
| `POST`   | `simplification/subsumption`     | `FormulaInput`      | `FormulaResult`   | -                                    |
| `POST`   | `simplification/unitpropagation` | `FormulaInput`      | `FormulaResult`   | -                                    |
| `POST`   | `solver/backbone`                | `FormulaInput`      | `BackboneResult`  | -                                    |
| `POST`   | `solver/maxsat`                  | `MaxSatInput`       | `MaxSatResult`    | MaxSAT Algorithm                     |
| `POST`   | `solver/predicate/contradiction` | `FormulaInput`      | `BoolResult`      | -                                    |
| `POST`   | `solver/predicate/equivalence`   | `FormulaInput`      | `BoolResult`      | -                                    |
| `POST`   | `solver/predicate/implication`   | `FormulaInput`      | `BoolResult`      | -                                    |
| `POST`   | `solver/predicate/tautology`     | `FormulaInput`      | `BoolResult`      | -                                    |
| `POST`   | `solver/sat`                     | `FormulaInput`      | `SatResult`       | UNSAT Core Flag                      |
| `POST`   | `substitution/anonymization`     | `FormulaInput`      | `FormulaResult`   | Variable Prefix                      |
| `POST`   | `substitution/variables`         | `SubstitutionInput` | `FormulaResult`   | -                                    | 

## Chaining

The API is designed in a way, that the output of many of the endpoints can be used as input to many other endpoints.
So e.g. you can simply pass a set of formulas, substitute some variables, anonymize the formula, compute a normal form,
and generate a Mermaid.js visualization of the resulting formula without ever manipulating the input/output.  In the table 
above this means that you always can use a `FormulaResult` from one endpoint directly as `FormulaInput` for another endpoint.

## Disclaimer

This is a funny little side project for playing around with Go Web Services (Standard Library only, no frameworks).  
It is by no means a production-ready piece of software with bells and whistles.  But if you just need some logical 
computations and don't want to integrate LogicNG in Go/Rust/Java in you project - perhaps it's for you :)

