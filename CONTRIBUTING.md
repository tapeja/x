# Contribution Guidelines

Please follow these guidelines, if you intend to contribute code to the project.

The document is heavilly inspired by :
- [Go's common list of
  mistakes](https://github.com/golang/go/wiki/CodeReviewComments)
- [ThoughtBot's style
  guide](https://github.com/thoughtbot/guides/tree/master/style)

**If you would like to suggest changes to these guidelines,** open an issue or
pull request with the `discourse` label and propose a change.

---

### General

#### Documentation

- Be verbose while documenting.
	- What is the intent?
	- What kind of errors can be expected?
	- Any special design considerations?
	- Any difficulties while implementing?

#### Testing

- Aim for 100% test coverage.
- Document tests just like the rest of the codebase.

#### Formatting

- Break long lines after 80 characters.
- Delete trailing whitespace.
- Use proper capitilazation, punctuation, and spelling.

#### Naming

- Avoid abbreviations.
- Avoid object types in names, i.e. `UserSlice`, `CarStruct`, `api_package`.

---

### Git

- Squash trivial commits into a single commit.
- Write [good commit
  messages](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).
	- Summary line should aim for 50 characters, with 72 as a max.
	- Follow summary line with an empty line.
	- Details follow the empty line and should break every 72 characters.
	- Be verbose, describe the change.
	- Don't get cute ;) i.e. emoticons.
- Avoid merge commits by using rebase, i.e. `git pull --rebase`.
- Avoid branch pollution in the main repository, work on a personal fork and
  file a pull request when ready.
- Avoid polluting the `.gitignore` with personal settings, i.e. ignoring IDE or
  OS specific files.
- Think twice before placing a binary file under version control.

---

### Go

- [Write idiomatic Go](https://golang.org/doc/effective_go.html), i.e. don't
  use `self` or `this` as a receiver.
- Run [`goimports`](https://godoc.org/golang.org/x/tools/cmd/goimports) to
  format your code before filing a pull request.
- [Pass values, unless the value is a large `struct` or the semantics of the
  method requires a
pointer.](https://github.com/golang/go/wiki/CodeReviewComments#receiver-type)
- [Follow these rules of thumb when choosing a receiver
  type.](https://github.com/golang/go/wiki/CodeReviewComments#receiver-type)
- [Avoid using dot
  imports.](https://github.com/golang/go/wiki/CodeReviewComments#import-dot)

#### Errors

- Do not discard errors using `_` variables. If a function returns an error,
  check it to make sure the function succeeded. Handle the error, return it, or,
in truly exceptional situations, panic.
- Don't use panic for normal error handling. Use error and multiple return
  values.
- Error strings should not be capitalized.
- [Indent error
  flow.](https://github.com/golang/go/wiki/CodeReviewComments#indent-error-flow)

#### Naming

- [Use mixed
  caps.](https://github.com/golang/go/wiki/CodeReviewComments#mixed-caps)
- [Use consistent acronyms, i.e. URL should be either URL or url never
  Url.](https://github.com/golang/go/wiki/CodeReviewComments#initialisms)
- [Omit package names from
  identifiers.](https://github.com/golang/go/wiki/CodeReviewComments#package-names)

#### Documentation

- Use the [godoc convention](http://blog.golang.org/godoc-documenting-go-code)
  for documentation.
- All top-level, exported names should have doc comments, as should non-trivial
  unexported type or function declarations.

#### Tests

- Use [Go's native testing package](http://golang.org/pkg/testing/), only use a
  testing framework if it becomes absolutely necessary.
- Use [table driven tests](https://github.com/golang/go/wiki/TableDrivenTests)
  when appropriate, instead of writing separate unit tests.
