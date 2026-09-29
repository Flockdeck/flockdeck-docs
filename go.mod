// The site is static files and has no Go code. This file exists so that
// versions published under a different licence can be retracted: a
// retraction is only readable from the go.mod of the highest version.
//
// v0.1.0 and v0.1.1 carried the MIT licence; every later version is
// PolyForm Noncommercial 1.0.0 (see LICENSE).

module github.com/Flockdeck/flockdeck-docs

go 1.16

retract [v0.1.0, v0.1.1] // Published under a different licence (MIT); use a later version.
