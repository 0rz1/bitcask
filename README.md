# bitcask

This is my toy project to learn golang.

Bitcask is a key-value log database. It produce value in file and store key to value's location to find it. location is consist of

1. fileno
2. offset
3. length

Every time add a key value would append log to a value file and a header file with locations.

And when reading, just find the location and read it from disk. A cache for recent key makes reading faster.

In option.go, there are some options for DB. such as go's cnt for reader, max file size, max key and value length.

I also add location header & checksum to avoid disk fail when loading.

I learn a lot in this project.

1. R/W files and handle errors
2. channel works like a thread safe queue.
3. again importance of tests 

There still some works to do. Such as
1. Managing readers to reduce system call of opening same file.
2. A monitor for deleting not outdated logs in files when idle.
3. Appending logs in multiple files, which makes more writers possible.(currently only one)
4. Recover for db's failure. A bug now is that it won't check sucess in creating new value files but not in header files. When this happen, loading return error.

But I want to do some more influential work and have no interest in this toy project anymore.