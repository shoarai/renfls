# renfls
renfls renames all files or files matching patterns in directories.

![](https://travis-ci.org/shoarai/renfls.svg?branch=master)

#### Before
<pre>
root/
├── dir1/
│   ├── text.txt
│   └── image.jpg
├── dir2/
│   ├── text.txt
│   └── あいうえお.txt
└── dir3/
    ├── dir3-1/
    │   └── music.mp3
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
└── dir3.txt
</pre>

## Installation
```sh
$ go get github.com/shoarai/renfls
```

## Usage
```sh
$ renfls root 
```
### Option
|Option   |Description                      |
|---------|---------------------------------|
|-ext     |Rename files only matching extension list separated by ","|
|-ignore  |Exclude files matching patterns|

For example, the following command renames files not matching extension of "jpg" and "mp4" in "root" directory.

```sh
$ renfls -ext=jpg,mp4 -ignore root
```
