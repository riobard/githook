# GitHook

Execute commands upon git webhook. Currently only GitHub webhooks are supported.


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

    githook -addr=:4008 -conf=/etc/githook.conf


The command to be executed will have access to the following environment
variables:

    GITHOOK_SOURCE              # currently "github"
    GITHOOK_GITHUB_EVENT        # e.g. "push"
    GITHOOK_GITHUB_DELIVERY     # GUID of the GitHub webhook event


## TODO

- BitBucket webhook support
