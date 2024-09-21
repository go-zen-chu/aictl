# aictl

Handy CLI accessing generative AI.

[![Documentation](https://pkg.go.dev/badge/github.com/go-zen-chu/golang-template)](http:///pkg.go.dev/github.com/go-zen-chu/golang-template)
[![Actions Status](https://github.com/go-zen-chu/golang-template/workflows/ci/badge.svg)](https://github.com/go-zen-chu/golang-template/actions)
[![GitHub issues](https://img.shields.io/github/issues/go-zen-chu/golang-template.svg)](https://github.com/go-zen-chu/golang-template/issues)

## Usage

### Authentication

Currently aictl supports only environment variable.

```bash
export AICTL_OPENAI_API_KEY="your openai api key here"
```

### In terminal

```console
# A simple usage
$ aictl query "How is the weather today?"
I'm sorry, but I cannot provide real-time weather information. You can check your local weather service or a weather app for the most accurate updates.

# Deliver a command result to stdin with -i option
$ echo "How is the weather today?" | aictl query -i
I'm sorry, but I cannot provide real-time weather information. You can check your local weather service or a weather app for the most accurate updates.

# Specify a result in text or json
$ echo "How is the weather today?" | aictl query -o json -i 
{
  "error": "Weather information is not available."
}

# We get response as stdout, aictl's error as stderr so you can pipe
$ aictl query "Hello" | aictl query -i | aictl query -i
Hello! How can I assist you today?

# To see how AI responding, tee to stderr with pipe
$ echo $(./aictl query "Hello") | tee /dev/stderr | ./aictl query -i
Hello! How can I assist you today?
Hello! I'm here to assist you. How can I help you today?

# Specify which language you want to have a response
$ aictl query -ljapanese "Hello"
こんにちは！いかがお過ごしですか？何かお手伝いできることはありますか？

# You can specify in any language too if double quoted
$ aictl query -l"中文" "Hello"
你好！请问有什么我可以帮助你的？

# Give text files and ask about the file. If you have large files, you may get an error from API because of a number of tokens restriction
$ aictl query "Why I got error in this Golang file?" -t ./testdata/go_error_sample1.go,./testdata/go_error_sample2.go
In your first Golang file, the error arises because the `fmt.Printf` function is called without providing the necessary arguments...

# Ask about git diff files
$ aictl query -git-diff "Could you give me a code review to these files?"
## compare with <commit>
$ aictl query -git-diff <commit> "Could you give me a code review to these files?"
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
