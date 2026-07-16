## TODO

### Major features


* Make Python code based on generated models more idiomatic
* See about adding more layers of abstraction for the code generation bit
* This python AST is not good for go reflection. Make more particular.
* Extract type mappings to config file so that they can be customized
* Find all the structs in the code and add them to design language files for this project
* Javascript target
* Typescript target
* WASM target
* Add default implementations for entities. Implementations will use best practices and design patterns.
* Interoperate with OpenAPI and other messaging systems
* Interoperate with Database systems
* Support alternate syntaxes as needed

### Other features

* Implement designlanguage-based flow for ClimateComms, PersonalKB, etc.
* Consider making parser less strict
  * Allow trailing space
  * Allow no space before parens

## Done

* Python target
* Find an intro article to publish: https://nonzerosumsolutions.com/library/gettingStarted.html
* Find a second article to publish: https://nonzerosumsolutions.com/library/sustainability.html
* Make documentation: X Getting started, X Contriubuting, etc.
* As part of this, make example repo of a hello world service. hellocli
* X Use it for an existing project -- appkit -- generate model interfaces
* X ANTLR parser no longer worth the trouble. 
* X Add inline documentation in the source file ending up in the output source code
* X comment of type in first line of body
* X comments after -- mark on line
* X Document comment including author attribution
* X Implement all the appkit syntax. TDD
* X Make a single-pass parser, line by line, with context
* X TDD test: entity with super interface, object with superinterface, renders to superType as first field
* Make a parser interface
* Move appkit to its own repository, and later rework sundries to use it.
* X Make an ANTLR grammar file
* X Generate a Go parser
* X Make a simple test and get it to work
* X Make more complex tests and get them to work
* X Make a gen command to parse design to model, and then write out interfaces
* Write A) NZS article introducing designlanguage, appkit, and helloapp but don't publish
  * Update B) NZS Stack page but don't publish
* Add enum type to model and go generator
ExprContext = {Load, Store, Del, AugLoad, AugStore, Param}
    * X Try this AST: https://github.com/go-python/gpython/blob/main/ast/ast.go
    * X AAAAH this AST Sucks. DDefine the Python AST in design language, generate, go from there. MAybe can advertise Python AST as a separate project
  