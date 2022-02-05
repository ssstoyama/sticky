# sticky

A sticky is simple and minimal dependency injection container for Go.

## Installation

Install sticky by running:

```
$ go get github.com/ssstoyama/sticky
```

## Usage

### sticky.New

Create a new container. Can also specify container options.

```go
// no option
c := sticky.New()

// cache option
c := sticky.New(sticky.Cache(false)) // cache is default true
```

### sticky.Register

Register dependencies. Constructor Function(Factory Function) and Parameter can be registered.

```go
/* service.go */

func NewService(db Repository) *Service {
  return &Service{db}
}


/* container.go */

var c sticky.Container
...
// register a constructor
err := sticky.Register(c, sticky.Constructor(NewService))
// specify options
err := sticky.Register(c, sticky.Constructor(NewService, sticky.Implements[IService](), sticky.Tag("service_tag")))
// register a parameter
err := sticky.Register(c, sticky.Param("http://localhost", "endpoint_tag"))
```

### sticky.Resolve

Resolve will resolve the registered dependencies.

```go
var c sticky.Container
...
// resolve the service
service, err := sticky.Resolve[*Service]()
// with tag
service, err := sticky.Resolve[*Service](sticky.Tag("service_tag"))
```

### sticky.Extract

Extract allows to pass registered dependencies to functions.

```go
var c sticky.Container
...
err := sticky.Extract(c, func (service *Service) {
  /* some code */
})
```

### sticky.Decorate

Decorate allows to modify the registered dependencies.

```go
var c sticky.Container
...
err := sticky.Decorate(c, func(s *Service) *Service {
  /* modify service */
  s.param = "decorated"
  return s, nil
})
```

### sticky.Validate

Validate allows to verify that the dependencies are registered correctly.

```go
var c sticky.Container
...
err := sticky.Validate(c)
```
