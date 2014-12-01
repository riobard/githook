# GitHook

Execute commands upon git webhook


## Configuration

A sample config file (e.g. `/etc/githook.conf`)

```JSON
{
    "/my-hook-path": {
        "source": "github",
        "secret": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
        "command": "/absolute/path/to/command"
    },
    "/another-hook-path": {
        "source": "github",
        "secret": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
        "command": "/absolute/path/to/command"
    }
}
```


## Usage

    githook -adr=:4008 -conf=/etc/githook.conf
