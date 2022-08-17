# How To Contribute

Thanks for your interest in contributing to this project and for taking the time
to read this guide.

## TL;DR

Pull requests should have the following guidelines:

+ Fork it `https://github.com/architeacher/colorize/fork`.

+ Create your feature branch `git checkout -b feature/my-awesome-feature`.

+ Although it is highly suggested including tests, they are not a hard
  requirement in order to get your contributions accepted.

+ [To be signed off](https://git-scm.com/book/en/v2/Git-Tools-Signing-Your-Work)

+ Please be generous describing your changes
  ```bash
  git commit -S -am "feat(api)!: :tada: Added an awesome feature." -m "[minor]" -m "Signed-off-by: John Doe <john.doe@bar.foo>"
  ```
  We follow [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) specification.

+ A logical series of [**squashed** well written commits](https://github.com/alphagov/styleguides/blob/master/git.md)

+ Push your changes `git push origin feature/my-awesome-feature`.

+ Create a pull request, and keep it small so other developers can review it quickly.

+ Keep each pull request focused on a specific topic. If you have two things
  to change, create two pull requests.
  This helps reviewers to understand the meat of your contribution.

## Coding Style

Unless explicitly stated, we follow all coding guidelines from the Go
community. While some of these standards may seem arbitrary, they somehow seem
to result in a solid, consistent codebase.

It is possible that the code base does not currently comply with these
guidelines. We are not looking for a massive PR that fixes this, since that
goes against the spirit of the guidelines. All new contributions should make a
best effort to clean up and make the code base better than they left it.
Obviously, apply your best judgement. Remember, the goal here is to make the
code base easier for humans to navigate and understand. Always keep that in
mind when nudging others to comply.

The rules:

1. All code should be pass validation checks with `make validate` which by default
   should pass the default levels of [`golangci-lint`](https://golangci-lint.run/).

2. All code should follow the guidelines covered in [Effective Go](http://golang.org/doc/effective_go.html)
   and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

3. Comment the code. Tell us the why, the history and the context.

4. Document _all_ declarations and methods, even private ones. Declare
   expectations, caveats and anything else that may be important. If a type
   gets exported, having the comments already there will ensure it's ready.

5. Variable name length should be proportional to its context and no longer.
   `noCommaALongVariableNameLikeThisIsNotMoreClearWhenASimpleCommentWouldDo`.
   In practice, short methods will have short variable names and globals will
   have longer names.

6. No underscores in package names. If you need a compound name, step back,
   and re-examine why you need a compound name. If you still think you need a
   compound name, lose the underscore.

7. No utils or helpers packages. If a function is not general enough to
   warrant its own package, it has not been written generally enough to be a
   part of a util package. Just leave it unexported and well-documented.

8.  All tests should run with Go race detection by running `make race`
    or `go test -race` and outside tooling should not be required.
    No, we don't need another unit testing framework. Assertion packages are
    acceptable if they provide _real_ incremental value.

9. Even though we call these "rules" above, they are actually just guidelines.
    Since you've read all the rules, you now know that.

If you are having trouble getting into the mood of idiomatic Go, we recommend
reading through [Effective Go](https://golang.org/doc/effective_go.html). The
[Go Blog](https://blog.golang.org) is also a great resource. Drinking the
kool-aid is a lot easier than going thirsty.

## Acknowledgements

Much of this taken from [https://github.com/cloudescape/gowsdl/blob/master/CONTRIBUTING.md](https://github.com/cloudescape/gowsdl/blob/master/CONTRIBUTING.md)
