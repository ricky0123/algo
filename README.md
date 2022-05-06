# Algo: Bash Alias Manager

Algo allows you to configure multi-level bash aliases through yaml. It reads yaml in from stdin and writes bash source code to stdout.

## Installation

```sh
go install github.com/ricky0123/algo@latest
```

## Example usage

Create a file `~/.algo.aliases.yml` with the following contents

```yml
# ~/.algo.aliases.yml

cm: chemzoi

g:
    $: git
    s: status
    q: commit -am "$(date)" && git push
```

Then add this line to your `~/.bashrc`:

```bash
# ~/.bashrc

if [ -f ~/.algo.aliases.yml ] && command -v algo &>/dev/null
then
    source <(cat ~/.algo.aliases.yml | algo) &> /dev/null
fi
```
