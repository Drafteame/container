# container

Flexible and opinionated IoC container for dependency injection with black magic

## Require

- Go >= 1.22

## Install

```bash
go get github.com/Drafteame/container@latest
```

## Usage

Using the global container you can access to all the container methods to manage dependency factories.

### Inline way

```go
package main

import (
	"github.com/Drafteame/container"
)

type param struct {}

type someType struct{
	p *param
}

func (*someType) SayHello() {
	println("hello")
}

func (*someType) SayGoodBye() {
	println("good bye")
}

type mainInterface interface{
	SayHello()
	SayGoodBye()
}

type subInterface interface {
	SayHello()
}

func someConstructor(p *param) *someType {
	return &someType{p: p}
}

func regularInstance() *someType {
	// Register a regular dependency with options
	err := container.Register("someName", someConstructor, container.WithArgs(container.MustGet[*param]("someParam")))
	if err != nil {
		panic(err)
	}
	
	return container.MustGet[*someType]("someName")
}

func singletonInstance() mainInterface {
	// Register a singleton dependency with options
	err := container.Singleton("someNameSingleton", someConstructor, container.WithArgs(container.MustGet[*param]("someParam")))
	if err != nil {
		panic(err)
	}

	return container.MustGet[mainInterface]("someNameSingleton")
}

func getSingletonAsSubInterface() subInterface {
	return container.MustGet[subInterface]("someNameSingleton")
}
```

### Functional way

```go
package main

import (
	"github.com/Drafteame/container"
)

type param struct {}

type someType struct{
	p *param
}

func (*someType) SayHello() {
	println("hello")
}

func (*someType) SayGoodBye() {
	println("good bye")
}

type mainInterface interface{
	SayHello()
	SayGoodBye()
}

type subInterface interface {
	SayHello()
}

func someConstructor(p *param) *someType {
	return &someType{p: p}
}

func regularInstance() *someType {
	err := container.Register("someName", func() *someType {
		p := container.MustGet[*param]("someParam")
		return someConstructor(p)
    })
	
	if err != nil {
		panic(err)
	}
	
	return container.MustGet[*someType]("someName")
}

func singletonInstance() mainInterface {
	err := container.Singleton("someNameSingleton", func() *someType {
		p := container.MustGet[*param]("someParam")
		return someConstructor(p)
	})
	
	if err != nil {
		panic(err)
	}

	return container.MustGet[mainInterface]("someNameSingleton")
}

func getSingletonAsSubInterface() subInterface {
	return container.MustGet[subInterface]("someNameSingleton")
}
```

## Dependencies

There are two types of dependencies: regular dependencies and singleton dependencies.

Regular dependencies are instances that each time that are required to be injected or retrieved, they will create a new
instance from the provided factory each time. This means that with this type of dependencies, you will have multiple
instances of the same type and this will not share any context. Basically, there is a fresh new instance each time we 
inject it.

```go
package main

import (
	"fmt"
	"github.com/Drafteame/container"
	"github.com/Drafteame/container/dependency"
)

type User struct {
	Name string
	Age  int
}

func newUser(name string, age int) *User {
	return &User{
		Age:  age,
		Name: name,
	}
}

func main() {
	depName := "someDep"
	
	// Register a regular dependency with options
	if err := container.Register(depName, newUser, container.WithArgs("John", 21)); err != nil {
		panic(err)
	}

	userInstance, err := container.Get[*User](depName)
	if err != nil {
		panic(err)
	}

	fmt.Println(userInstance)
}
```

Singleton dependencies are pretty much the same as a regular dependency with the particularity that the container will
keep the result obtained from the factory internally. They if a new instance of the same dependency is called to be
injected, instead of creating a new one from the factory will inject the previously created instance.

Keep in mind that this cannot work as a real singleton if the returned value of the factory is not a pointer or
interface.

```go
package main

import (
	"fmt"
	
	"github.com/Drafteame/container"
)

type User struct {
	Name string
	Age  int
}

func newUser(name string, age int) *User {
	return &User{
		Age:  age,
		Name: name,
	}
}

func main() {
	depName := "someDep"
	
	// Register a singleton dependency with options
	if err := container.Singleton(depName, newUser, container.WithArgs("John", 21)); err != nil {
		panic(err)
	}

	userInstance, err := container.Get[*User](depName)
	if err != nil {
		panic(err)
	}
	
	userInstance2, err := container.Get[*User](depName)
	if err != nil {
		panic(err)
    }
	
	if userInstance == userInstance2 {
		fmt.Println("same instance")	
    }
}
```

## Must Methods

The container package provides a set of "Must" methods that are variants of their regular counterparts. These methods panic instead of returning errors, which can be useful in scenarios where you want to fail fast if something goes wrong, such as during application initialization.

### MustGet

`MustGet` retrieves a dependency from the container and panics if it cannot be found or cast to the requested type:

```go
package main

import (
	"github.com/Drafteame/container"
)

func main() {
	// Register a dependency
	if err := container.Register("config", newConfig, container.WithArgs("production")); err != nil {
		panic(err)
	}
	
	// Get the dependency with MustGet - will panic if not found or cannot be cast
	config := container.MustGet[*Config]("config")
	
	// Use the dependency directly without error checking
	fmt.Println(config.Environment)
}
```

