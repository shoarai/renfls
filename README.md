# renfls
`renfls` renames files that match patterns to the directory name of each file.<br>
`renfls` doesn't delete the files, it just renames them.

![](https://travis-ci.org/shoarai/renfls.svg?branch=master)

#### Before
<pre>
root/
├── dir1/
│   ├── text.txt
│   └── image.jpg
├── dir2/
│   ├── text.txt
│   ├── あいうえお.txt
│   └── data.dat
└── dir3/
    ├── dir3-1/
    │   ├── music.mp3
    │   └── tmp.csv
    └── dir3-2/
        └── text.txt
</pre>
#### After
<pre>
root/
├── dir1.txt
├── dir1.jpg
├── dir2.txt
├── dir2-2.txt
├── dir3.mp3
├── dir3.txt
└── ignore/
    ├── data.dat
    └── tmp.csv
</pre>

## Installation
```sh
$ go get github.com/shoarai/renfls
```

## Usage
#### CLI
```sh
$ renfls -dest=dest root
```

#### go
```go
package main

import "github.com/shoarai/renfls"

func main() {
    // Move files that match the condition in the "root" directory to the "dest" directory.
    condition := renfls.Condition{Exts: "jpg", Reg: "image*", Ignore: false}
    if e := renfls.WalkToRootSubDirName("root", "dest", condition); e != nil {
        // fmt.Println(e)
    }
}
```
### Option
|Option   |Description                      |
|---------|---------------------------------|
|-dest    |Destination to which renamed files are moved|
|-ext     |Rename files only matching extension list separated by ","|
|-ignore  |Exclude files matching patterns|

For example, the following command renames files whose extension is not "jpg" or "mp4" in the "root" directory and moves them to the "dest" directory.

```sh
$ renfls -dest=dest -ext=jpg,mp4 -ignore root
```
