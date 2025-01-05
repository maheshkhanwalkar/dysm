# dysm
An object file and disassembler utility

### Background

This tool is conceptually similar to `objdump` but intends to be more powerful and easier
to use over time. It is currently in active initial development phase, so stay tuned for updates and
additional functionality!

### Supported Types
#### Object File Formats
* Mach-O \[non-universal\]

#### Architectures
* x86_64
* ARM64

### Supported Commands

```shell
dysm --hdr <file>
```
Print out header metadata for the given file

```shell
dysm --dump '<section-name-or-regex>' <file>
dysm --dump-all <file>
```

Print out a hex dump of section(s) matching the name or regex passed in,
or print out all sections of the file.

```shell
dysm --symtab <file>
```

Print out the symbol table for the file.
