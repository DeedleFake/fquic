fquic
=====

[![Go Reference](https://img.shields.io/github/v/tag/DeedleFake/fquic)](https://pkg.go.dev/pkgs.go.dev/mod/github.com/DeedleFake/fquic)
[![Go Report Card](https://goreportcard.com/badge/github.com/DeedleFake/fquic)](https://goreportcard.com/report/github.com/DeedleFake/fquic)

*Note: I wouldn't particularly recommend using this package. It was mostly just an experiment for myself. In writing it, I got a little more used to quic-go's API and found that a lot of the biggest problems that I have with it are actually documentation issues, so this package is pretty much pointless overall.*

fquic provides an experimental wrapper for quic-go that simplifies the API and makes it a little friendlier. It attempts to match the interfaces and practices of the standard libraries `net` package as much as possible.
