package io

// Don't use standard file detection software / libmagic as it requires >= 128 bytes to be read.
// https://en.wikipedia.org/wiki/List_of_file_signatures

// ZIP files / tar files - return concat of all contained files.
