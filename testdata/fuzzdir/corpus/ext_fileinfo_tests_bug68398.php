<?php

$f = new finfo(FILEINFO_MIME);
var_dump($f->file(dirname(__FILE__) . DIRECTORY_SEPARATOR . '68398.zip'));
?>
+++DONE+++
