make a CLI tool that will print the sizes of all dirs to console
but in a manner that shows the relative size of each dir to each other
also maybe with colors

think Sequoia View (Windows) but in the console

because regular `du -shc` is annoying in how you cant easily tell which dir is the largest without sorting, and you cant easily see the relation of sizes in the output

maybe start with output like this

```
|                           1.9M	           Creative Cloud Files
||||||||||                  173M	  Desktop
||                          24M	Documents
||||||||||||||||||||||||    5.1G	Downloads
```

https://theasciicode.com.ar/
https://stackoverflow.com/questions/52937816/how-to-print-utf-8-or-unicode-characters-in-go-golang-on-windows

use dupefinder code as base for dir finding and size measures