# designlanguage
Nonzero Sum Design Language definition and tools

## Overview

Software can be seen to have structural elements as well as aspects that are more fluid. It can be useful when defining the structure of an application to write down its design. This design might then be implemented in multiple languages and platforms.

The Nonzero Sum Design Language is a part of the [Nonzero Sum Stack](https://nonzerosumsolutions.com/stack/index.html). The Nonzero Sum Stack is a set of recommended technologies and engineering processes for sustainable software development. For examples of apps built using the Stack (including Design Language) see [helloapp](https://github.com/mattmunz/helloapp) (Go) and [helloapppy](https://github.com/mattmunz/helloapppy) (Python).

Design Language is a language and toolset for working with software designs. By automating the process of developing software structure from an abstract design, the cost of developing and maintaining an application can be reduced.

This workflow can nicely complement a workflow that includes AI code generation, by removing from the AI coder the task of generating the application structure. This allows the agent and the AI to focus on the high-level design of the application without wasting effort on the generation of the application structure, which can be easily automated with the design language tools. Once the structure is laid down, both human and AI coders can fill in the details, confident in the accuracy and reliability of the overall application structure.

The idea of the design language toolchain is that it is not specific to any (other) programming language. It will be used to produce the structural elements of programs written in many different programming languages. Unlike similar frameworks like UML it has a very simple syntax. This simplicity is a key feature. The language is intended to be more concise than most data structure syntaxes for most programming languages. There is room to make it more concise over time.

## Getting started

### 1. Create your design

Here's an example design for a commandline application.

```nzsd
User

Console
* Version String
* CommandHandlers []CommandHandler

CommandHandler
* Description String
* Interpret (Command String) -> (Response String)
* GetConfigs (Time DateTime) -> (Configs []String)

Browser
* GetElement (ID String) -> (Element Element)
```

### 2. Save your design

Your design should be saved to a path matching `$PROJECT_DIR/documentation/design/$DESIGN_NAME/$DESIGN_NAME.nzsd.txt`. For example, see the example [Geometry design](test/data/project1/documentation/design/geometry/Geometry.nzsd.txt).

### 3. Install

```bash
go install github.com/mattmunz/designlanguage
```

### 4. Go to your project directory

```bash
cd $PROJECT_DIR
```

### 3. Test run

First print the code to the commandline.

In Go...

```bash
designlanguage gen -t go
```

Or Python...

```bash
designlanguage gen -t python
```

Review the code and address any issues in the design that arise.

### 4. Generate code

```bash
designlanguage gen -d=false -t <language>
```

### 5. Update your code to use the generated model code

(The generated code can be found under `$PROJECT_DIR/model/gen/$DESIGN_NAME)`.)

### 6. Compile and run your application

Presumably this is something like...

In Go...

```bash
go build
go install
myapp 
```

in python...

```bash
python3 my_app.py
```

## Syntax

The current syntax is markdown-inspired, with elements that may be familiar to Gophers or Haskell aficionados. It is intended to be simple, intuitive, and easy to read.

In the future alternate syntaxes may be supported, as needed by the user base. Any syntax will be required to be simple to use and parse, and to fully support the underlying model of the langiage.

## Getting Involved

* Try Design Language on one of your projects and share how it went.
* Pick an item from the TODO section to implement. Make sure to get in touch with the project before starting anything ambitious!

## FAQ

### What target languages are supported?

Currently, Go and Python are supported but support for other languages is planned. Most likely the next one will be Javascript/Typescript. If you have a language that you would like to be supported please [get in touch with us](contact).

### Can changes be made to the syntax?

Sure. Conceivably multiple front ends can be used to allow multiple syntaxes. Simple changes to the default syntax may also be considered.

### Why not use UML?

UML is rather complex. Design Language is intended to be extremely simple, so as to avoid operational errors. 

### Doesn't OpenAPI do this?

OpenAPI is focused on data exchange. Design Language is more general, focused on application structure. It is possible to integrate Design Language with OpenAPI, for example by generating OpenAPI model object from Design Language designs.

### What about product X, which does something similar?

It would be great for us to learn about similar products. Please [contact us](contact) to tell us all about it.

## Is this a tool to guide LLMs?

In a manner of speaking, yes. Ideally LLMs would be used to generate software descriptions in the design language, and the tool would be used to generate working code deterministically. From that point LLMs could be used to write the non-structural parts of the codebase.

## Why generate at all?

Simple structures and interfaces are easy to write in Go, Python Javascript, etc. So why introduce this special language for them?

Simply put, the structure of your application is a special and important aspect of what you are building. By using a specialized language, you can avoid errors. By using a toolchain that supportas multiple targets, you can harmonize multiple cooperating applications written in different languages. As the norm these days is multi-lingual programming, having a rock-solid core to build upon is quite important!

## Contributing

Contributions are welcome! Please [contact us](contact) to learn about the best ways to contribute, report bugs, submit fixes, plan new features, etc.

## Contact

See the [Nonzero Sum contact page](https://nonzerosumsolutions.com/contact.html) for ways to get in touch with us.

## TODO

See [TODO](documentation/Todo.md)
