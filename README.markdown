# UniTSV

UniTSV is a simple Go library to read TSV-formatted files. It does not comply to
the [standard](http://www.iana.org/assignments/media-types/text/tab-separated-values),
as the standard disallows newlines and tab characters in fields. But TSV files
that do not contain backslashes (`\`) should be parsed correctly.

I have extended TSV in the following way:

  * Using escape characters (literal `\\`, `\n` and `\t`) to be able to
    represent any character.
  * Using (only) UTF-8 (hence 'Uni').

For now, only reading TSV-formatted files is supported, and no tests have yet
been written.
