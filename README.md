# aictl

Handy CLI accessing generative AI.

[![Documentation](https://pkg.go.dev/badge/github.com/go-zen-chu/golang-template)](http://pkg.go.dev/github.com/go-zen-chu/golang-template)
[![Actions Status](https://github.com/go-zen-chu/golang-template/workflows/ci/badge.svg)](https://github.com/go-zen-chu/golang-template/actions)
[![GitHub issues](https://img.shields.io/github/issues/go-zen-chu/golang-template.svg)](https://github.com/go-zen-chu/golang-template/issues)

## Usage

### Authentication

Currently aictl supports only environment variable.

```bash
export AICTL_OPENAI_API_KEY="your openai api key here"
```

### In terminal

#### A simple usage

```console
$ aictl query "How is the weather today?"
I'm sorry, but I cannot provide real-time weather information. You can check your local weather service or a weather app for the most accurate updates.
```

#### Deliver a command result to stdin with -i option

```console
$ echo "How is the weather today?" | aictl query -i
I'm sorry, but I cannot provide real-time weather information. You can check your local weather service or a weather app for the most accurate updates.

# Since, we got response as stdout, aictl's error as stderr you can pipe command results
$ aictl query "Hello" | aictl query -i | aictl query -i
Hello! How can I assist you today?

# To see how AI respond in each command, tee to stderr with pipe
$ echo $(./aictl query "Hello") | tee /dev/stderr | ./aictl query -i
Hello! How can I assist you today?
Hello! I'm here to assist you. How can I help you today?
```

#### Specify a query result format in text or json

```console
$ echo "How is the weather today?" | aictl query -o json -i 
{
  "error": "Weather information is not available."
}
```

#### Specify which language you want to have a response

```console
$ aictl query -ljapanese "Hello"
こんにちは！いかがお過ごしですか？何かお手伝いできることはありますか？

# You can specify in any language too if double quoted
$ aictl query -l"中文" "Hello"
你好！请问有什么我可以帮助你的？
```

#### Give text files and ask about the file

> [!NOTE]
> If you give large files, you may get an error from API because of the number of tokens limit

```console
$ aictl query "Why I got error in this Golang file?" -t ./testdata/go_error_sample1.go,./testdata/go_error_sample2.go
In your first Golang file, the error arises because the `fmt.Printf` function is called without providing the necessary arguments...

# By handling a result of `git diff`, you can do a code review for that
$ aictl query -t "$(git diff --name-only HEAD origin/main | tr '\n' ',' | sed 's/,$/\n/')" "Could you give me a code review for each files with filename?"
Here's a code review for each file you've provided. The review will focus on structure, style, best practices, potential improvements, and any other relevant aspects.
---
### File: `README.md`
#### Review
1. **Structure**: The README has a clear structure that helps users understand what the project is about, how to authenticate, usage in the terminal, and usage in CI.
```

### In CI

You can use aictl in any CI, but we show an example for GitHub Actions.

```yaml
jobs:
  check-aictl:
    runs-on: ubuntu-latest
    steps:
    - name: Git checkout
      uses: actions/checkout@v4
    - name: Fetch base branch
      run: git fetch origin ${{ github.event.pull_request.base.ref }}
    # Example 1, simple query
    - uses: go-zen-chu/aictl@main
      run: query "How are you?"
    # Example 2, a query with multiline
    - uses: go-zen-chu/aictl@main
      run: |
        query " \
        Let me ask a question. \
        Why I got error in this Golang file? \
        "
    # Example 3, specifing an output format
    - uses: go-zen-chu/aictl@main
      run: query -o json "How are you?"
    # Example 4, specifing an output format
```
