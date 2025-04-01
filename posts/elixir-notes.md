# Elixir Notes

These are notes on the Elixir programming language.

## Basics

Elixir is a general purpose, functional language that executes on the BEAM virtual machine (similar
to Erlang).

- It runs on the BEAM.
- The language is generally interoperable with Erlang and Erlang libraries (Erlang libraries may be
  referenced by `:lower_case_atoms`).
- `iex` is an interactive shell to execute ad-hoc Elixir statements.
- Elixir uses first-class documentation syntax, including `#` for inline, and `@doc` and
  `@moduledoc` for annotated documentation for functions and modules, respectively.

Documentation is at <https://hexdocs.pm/elixir>.

## Types

-   **Integers**, which may include binary, octal, or hexadecimal.

-   **Floats**, which also support `e` notation.

-   **Atoms**, similar to Erlang atoms, is a constant whose value is the name. They are prefixed with
    `:`, e.g. `:foo` or `:bar`.

-   **Booleans**, which are `true` and `false`. Everything evaluates as `true` except for `false` and
    `nil`. Boolean values may also be written as atoms, `:true` and `:false`. `||` and `&&` and `!`
    are relevant operators . `and` and `or` also exist but require the first argument be a boolean.

-   Strings are represented as **binaries** and use "double quotes" and are UTF-8.

    String interpolation can be done with `#{}`.

        name = "Jimbo"
        IO.puts("Hello #{name}")

    Strings may be concatenated with `<>`.

    **Binaries** (and thus, strings) under the hood are sequences of bytes. This can be made explicit
    by using `<< >>`. Digits in byte sequences are UTF-8 character codes.

        <<104,101,108,108,111>>  # "hello"

    Strings may also be represented as character lists. Strings using 'single quotes' are defined this
    way.

        'hello' = [104, 101, 108, 108, 111]

-   **Lists** can contain multiple types.

        [123, "abc", :foo]

    All collection types are also enumerables, with the exception of tuples, and can use functions in
    the Enum module. Enum functions are bread and butter functions for functional programming (e.g.
    `any?/2`, `all?/2` `each/1`, `map/2`, `sort/1`).

### Lists

List concatenation uses `++`.

    [123] ++ [123, "abc"]  # [123, 123, "abc"]

The `--` operator may be used to remove matching list values, from left to right, using strict
comparison.

    [1, 1, 2, 1, 2] -- [1, 2, 1]  # [1, 2]

Like other functional languages, referencing the head or tail of a list is common. The `hd` and `tl`
functions accomplish this.

    hd [1, 2, 3]  # 1
    tl [1, 2, 3]  # [2, 3]

Head and tail may be assigned to variables using the `|` operator and pattern matching.

    [h | t] = [1, 2, 3]  # h is 1, t is [2, 3]

Lists are stored internally as linked lists, so traversal and appending may be expensive.

### Tuples

Unlink lists, tuples are stored in contiguous memory, making reads fast but modification expensive.
They are commonly used as return values for functions.

    {:ok, 123}
    {:error, nil}

### Keyword Lists

