# `duh`

Like `du -h`, except more human readable. Because thats the way you wanted `du` to print results anyway. Duh.

# Features

Prints the human-readable sizes of each item inside a directory, along with a text graphic to show the relative size of each item compared to the total directory size. So that you can quickly and easily find the items taking up the most space.

Console output on non-Windows systems will be colorized based on both absolute size (file or subdir total bytes) and relative size (percent of total dir space usage) to aid in quickly identifying the most significant storage usages amongst directory contents.

For best results, the input argument to `duh` should be a directory. 

Note that due to the need to compute relative sizes, output cannot be displayed until all files and subdirs inside the directory have been scanned, which might take a while if you have a lot of subdirs with a lot of files (or a slow disk). Scanning speed will be limited by your hardware. 

# Usage

Download a pre-compiled binary from [here](https://github.com/stevekm/duh/releases) and run it like this:

```
$ ./duh .
85K	.git	|
55B	.gitignore	|
2K	Makefile	|
135B	README.md	|
11.8M	build	|
569.1M	dir1	||||||||||||||||||||||||||||||
706.2M	dir2	||||||||||||||||||||||||||||||||||||||
432M	dir3	|||||||||||||||||||||||
1.9M	duh	|
156B	go.mod	|
8.2K	go.sum	|
137.1M	go1.18.3.darwin-amd64.tar.gz	|||||||
4.7K	main.go	|
8.6K	main_test.go	|
816B	notes.md	|
-----
1.8G	.

```

Otherwise, build from source using Go version 1.18+:

```
$ git clone https://github.com/stevekm/duh.git
$ cd duh
$ go build -o ./duh ./main.go
```

# Example Output

<img width="578" alt="Screen Shot 2022-08-12 at 7 09 44 PM" src="https://user-images.githubusercontent.com/10505524/184455667-3d58d014-c899-488a-b407-808c66827ebb.png">

<img width="731" alt="Screen Shot 2022-08-12 at 7 10 51 PM" src="https://user-images.githubusercontent.com/10505524/184455711-646c5629-f09e-4e97-a428-0725290f3b67.png">
