# spring

### Simple Framework designed to build scalable GraphQL services

### Main features:

 * Build on top of [GraphQL Server](https://gqlgen.com/]) and [Gin Framework](https://github.com/gin-gonic/gin])
 * Easy to integrate with [Spring ORM](https://github.com/summer-solutions/orm])
 * Follows [Dependency injection](https://en.wikipedia.org/wiki/Dependency_injection) pattern
 
 
### Starting GraphQL Server

```go
package main
import "github.com/summer-solutions/spring"

func main() {
	
    graphQLExecutableSchema := ... // setup graphql.ExecutableSchema 
    // run http server
    spring.NewServer().Run(8080, graphQLExecutableSchema)
}

``` 

#### Setting server port
By default, spring server is using port defined in environment variable "**PORT**". If this variable is not
set spring will use port number passed as fist argument.

#### Setting mode

By default, spring is running in "**spring.ModeLocal**" mode. Mode is a string that is available in: 
```go
    s := spring.NewServer(...)
    s.GetMode() // returns current spring mode
```

You can define spring mode using special environment variable "**SPRING_MODE**".

Spring provides by default two modes:

 * **spring.ModeLocal**
   * should be used on local development machine (developer laptop)
   * errors and stack trace is printed directly to system console
   * log level is set to Debug level
   * log is formatted using human friendly console text formatter
   * Gin Framework is running in GinDebug mode
  * **spring.ModeProd**
    * errors and stack trace is printed only using Log
    * log level is set to Warn level
    * log is formatted using json formatter   
    
Mode is just a string. You can define any name you want. Remember that every mode that you create
follows **spring.ModeProd** rules explained above.
    
    
In code you can easly check current mode using one of these methods:    

```go
    s := spring.NewServer(...)
    s.IsInLocalMode()
    s.IsInProdMode()
    s.IsInMode("my_mode")
```

#### Defining DI services

Spring builds global shared Dependency Injection container. You can register new services using this method:

```go
package main
import "github.com/summer-solutions/spring"

func main() {
    server := spring.NewServer().RegisterDIService(
      // put service definitions here
    )
    server.Run(...)
}

``` 

Example of DI service definition:

```go
package main
import (
    "github.com/summer-solutions/spring"
    "github.com/summer-solutions/spring/di"
)
    
func main() {
    myService := &di.ServiceDefinition{
        Name:   "my_service", // unique service key
        Global: true, // false if this service should be created as separate instance for each http request
        Build: func() (interface{}, error) {
            return &SomeService{}, nil // you can return any data you want
        },
    }
    
    // register it and run server
    server := spring.NewServer().RegisterDIService(
      myService,
    )
    server.Run(...)
}

```

Now you can access this service in your code using:

```go
import (
    "github.com/summer-solutions/spring/di"
)

func SomeResolver(ctx context.Context) {
    myService := di.GetContainer().Get("my_service") 
    // or if service is defined with **Global** to false
    myContextService := di.GetContainerForRequest(ctx).Get("my_service_request") 
}

```

It's a good practice to create one package in your app that provides all services with simple Getters:

```go
package di
import (
    "github.com/summer-solutions/spring/di"
)

// optional global service    
func MyService() (MyServiceType, bool) {
    v, err := di.GetContainer().SafeGet("my_service")
    if err == nil {
        return v.(MyServiceType, true
    }
    return nil, false
}

// required global service    
func MyOtherService() MyOtherServiceType {
    return di.GetContainer().Get("my_other_service")
}

// required request service    
func MyLastServiceForContext() MyContextServiceType {
    return di.GetContainerForRequest(ctx context.Context).Get("my_other_service")
}

```
