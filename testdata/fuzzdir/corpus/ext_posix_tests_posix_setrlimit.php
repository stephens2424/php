<?php

var_dump(posix_setrlimit(POSIX_RLIMIT_NOFILE, 128, 128));
var_dump(posix_setrlimit(POSIX_RLIMIT_NOFILE, 129, 128));

?>
