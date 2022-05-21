# greg

GReG stands for Git REpository Grepper, or Git Remote Entry Grepper if you prefer.
It is a Golang tool to search for arbitrary strings in GitHub orgs/repos.

_Other git hosting providers may be supported in the future._

## Motivation

Based on my experience of GitHub Enterprise at work,
I often find GitHub indexing to be stale (potentially very much so)
and to be highly particular about how it matches search terms against content.

To circumvent these limitations, `gitfind` queries GitHub repositories
(ensuring up-to-date content) and applies fast matching via the Aho-Corasick
algorithm.

The tool is written in Golang to provide a portable, lightweight executable.

