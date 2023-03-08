# Future-proof Go packages

## Elevator Pitch

Future-proof packages don't break their surface area.
They grow and evolve without needing to change code that consumes them.
They can be part of a library or a large application.
In this talk, we will cover how you go about writing future-proof packages.

## Description

Future-proof packages are well-designed Go packages that don't break their
surface area, and act almost like independent libraries in whatever codebase
that they're in.
They define clean boundaries and contracts with the rest of the codebase,
and therefore make it easy to grow, extend, and evolve the package
without needing to change code that consumes them.
They exemplify the idea that testing is a design problem:
anything can be tested if you design for it.

In this talk, I will discuss requirements for future-proof packages
(flexibility and maintainability among them)
and practices you can adopt in your code to make your packages future-proof.

## Notes

### Why I'm qualified

I believe I'm a good person to speak on this topic because
I've spent 5+ years writing, maintaining, and extending libraries in Go.
These libraries were used regularly by a large user base
(the majority of the Go engineers at a large company)
with knowledge and expertise in Go ranging from beginner to advanced.
The key point I want to focus on in the talk is maintainability of the public API of a package.
Maintainability was a hard requirement for libraries that I wrote or maintained
because breaking changes were not an option.
Code maintainability was also the area I focused on in the internal educational
materials I produced for engineers at the company.

### Rough outline

Things I intend to cover go from higher level design advice like the following
(not exhaustive; may change):

- Starting with the public interface:
  Don't worry about what code already looks like.
  Start with what you want the usage to look like,
  and then make reality match what you want.
- Features are not additive:
  You can't just "add" a new feature.
  How does the feature interact with the other features?
  What happens when I set both those flags to true?
  It's not an addition, but a cross product.
- Don't leak implementation. Define the domain language.
  Internals store a JSON object? Don't leak that representation outside.
  Map to a domain-specific object that you expose as part of the API.
  Leave yourself room to change the JSON shape or migrate to Protobuf
  without changing the public API.

Going into Go-specific design thoughts:

- Accept interfaces, return structs: This is a popular Go adage,
  so I'll cover the what and why of it quickly.
- Build on top of a small interface:
  This will echo the "The bigger the interface, the weaker the abstraction."
  Go proverb and help programmers expand on it.
  Make a small interface with the bare minimum functionality -- the most
  powerful base piece upon which you can build the rest.
  Then wrap it with a struct and build convenience methods on that struct
  that "translate" to that interface. (This is best demonstrated with an
  example. One example from the standard library is how `http.Client` is a
  wrapper around the small `http.RoundTripper` interface.)

To more concrete recommendations and patterns like the following
(will include **short** demonstrative or comparative code samples):

- Errors:
  Expose sentinel errors for leaf errors, and error structs for *meaningful*,
  structured information about the failure (e.g. here's the object the
  AlreadyExistError conflicted with.)
  If it's just a wrapper around another error, use `fmt.Errorf` with `%w`.
  Tests should be able to `errors.Is` or `errors.As` the original failure in
  that case.
- Goroutines: Don't start without permission. Don't ever spawn in `init()`.
- Channels that cross API boundaries should be rare:
  Keep these an internal detail unless there's absolute need to expose a
  channel.
  If you do, make them uni-directional (either read or write).
  If you expose a write channel, make sure the caller knows to close it
  (or to not close it).
  Better: don't expose a channel and use a callback. Which brings us to:
- APIs that accept callbacks (`onFoo func(*Foo)`) are better off accepting
  interfaces (`Fooer interface{ OnFoo(*Foo) }`).
  This leaves room for upgrades (with upcasting) in the future.
- You can always add more fields to a struct:
  Make a struct to hold options for a method,
  and you'll never break that method signature again.
- Functional options:
  (There's plenty of literature on this. I won't repeat that here but I'll
  provide an overview in the talk and when you might use them.
  I'll mention that you can do closure-based options and interface-based.
  Then after the next section, I'll bring it back with "therefore, prefer
  interface-based.")
- Prefer objects to closures:
  This ends up being more about internal maintainability of the package rather
  than the public API, but these closures often pass through API boundaries
  (e.g. as callbacks from another library).
  Objects are comparable, closures are not.
  Making them into objects make them easier to test and maintain,
  and makes it easy to expose a meaningful object in the public API
  should the need arise.

After talking about code concretely, we'll step back and talk about naming:

- Package names are part of the name. Don't stutter. No http.HTTPClient.
- Field and methods names are in their own scope. Don't be redundant.
  No Post.PostTitle or User.DeleteUser().
- Interfaces methods can sometimes have redundant names
  when they're meant to be "just another thing" that an object supports.
  e.g. MarshalJSON, not Marshal.
- No generic "helpers". Everything with a well-defined purpose.
- Consistency and re-use are important.
  Each new term you introduce in function and type names is a new concept users
  have to incorporate into their mental model of your package.
  Why is that called UserDetails and that called PostInfo?
  Why is that a AardvarkResponse and that a PangolinReply?
  Why is that the UserController and that the PostHandler?
  Re-use terms for similar concepts. Establish a precedent and stick to it.

And then I'll wrap it up with the rough summary
(subject to change):

- decide what you want before you write it
- leave room for it to grow in inputs and outputs
- leave room for its internals to change
- if all else fails, you can always create a new function

Thanks!
