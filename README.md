![Check & test](https://github.com/summer-solutions/spring/workflows/Check%20&%20test/badge.svg)
[![codecov](https://codecov.io/gh/summer-solutions/springspring/branch/master/graph/badge.svg)](https://codecov.io/gh/summer-solutions/orm)
[![Go Report Card](https://goreportcard.com/badge/github.com/summer-solutions/spring)](https://goreportcard.com/report/github.com/summer-solutions/spring)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)



# spring

### Simple Framework designed to build scalable GraphQL services

### Main features:

 * Build on top of [Gqlgen](https://gqlgen.com/]) and [Gin Framework](https://github.com/gin-gonic/gin)
 * Easy to integrate with [Spring ORM](https://github.com/summer-solutions/orm)
 * Follows [Dependency injection](https://en.wikipedia.org/wiki/Dependency_injection) pattern
 
 
### Starting GraphQL Server

```go
package main
import "github.com/summer-solutions/spring"

func main() {
	
    graphQLExecutableSchema := ... // setup graphql.ExecutableSchema 
    // run http server
    spring.New("app_name").Run(8080, graphQLExecutableSchema)
}

``` 

#### Setting server port

By default, spring server is using port defined in environment variable "**PORT**". If this variable is not
set spring will use port number passed as fist argument.

#### Application name

When you setup server using **New** method yo must provide unique application name that can be
checked in code like this:

```go
    s := spring.New("app_name").Run(...)
    ioc.App().Name
```

#### Setting mode

By default, spring is running in "**spring.ModeLocal**" mode. Mode is a string that is available in: 
```go
    s := spring.New("app_name").Run(...)
    
    // now you can access current spring mode
    ioc.App().Mode
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
    s := spring.New(...)
    ioc.App().IsInLocalMode()
    ioc.App().IsInProdMode()
    ioc.App().IsInMode("my_mode")
```

#### Defining DI services

Spring builds global shared Dependency Injection container. You can register new services using this method:

```go
package main
import "github.com/summer-solutions/spring"

func main() {
    server := spring.New("my_app").RegisterDIService(
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
    "github.com/summer-solutions/spring/ioc"
)
    
func main() {
    myService := &ioc.ServiceDefinition{
        Name:   "my_service", // unique service key
        Global: true, // false if this service should be created as separate instance for each http request
        Build: func() (interface{}, error) {
            return &SomeService{}, nil // you can return any data you want
        },
    }
    
    // register it and run server
    server := spring.New("my_app").RegisterDIService(
      myService,
    )
    server.Run(...)
}

```

Now you can access this service in your code using:

```go
import (
    "github.com/summer-solutions/spring/ioc"
)

func SomeResolver(ctx context.Context) {

    ios.HasService("my_service") // return true
    
    // return error if Build function returned error
    myService, has, err := ios.GetContainer().GetServiceSafe("my_service") 
    // will panic if Build function returns error
    myService, has := ios.GetContainer().GetServiceOptional("my_service") 
    // will panic if service is not registered or Build function returned errors
    myService := ios.GetContainer().GetServiceRequired("my_service") 

    // if you registered service with field "Global" set to false (request service)

    myContextService, has, err := ioc.GetServiceForRequestSafe(ctx).Get("my_service_request")
    myContextService, has := ioc.GetServiceForRequestOptional(ctx).Get("my_service_request") 
    myContextService := ioc.GetServiceForRequestRequired(ctx).Get("my_service_request") 
}

```

It's a good practice to create one package in your app that provides all services with simple Getters:

```go
package ioc
import (
    "context"
    "github.com/summer-solutions/spring/ioc"
)

// optional global service    
func MyService() (MyServiceType, bool) {
    v, has := ioc.GetServiceOptional("my_service")
    if has {
        return v.(MyServiceType), true
    }
    return nil, false
}

// required global service    
func MyOtherService() MyOtherServiceType {
    return ioc.GetServiceRequired("my_other_service").(MyOtherServiceType)
}

// required request service    
func MyLastServiceForContext(ctx context.Context) MyContextServiceType {
    return ioc.GetServiceForRequestRequired(ctx, "my_contect_service").(MyContextServiceType)
}

```
