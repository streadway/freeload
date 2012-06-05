Accepts an aggregate list of URLs and returns a JSON structure of a URI keyed
object to either data URI or original URI.

Useful for doing batch downloads of images and JSON based expansion of the
results into CSS backround URLs.

Checkout [doc.go documentation](https://github.com/streadway/freeload/blob/master/doc.go) for
details.

# Installing

Freeloader requires Go1 installable from these instructions:

http://golang.org/doc/install

Then install binary which you'll find in your GOPATH/bin:

  go get github.com/streadway/freeload/freeloader

# Contributing

Patches and enhancements welcome!  Check the issues at
https://github.com/streadway/freeload/issues and see if there is anything you
can contribute.

Make your change and tests in a branch other than `master` and use `go fmt`
before your pull request.

# License

Copyright (C) 2012 Sean Treadway <treadway@gmail.com>, SoundCloud Ltd.

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
