# wrdcntr - simple word occurence counter with sentence indices saving

This application reads file given as an input param (command line flag) and counts number of occurencies of each word found.

## Application running

To launch application with "~/some_file.txt" as input launch command below in main directory of this repo:

```bash
go run ./ -f ~/some_file.txt
```

## Test running

This project contains a set of unit tests for statistics collection object. These tests are written using default Go testing framework and can be launched using following command:

```bash
go test -v ./...
```

## Information about limitations, assumptions and shortcuts done intentionally

The goal of this application is touching area of language processing. The NLP (Natural Language Processing) topic is large and complex, much beyond this simple exercice. To make this exercise possible to be finished following assumptions/shortcuts has been done:

- all words are processed in lower-character version only (including "I", "Mr." and so on)
- application defines a list of special "words" that are processed in different way ("I'm" or "You're" are expanded to two words "I am" and "You are" respectively, "i.e." is stored in statistics as a one string with dots)
- 