### MustRegister

`MustRegister` registers a dependency and panics if registration fails:

```go
package main

import (
	"github.com/Drafteame/container"
)

func main() {
	// Register a dependency with MustRegister - will panic if registration fails
	container.MustRegister("logger", newLogger, container.WithArgs("debug"))
	
	// Use the registered dependency
	logger := container.MustGet[Logger]("logger")
	logger.Info("Application started")
}
```

### MustSingleton

`MustSingleton` registers a singleton dependency and panics if registration fails:

```go
package main

import (
	"github.com/Drafteame/container"
)

func main() {
	// Register a singleton with MustSingleton - will panic if registration fails
	container.MustSingleton("database", newDatabase, container.WithArgs("connection-string"))
	
	// Use the registered singleton
	db := container.MustGet[Database]("database")
	db.Connect()
}
```

### MustOverride

`MustOverride` overrides an existing dependency and panics if the operation fails:

```go
package main

import (
	"github.com/Drafteame/container"
)

func main() {
	// Register the original dependency
	container.MustRegister("emailService", newRealEmailService)
	
	// Later, override it with a mock for testing
	container.MustOverride("emailService", newMockEmailService)
	
	// Use the overridden dependency
	emailService := container.MustGet[EmailService]("emailService")
	emailService.SendEmail("user@example.com", "Test Subject", "Test Body")
}
```

### When to Use Must Methods

Must methods are particularly useful in the following scenarios:

1. **Application Initialization**: When setting up your application, you often want to fail fast if a critical dependency cannot be registered or retrieved.

2. **Testing**: In tests, you may want to simplify error handling and focus on the test logic.

3. **Simple Applications**: In small applications or scripts where comprehensive error handling is not necessary.

However, be cautious when using Must methods in production code, especially in request handlers or other code that should gracefully handle errors, as panics can crash your application if not properly recovered.

## Options

The container package provides a flexible options pattern for configuring dependency registration and retrieval.

### WithArgs

The `WithArgs` option allows you to provide arguments to the factory function when registering a dependency:

```go
package main

import (
	"github.com/Drafteame/container"
)

func main() {
	// Register with direct arguments
	if err := container.Register("user", newUser, container.WithArgs("John", 21)); err != nil {
		panic(err)
	}
	
	// Register with a mix of direct arguments and dependencies
	if err := container.Register("service", newService, container.WithArgs(
		container.MustGet[*User]("user"),
		"api-key-123",
	)); err != nil {
		panic(err)
	}
}
```

### WithContainer

The `WithContainer` option allows you to specify a custom container instance to use instead of the global container:

```go
package main

import (
	"github.com/Drafteame/container"
)

func main() {
	// Create a custom container
	customContainer := container.New()
	
	// Register a dependency in the custom container
	if err := container.Register("user", newUser, 
		container.WithContainer(customContainer),
		container.WithArgs("John", 21),
	); err != nil {
		panic(err)
	}
	
	// Get the dependency from the custom container
	user, err := container.Get[*User]("user", container.WithContainer(customContainer))
	if err != nil {
		panic(err)
	}
}
```

## Override Dependencies

The container package provides a way to override existing dependencies, which is particularly useful for testing:

```go
package main

import (
	"github.com/Drafteame/container"
)

func main() {
	// Register the original dependency
	if err := container.Register("database", newRealDatabase, container.WithArgs("connection-string")); err != nil {
		panic(err)
	}
	
	// Later, override it with a mock for testing
	if err := container.Override("database", newMockDatabase, container.WithArgs()); err != nil {
		panic(err)
	}
	
	// The container will now return the mock implementation
	db, err := container.Get[Database]("database")
	if err != nil {
		panic(err)
	}
}
```

## Using Injectable Dependencies

You can use the `Inject` function to reference dependencies that are already registered in the container:

```go
package main

import (
	"github.com/Drafteame/container"
	"github.com/Drafteame/container/dependency"
)

func main() {
	// Register a database dependency
	if err := container.Register("database", newDatabase, container.WithArgs("connection-string")); err != nil {
		panic(err)
	}
	
	// Register a user repository that depends on the database
	if err := container.Register("userRepo", newUserRepository, container.WithArgs(container.Inject("database"))); err != nil {
		panic(err)
	}
	
	// Get the user repository
	repo, err := container.Get[*UserRepository]("userRepo")
	if err != nil {
		panic(err)
	}
}
```

## Custom Containers

You can create and use custom containers instead of the global container:

```go
package main

import (
	"github.com/Drafteame/container"
	"github.com/Drafteame/container/dependency"
)

func main() {
	// Create a custom container
	customContainer := container.New()
	
	// Register dependencies in the custom container
	if err := customContainer.Provide("user", dependency.New(newUser, "John", 21)); err != nil {
		panic(err)
	}
	
	// Get dependencies from the custom container
	user, err := customContainer.Get("user")
	if err != nil {
		panic(err)
	}
	
	// Type assertion is needed when using the container directly
	typedUser := user.(*User)
}
```