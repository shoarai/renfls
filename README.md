# renfls
renfls renames files in a directory.
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
