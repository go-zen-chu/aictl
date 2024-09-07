# aictl

> [!CAUTION]
> Still in progress, it does not work yet. Please wait

Handy CLI accessing generative AI.

[![Documentation](https://pkg.go.dev/badge/github.com/go-zen-chu/golang-template)](http:///pkg.go.dev/github.com/go-zen-chu/golang-template)
[![Actions Status](https://github.com/go-zen-chu/golang-template/workflows/ci/badge.svg)](https://github.com/go-zen-chu/golang-template/actions)
[![GitHub issues](https://img.shields.io/github/issues/go-zen-chu/golang-template.svg)](https://github.com/go-zen-chu/golang-template/issues)

## Usage

### Authentication

Currently aictl supports only environment variable.

```bash
AICTL_OPENAI_TOKEN="your token here"
```

### In terminal

```bash
# simple usage
aictl query "How is the weather today?"

# you can pass to stdin with -i option
echo "How is the weather today?" | aictl query -i

# you can get result in any format as long as the generative AI response properly
echo "How is the weather today?" | aictl query -o json -i 

# you can pass text files and ask about the file
aictl query "Why I got error in this Golang file?" -f error_sample.go,error_sample2.go

# you can get your result in any language (default is English)
aictl query "Why I got error in this Golang file?" -f error_sample.go -l Japanese
aictl query "Why I got error in this Golang file?" -f error_sample.go -l "日本語"

# we get response as stdout, aictl's error as stderr so you can pipe
aictl query "Hello!" | aictl query -i | aictl query -i
```

### In CI

You can use aictl in any CI, but we show an example for GitHub Actions.

```bash
jobs:
  check-aictl:
    runs-on: ubuntu-latest
    steps:
    - uses: go-zen-chu/aictl@main
      run: query "Why I got error in this Golang file?"
      with:
        files: ./test
```
