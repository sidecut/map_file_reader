# map_file_reader

Read `.map` files and write their embedded contents out.

## Sample use

### List embedded sources

This will list all embedded sources.

```bash
$ map_file_reader -f lib/big_sdk.js.map -s
```

### Extract embedded source

Assuming that the js.map file contains the original soure code, this will completely extract the embedded sources to the files whose names are embedded in the .map file.  Those are the same files output by the `-s` command, illustrated above.

```bash
$ map_file_reader -f lib/big_sdk.js.map -o
```

## TODO

- [ ] Clean up command-line arguments to make them more symmetrical
