# Go Notes

These are notes on the Go programming language.

## Basics

Go is a statically typed, general purpose, imperative language that is compiled, garbage collected,
and features memory referencing, cross-compilation, and concurrency.

- Style is enforced by the `go fmt` tool. Go expects certain styling to compile.
- Each Go file declares a package.
- All executable Go programs have a main function in a main package, otherwise they are libraries.
- Uppercase methods are public and exportable, lowercase methods are private.
- Go constructs strings from **runes**, which are unicode characters. Unicode characters are also
  useable in the language syntax itself.
- Go uses pointers and pointer de-referencing.

Documentation is at [golang.org](https://golang.org). It's also available using the `go doc` tool.

    go doc <module> <function>

Go source is compiled.

    go build cmd/project/main.go

`GOPATH` holds Go dependencies. It can be anything, like the current working directory, or something
like`$HOME/lib/go`.

    cd ./project
    export GOPATH=`pwd`
    go install

Go dependencies are declared and managed as Go modules, essentially a declared collection of
packages. These are often managed by the Go toolchain in a `go.mod` file. Packages may be explicitly
retrieved and can be detected within code using `go mod`.

    cd ./project
    go mod init project
    go get <package>
    go mod tidy

Projects with a `main` package and a `main.go` file can be directly executed with explicit
compliation (they are still compiled).

    go run main.go

## Variables

Go has static types with hinting. You can explicitly declare types or you can have Go figure it out.

The following are equivalent:

    var message string
    message = "Hello"

    // or

    var message string = "Hello"

    // or

    message := "Hello"

The assignment operator `:=` will also allow you to redeclare some variables, as long as part of the
declaration statement includes new variables. Redeclared variables must respect type and use `=`.

    old_var := 3
    new_var, old_var := 3, 4
    new_var = 5

Outside functions, `const` blocks declare global constants and `var` blocks declare global
variables. Capitalized constants and variables export them for use by outside modules.

    const (
      foo = "bar"
    )

    var (
      Fizz = "buzz"
    )

The `_` variable can be used to throw away unneeded information.

    var foo, _ = myFunc()

Go will error on unused or over-declared variables.

All variables, when declared, are initialized to their "zero value," which depends on the type
(e.g., booleans get `false` and integers get `0`).

## Types

Types can be defined using type functions.

    var number = uint(9)  // 9 is an argument to the uint type function.

Type functions can also do some conversions. Each conversion creates a copy.

    my_string := "This is a test"
    my_byte_slice := []byte(my_string)
    back_to_string = string(my_byte_slice)

Some times include:

- **int** and **uint8** and **unint64** are integers.

- **float** and **float32** and **float64** are floating point types.

        var pi float64 = 3.1415926
        pi = float64(3.14159)

- **string** is a string.

- **bool** is a boolean. `!` is negation.

- Also **byte**.

- **iota** can create constant enumerations and is an auto incrementing value.

        const (
          zero int = iota
          one
          two
        )

        fmt.Printf("%d %d %d", zero, one, two)   // Prints '0 1 2'.

        const (
          zero int = iota * 4
          four
          eight
        )

        fmt.Printf("%d %d %d", zero, four, eight)   // Prints '0 4 8'.

- Also **error**--the built-in Go error type.

- The empty value in Go is `nil`.

- A "catch-all", generic type is `interface{}` or `any`, which can be used to variable multiple
  types.

## Defined Types

**Structs** can be used to define types.

    type something struct {
      foo string
      bar int
    }

    s := &something{foo: "test"}
    s.bar = 3

Both examples below will instantiate a struct with zero-values.

    s := &something{}
    s := new(something)

Types and the `type` keyword can also be used to alias and extend existing
types with methods.

    type SuperIntSlice []int
    s := &SuperIntSlice{1, 2, 3, 4, 5, 6}

## Strings

Strings can be indexed using array brackets and `:`.

    astr := "This is my string.\n"

    astr[0:4]
    astr[:4]

The `range` keyword allows k/v iteration.

    for _, r := range astr {
      ...
    }

Looping over a string rune by rune:

    for i, r := range astr {
      fmt.Printf("%d %c\n", i, r)
    }

Get string length with `len()`.

Strings in back quotes, \`, are interpreted literally. Double quotes, `""`, will interpret special
characters, like `\n`.

## Arrays and Slices

Go handles normal array behavior with **arrays** and **slices**.

Declaring arrays can be done bounded or unbounded.

    unbounded := [...]string{"foo", "bar"}
    bounded := [2]string{"fizz", "buzz"}

Slices are a specific sections of arrays. They are defined using empty `[]`.

    my_slice := []string{"the", "quick", "brown", "fox"}

Slices can be formed on the fly from arrays or other slices. The first value is the starting index
(inclusive), the second value is the ending index (exclusive).

    foo = unbounded_slice[1:2]  // Includes only unbounded_slice[1]
    bar = unbounded_slice[:3]   // Includes the first three items.

Arrays are fixed in memory and are only passed between functions by copying, whereas slices are
passed by reference. Slices are more performant when passed between functions because it doesn't
copy data, and allows array data to be directly modified outside of its declaration scope.

Dynamic slices can be created using the `make` function, which allocates memory.

    words := make([]string, 4)      // Initializes with four empty strings.
    words := make([]string, 0, 4)   // Intiializes with no entries, but a max
                                    // capacity of four.

    words = append(words, "Added string")

- `len()` will get length and `cap()` will show capacity.
- `append()` will automatically expand slice capacity if it's exceeded.
- Slices can be explicitly copied using `copy()`.

Attempting to access sections of a slice or array that do not exist will cause "out of bounds"
runtime errors.

Slices have no equality property, so they cannot be compared, except to `nil`.

## Maps

Maps can be instantiated using the `map` keyword, either by `make` or by explicit definition.

    my_map := make(map[key_type]value_type)

    my_other_map := map[key_type]value_type{
        "foo": "bar",
        ...
    }

Accessing an element of the map that doesn't exist will return the zero type. An error can be
explicitly thrown by checking for a return value.

    foo, ok := my_map["bar"]  // "bar" doesn't exist.
    if !ok {
      // This was an errorneous lookup.
    } else {
      // The lookup worked.
    }

Iterating over a map uses `range`. Maps are not sorted by default, so results will come out in
arbitrary order.

    for month, days := range days_in_months {
      fmt.Printf("There are %d days in %s.\n", days, month)
    }

Adding a value to a map just requires an assignment.

    my_map["foo"] = "bar"

Items can be deleted from maps. Deletion will silently pass if the key does not exist.

    delete(my_map, "foo")
    delete(my_map, "foo")  // Does not delete, does not error.

Keys can be any type that has equality.

## Conditionals

Go can use if statements for error handling. Variables declared within if statements are only
accessible within the conditional block.

    if readCount, err = module.Format("Hello.\n"); err != nil {
      // Respond to an error.
    } else if readCount == 0 {
      fmt.Printf("Read nothing!")
    }else {
      fmt.Printf(readCount)
    }

Go also has switch statements.

    switch {
    case foo == "bar":
      ...
    case foo == "bizz":
      ...
    default:
      os.Exit(1)
    }

Switch statements do not have a necessary `break` and do not fall-through. Instead, the
`fallthrough` keyword will cause cases to collapse.

    case foo == "bar":
      ...
      fallthrough

You can also switch on a particular variable.

    r = someRandomCharacter()

    switch r {
    case "a", "e", "i", "o", "u":
      fmt.Printf("r turned out to be a vowel.")
    default:
      fmt.Printf("r is weird!")
    }

## Loops

The `for` loop covers all cases.

An infinite loop:

    for {
      fmt.Printf("Hello, world.\n")
    }

A constrained loop:

    counter := 1
    for counter < 100 {
      counter += 1
    }

    // or ...

    for counter := 0; counter < 10; counter++ {
      fmt.Printf("Foo bar.\n")
    }

Go does have post-fix updating, `foo++`, but it isn't a statement and is only valid inside loop
clauses.

Loops can simultaneously use assignments and multiple conditions.

## Functions

Calling a function takes this form:

    <package>.<function>([args ...])

A function definition looks something like this:

    func name(param1 type, param2, param3 shared_type) (return_type1, return_type2) {
      ... statements ...

      return foo, bar
    }

- Functions can return multiple things.
- Parameters declare their type after the variable name in the argument list.
- Return type parentheses `()` are optional if there is only one return type.

The `defer` keyword queues up actions to be taken after a function exits. Statements which are
defered execute in stack order. Useful for function clean-up.

    func printer(msg []byte) error {
      f, err := os.Create("helloworld.txt")
      defer f.Close()

      if err == nil {
        f.Write(msg)
      }

      return err
    }

Return values in Go can also be defined in the function clause for shorthand returns.

    func tester(msg []byte) (e error) {
      _, e = fmt.Printf("%s\n", msg)
      return  // Automatically returns `e`.
    }

Go supports variadic functions with the `...` parameter.

    func tester(fmt string, msgs ...string) {
      for _, msg := range msgs {
        fmt.Printf(fmt, msg)
      }
    }

In general, pass-by-reference is accomplished by using pointers. Otherwise, Go usually uses
pass-by-copy. Some types will always pass-by-copy (e.g. structs).

## Errors

Go has errors which often return from functions and also (the very rarely used) panics to handle
extreme events.

    func my_printer(msg string) error {
      _, err := fmt.Printf("%s\n", msg)
      return err
    }

    if err := my_printer("Hello, world."); err != nil {
      // The function returned an error.
      os.Exit(1)
    }

By default, errors use the built-in standard library. Defining errors can be done using error
formatting, `Errorf`, or the `errors` library.

- `Errorf` can be used to say unique error messages.

        fmt.Errorf("This is an error message!")

- The `errors` library can be used to define custom error types. Downstream behavior can check the
  type of error and respond accordingly.

        var (
          errorTerribleThing = errors.New("Something terrible happened!")
        )

        func do_something() error {
          ...
          return erroTerribleThing
        }

`panic` and `recover` in Go can stop execution or return stack traces if something happens that is
not recoverable.

## Channels and Goroutines

Channels and Goroutines are two features Go offers to manage concurrency.

**Channels** are typed, representing what they contain, and can act as pipes. They can be buffered
or unbuffered, which have different blocking behaviors. Channels generally behave like synchronized
queues. Sending to and receiving from channels use the `<-` operator.


**Goroutines** are lightweight threads of execution that execute a function concurrent to the thread
which launched them. They are triggered with the `go` keyword.

Below is an example program employing both. Note that there is no inherent synchronization or "wait"
mechanism, hence the sleep.

```go
package main

import (
	"fmt"
	"time"
)

func consume(numbers chan int, prefix string) {
	var next int

	for {
		next = <-numbers
		fmt.Printf("Got number %d from consumer %s\n", next, prefix)
	}
}

func main() {
	var numbers chan int = make(chan int, 10)
	defer close(numbers)

	for i := 0; i < 10; i++ {
		fmt.Printf("Put numbers %d\n", i)
		numbers <- i
	}

	go consume(numbers, "a")
	go consume(numbers, "b")
	go consume(numbers, "c")

	time.Sleep(10 * time.Second)
}
```
