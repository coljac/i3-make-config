# i3-make-config
Small tool to assemble my i3wm config file with a specific local config.

Expects `config-base` and an optional `$HOST.config` file; writes resulting config to stdout.

Looks for lines with the format `LET var_name <rest of line is value for variable>`; the last occurrence of a LET statement in config-base then in the host-specifc config file is the one with precedence. Any lines matching `LET..` or `#LET` are not written out.

Any occurence of `@<token>` will be replaced with a variables set above. For example: 

In `config-base`:

```
LET FONT Droid Sans Mono 14
LET WS1 1:Browser
...
font pango:@FONT
set $workspace1 @WS1
```
In `myhost.config`:
```
LET FONT Menlo 12
```
Then the resulting output will be:
```
font pango:Menlo 12
set $workspace1 1:Browser
```

 ## Installation

`go get github.com/coljac/i3-make-config`