Keyword lists are a type of ordered, associative array where keys are atoms. They are commonly used
to pass arguments to functions. Like lists, keys need not be unique.

    [foo: "bar", fizz: "buzz"]  # or ...
    [{:foo, "bar}, {:fizz, "buzz"}]

### Maps

Maps are similar to keyword lists, but keys may be any type and they are unordered. Keys must be
unique.

    map = %{:foo => "bar", "fizz" => :buzz}
    IO.puts(map[:foo])  # or ...
    IO.puts(map.foo)

Maps can be updated using `|` (keys must already exist).

    map = %{:foo => "bar", "fizz" => :buzz}
    map2 = %{map | :foo => "derp"}

`Map.put/3` is used to update maps.

## Pattern Matching

The `=` operator is not just for assignment, but also for matching. Matching returns the value of
the expression if both sides match (which always happens during assignment), or an error.

    x = 1
    1 = x
    a_list = [1, 2, 3]
    [1, 2, 3] = a_list

Assignment can be avoided with the pin operator, `^`.

    x = 1
    ^x = 2  # MatchError

The pin operator can also affect function clauses.

Note here that `some_thing` referenced in the function is dependent on the `some_thing` variable
existing in context, and that `_` prefixing it is used because the variable is technically unused in
the function.

    some_thing = "foo"
    some_function = fn
        (^some_thing) -> "I got foo"
        (_some_thing) -> "I didn't get foo"
    end

## Conditionals

Conditionals come from the Kernel module. There is `if` and `unless`.

    if true do
        "It's true"
    else
        "It's false"
    end

    unless false do
        ...
    end

There is also `case`. A failure to match anything will raise an error.

    case status do
        {:ok, result} -> result
        {:error} -> "Unfortunate"
        _ -> "Something weird happened"
    end

Case can also use guards.

    case {1, 2, 3} do
        {1, 2, x} when x < 10 ->
            "This matches"
        _ ->
            "No match"
    end

There is also `cond` that can match against conditions and `with` which can work easily with
compound clauses.

## Comprehensions

Elixir has comprehensions to manage enumerables.

    for x <- [1, 2, 3, 4, 5], do: x*2  # [2, 4, 6, 8, 10]

    for {k, v} <- [a: 1, b: 2, c: 3], do: IO.puts("#{k} : #{v}")

    for {k, v} <- %{"a" => 1, "b" => 2, "c" => 3}, do: IO.puts("#{k} : #{v}")

Comprehensions operate from generators which specify the next value. Generators use pattern matching
and ignore non-matches.

    for {"a", val} <- %{"a" => 1, "b" => 2}, do: val  # Only prints '[1]'

Multiple generators may be used in a comprehension.

Comprehensions may use a type of guard called filters.

    for n <- 1..100, is_odd(n), do: n

The output of a comprehension is always a list, but this can be changed by using `into`.

    for n <- 1..100, into: "", do: <<n>>

## Functions

Functions have a name and an arity (number of arguments). The proper way to refer to a function in
Elixir is to combine them, such as `foobar/2`, which represents a unique function in an Elixir
module.

Functions may be declared within modules (i.e. a "named" function). Note that the calling syntax is
different--(`func.()` vs. `Module.func()`.

    defmodule MyModule do
        def my_function(number), do: number * 2
    end

    MyModule.my_function(123)

Anonymous functions are often used as function arguments in transformations. It can be declared with
`fn` or use the capture operator, `&`.

    Enum.map([1,2,3], fn number -> number * 2 end)  # or ...
    Enum.map([1,2,3], &(&1 * 2))

The capture operator can also be used to reference previously declared functions.

    Enum.map([1,2,3], &MyModule.my_function(&1))  # or ...
    Enum.map([1,2,3], &MyModule.my_function/1)

Functions are first-class and may be assigned to variables.

    multiply_by_two = &(&1 * 2)
    Enum.map([1,2,3], multiply_by_two)

Functions may be invoked with `.`. Parenthesis are (mostly) optional, although they should probably
be used.

    multiply_by_two.(2)

Functions may be chained with the pipe operator, `|>`. This automatically provides the output of a
previous function as the first parameter of the second.

    "abcdefg" |> foo.() |> bar.("fizzbuzz")  # "fizzbuzz" is the second argument of bar

Functions may employ pattern matching directly from arguments.

    defmodule RecursiveCounter do
        def count([]), do: 0
        def count([_ | tail]), do: 1 + count(tail)
    end

Functions with the same name and arity can use guards to distinguish input.

    defmodule MyModule do
        def hello(names) when is_list(names) do
            ...
        end

        def hello(name) when is_binary(name) do
            ...
        end
    end

Function overloading in Elixir is equivalent to defining functions with the same name and different
arities.

    defmodule MyModule do
        def my_function(number), do: number * 2
        def my_function(numberA, numberB), do: numberA * numberB
    end

Pattern matching collections can be done by selecting expected members and using variable assignment
to capture more than the match.

    def hello_from_map(%{name: a_name} = a_person) do
        IO.puts "Hello, " <> a_name
        IO.inspect a_person
    end

    hello_from_map.(%{age: 38, name: "Jimbo"})

Functions are private if declared with `defp`, which only allows inter-Module usage.

Default arguments are defined using `\\` syntax. Note that this takes some special semantics when
combining use of default arguments with guards.

    def hello(name, greeting \\ "Hello") do
        IO.puts greeting <> ", " <> name
    end

Functions will always return the last evaluation.

    def hello() do
        _ = "Hello!"
    end

    hello.()  # Will return "Hello!"

Function return convention is to provide either a `{:ok, result}` or `{:error, reason}` tuple.

## Modules

Modules are the namespacing mechanism for functions. The calling syntax for a function in a module
is `<Module>.<function>(<args>)`, provided the function is not private.

Besides functions, modules may contain other things:

-   **Structs** are pre-defined maps which hold data for the module instance. Values provided in
    `defstruct` act as defaults.

        defmodule MyModule do
            defstruct foo: "bar", fizz: "buzz"
        end

        an_instance = %MyModule{foo: "bar2"}

    Structs may also be defined using the `__struct__` helper.

        defmodule MyModule do
            defstruct foo: "bar", fizz: "buzz"
            def new(args), do: __struct__(args)
        end

        MyModule.new(foo: "bar2")

    Structs can be compared to and updated like maps.

-   **Imports** which can import functions from mother modules.

        defmodule MyModule do
            import AnotherModule, only: [some_func: 1]
            ...
        end

    Importing modules allows direct use of their functions.

        import Integer
        is_even(2)  # actually Integer.is_even/1

-   **Alias** can be used for composing a module with another.

        defmodule MyModule do
            alias AnotherModule.Thing
            ...
        end

-   **Exceptions** to define custom exceptions.

        defmodule MyError do
            defexception message: "Well this sucks"
        end

Modules can also include other composition and meta-programming constructs.

## Errors

Functions should try to return `:ok` and `:error` tuples, but Elixir can also use exceptions.

    raise "This is an exception"
    raise ArgumentError, message: "This is an exception"

Exceptions can be handled using `try` clauses. Note that `RuntimeError` and `ArgumentError` are two
of Elixir's built-ins.

    try do
        raise "This is an exception"
    rescue
        e in RuntimeError -> IO.puts("We caught an error")
        e in ArgumentError -> IO.puts("We caught another type of error")
    after
        IO.puts("This runs regardless of whether or not we got an error")
    end

There is also `throw` and `catch`, used less often.

## Concurrency

Elixir employs a concurrency model very similar to Erlang's, based on the Actor model. The BEAM has
a special concept of a process (not the same as an OS process).

    spawn(MyModule, :my_function, [arg1, arg2])  # Create a new BEAM process

Processes may pass messages to each other. Pattern matching in receive clauses will skip messages
that do not match. Sending and receiving do not block.

    defmodule MyModule do
        def listen do
            receive do
                {:ok, payload} -> IO.puts(payload)
            end
            listen()  # Re-execute listener to receive another message
        end
    end

    pid = spawn(MyModule, :listen, [])
    send pid, {:ok, "hello"}

These are the very basic building blocks of Actor model concurrency. More sophisticated things are
also provided by the language to manage process management and state (e.g. the `Agent` module, or
async/await provided by the `Task` module).

Elixir also has access to Erlang's OTP model of designing distributed systems.
