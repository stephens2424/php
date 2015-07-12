<?php

putenv("HOME=/tmp");
var_dump(getenv("HOME"));

putenv("FOO=BAR");
var_dump(getenv("FOO"));
?>
