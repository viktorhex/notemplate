## notemplate

a program for creating notes

### How to use

```sh
# create entry at ~/notemplate/documents/star_entries/<DATE>-entry-<ENTRY_NUMBER>/info.toml
notemplate -t star_entries

# create entry at ~/notemplate/documents/notes/<DATE>-entry-<ENTRY_NUMBER>/info.toml
notemplate -t notes
```

### Side effects

- Creates folders at `~/notemplate/documents/*`(your notes) and `~/notemplate/templates/*`(generates default templates, you can add your own)
