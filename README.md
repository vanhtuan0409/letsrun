# Letsrun

Letsrun is a command runner that allow parallel execution of blocking process. Plus it combines standard output of each process and dispatch OS signal.

If you like my work, consider buy me a coffee :D

<a href="https://www.buymeacoffee.com/sHZbgvYh0" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>

### Installation

#### Build from source
```
go get github.com/vanhtuan0409/letsrun
go install github.com/vanhtuan0409/letsrun
```

### Example

```
letsrun "sleep 10" "sleep 10" "echo Hello"
```

### Usage manual

```
Usage: letsrun [OPTIONS] COMMANDS

Background command runner and combine output into stdout

Options:
  -c	Print colorized output (default true)
  -i	Print command index indicator (default true)
  -t	Print timestamp to output
```

### Screenshot

![](screenshot/image.png)
