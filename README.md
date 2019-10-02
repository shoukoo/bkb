<img src="https://github.com/shoukoo/bkb/workflows/build/badge.svg" class="image mod-full-width" /> <img src="https://img.shields.io/github/v/release/shoukoo/bkb?sort=semver" class="image mod-full-width" />

# bkb - `Buildkite Beaver
A CLI tool to search recent Buildkite builds

## Install
OSX
```
brew tap shoukoo/taps
brew install bkb
```
Linux & Windows can download the release on this page - https://github.com/shoukoo/bkb/releases

## Usage
```
Usage of bbk:
  bbk [flags] # run buildkite beaver
  bbk init # set token and org
  bbk show # show existing token and org

Flags:
  -help
        Print help and exist
  -version
        Print version and exit
```

## Start

Visit this page first to get a token - https://buildkite.com/user/api-access-tokens 

```
bkb init # to save org name and token in keyring 
bkb # to get the recent builds
```

