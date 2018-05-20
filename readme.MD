# icte

> if changes then execute

Simple tool to execute a command when a file(s) or a directory(ies) have a change

The idea is to enable the hot reload



## Usage

```
icte file1 directory1/ file2 file3 -c command args
```

The above command will monitor 

- file1
- file2 
- file3
- all files inside the directory1 (recursively)

If any of those files changes it executes the command with the args 

### Examples:

```bash
icte /etc/nginx/nginx.conf -c service nginx reload
```

```bash
icte main.go -c restart.sh
```

**Multiple commands**

```bash
icte main.go -c build.sh -c restart.sh
```

```bash
icte main.go -c 'grep //TODO main.go -c' -c build.sh -c restart.sh
```

## Install

If you have Go setup in your machine

```
go install github.com/zmalik/icte
```

Or

Download the binary from the latest release: <https://github.com/zmalik/icte/releases>