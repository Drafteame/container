# container

![gopher-container](https://user-images.githubusercontent.com/9085902/216511169-b28f488e-1f9c-4a8e-8ec9-4cb90db1c25a.png)

Flexible runtime dependency container inspired on go.uber.org/dig and based on reflection. It applies the dependency tree
concept to make flexible injections.

## Require

- Go >= 1.20

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
	err := container.Register("someName", someConstructor, container.Inject("someParam"))
	if err != nil {
		panic(err)
	}
	
	return container.MustGet[*someType]("someName")
}

func singletonInstance() mainInterface {
	err := container.Singleton("someNameSingleton", someConstructor, container.Inject("someParam"))
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

There two types of dependencies, regular dependencies and singleton dependencies.

Regular dependencies are instances that each time that are required to be injected or retrieved, they will create a new
instance from the provided factory each time. This means that with this type of dependencies, you will have multiple
instances of the same type and this will not share any context. Basically is a fresh new instance each time we inject it.

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
	dep := dependency.New(newUser, "John", 21)

	if err := container.Register(depName, dep); err != nil {
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
keep the result obtained from the factory internally and if a new instance of the same dependency is called to be
injected, instead of create a new one from the factory will inject the previous created instance.

Keep in mind that this can not work as a real singleton if the returned value of the factory is not a pointer or
interface.

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
	dep := dependency.NewSingleton(newUser, "John", 21)

	if err := container.Register(depName, dep); err != nil {
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

Arguments of the regular and singleton dependencies can be plain values, other `dependency.Dependency` objects or
`dependency.Injectable` instances. This last type of argument are objects that make reference to a dependency that was
registered in the container previously. This is specially helpful if you do not want to redefine a dependency many
times, and just reuse same specification of the dependency.

Example of plain values as dependency arguments:

```go
package main

import "github.com/Drafteame/container/dependency"

func main() {
	name := "foo"
	age := 21

	// Regular dependency
	depName := "test"
	dep := dependency.New(newUser, name, age)

	// Singleton dependency
	depName2 := "test2"
	dep2 := dependency.NewSingleton(newUser, name, age)
}
```

Example of dependency instances as arguments:

```go
package main

import (
	"os"
	
	"github.com/Drafteame/container/dependency"
)

func main() {
	driver := dependency.New(newDB, os.Getenv("DB_URL"))

	// Regular dependency
	dep := dependency.New(newUser, driver)

	// Singleton dependency
	dep2 := dependency.NewSingleton(newUser, driver)
}
```

Example of Injectable dependency as argument:

```go
package main

import (
	"os"

	"github.com/Drafteame/container"
	"github.com/Drafteame/container/dependency"
)

func main() {
	driverName := "database"
	driver := dependency.New(dbConstructor, os.Getenv("DB_URL"))

	if err := container.Register(driverName, driver); err != nil {
		panic(err)
    }
	
	// Regular dependency
	dep := dependency.New(userConstructor, dependency.Inject(driverName))

	// Singleton dependency
	dep2 := dependency.NewSingleton(userConstructor, dependency.Inject(driverName))
}
```

### Invoke

There is a method that can help you to bring some extra functionality to the container and obtain more than one instance
at a time.

This method will receive a callback that can or not return an error and can or not receive multiple arguments. This
arguments should be structs, defining on his fields the instances that the container should inject to it.

```go
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Drafteame/container"
	"github.com/Drafteame/container/dependency"
	"github.com/Drafteame/container/types"
)

type args struct {
	types.In
	User *user `inject:"name=user"`
}

func invoker(in args) error {
	if in.User == nil {
		return errors.New("empty instance of user")
	}

	fmt.Println("Hello ", in.User.GetName())
	return nil
}

func main() {
	driverName := "database"
	driver := dependency.New(newDB, os.Getenv("DB_URL"))

	if err := container.Register(driverName, driver); err != nil {
		panic(err)
	}

	depName := "user"
	dep := dependency.New(newUser, dependency.Inject(driverName))

	if err := container.Register(depName, dep); err != nil {
		panic(err)
	}
	
	if err := container.Invoke(invoker); err != nil {
		panic(err)
    }
}
```

Also you can use interface segregation to define the arguments:

```go
package main

import (
	"errors"
	"fmt"
	"os"

	// ...
	
	"github.com/Drafteame/container"
)

type namer interface{
	GetName()
}

type args struct {
	types.In
	User namer `inject:"name=user"`
}

func invoker(in args) error {
	if in.User == nil {
		return errors.New("empty instance of user")
	}

	fmt.Println("Hello ", in.User.GetName())
	return nil
}

func main() {
	// .....

	if err := container.Invoke(invoker); err != nil {
		panic(err)
	}
}
```

#### Optional arguments

When you define In structs to be used with the `Invoke` method you can mark optional fields if you expect that some
fields can or not be filled by the injector and avoid an error if there's no dependency registered with the required
name.

```go
package main

import (
	"errors"
	"fmt"
	"os"

	// ...
	
	"github.com/Drafteame/container/types"
)

type namer interface{
	GetName()
}

type args struct {
	types.In
	User namer `inject:"name=notExist,optional"`
}

func invoker(in args) error {
	if in.User != nil {
		fmt.Println("Hello ", in.User.GetName())
	} else {
		fmt.Println("Ups no namer instance found")
    }
	
	return nil
}

func main() {
	// .....

	if err := container.Invoke(invoker); err != nil {
		panic(err)
	}
}
```
