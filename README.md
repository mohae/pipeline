Pipeline
========

Pipeline is an exercise in pipelines with Go. It is based on http://blog.golang.org/pipelines and mostly follows the code in the post and the linked files. This merely implements it as a CLI using commands to execute the appropriate pipeline example, with some additional functionality and minor changes.

This was created as an exercise for me and to, eventually, compare processing speeds of varying degrees of parallelism. I may also add other algorithms in the future so I can compare speeds of various algorithms.

No optimization has been done to this code, e.g. adding a concurrent walker instead of using the stdlib, etc., since the main focus was exploring pipelines in a more consistent manner than I had previously been.

This is meant to assist in learning about pipelines, go routines, concurrency, and channels; including error communication, termination signalling, and using wait groups. This is not meant to be used for any other reason, do not expect it to fulfill any functional need other than being a learning tool.

## Square
The `square` command implements the squaring example and accepts a variadic list of ints to square. The output is the square of the list in the order processed, which may not be the order you specify.

    pipeline square 2 3 4
	
will ouput:

    4
	9
	16
	
## MD5
The `md5` command implements the bounded parallel version of the MD5 pipeline example. It accepts a path, which can either be a file or a directory and outputs each file's hash along with the filename. 

MD5 accepts 1 flag, `--parallel`, which allows you to specify the degree of parallelism for file processing. By default this is set to `10`.

    pipeline md5 path/to/directory
	
or 

    pipeline md5 --parallel=20 path/to/directory
	
_Even though `-p` is documented as a valid short flag for `--parallel`, it is not currently supported because of a bug in the underlying package_

There is also a serial implementation of the MD5 processing which has not been exposed. This will be exposed when timing is added so processing speed can be compared.

## SHA256
The `sha255` command implements the bounded parallel version of the MD5 pipeline example. It accepts a path, which can either be a file or a directory and outputs each file's hash along with the filename. 

SHA256 accepts 1 flag, `--parallel`, which allows you to specify the degree of parallelism for file processing. By default this is set to `10`.

    pipeline sha256 path/to/directory
	
or 

    pipeline sha256 --parallel=20 path/to/directory
	
_Even though `-p` is documented as a valid short flag for `--parallel`, it is not currently supported because of a bug in the underlying package_

When timing is added, a serial version will be added for comparative reasons.

## TODO
Add timing information so processing speeds of varying degrees of parallelism can be compared.`