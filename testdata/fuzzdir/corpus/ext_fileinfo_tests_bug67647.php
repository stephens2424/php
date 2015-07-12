<?php

$f = dirname(__FILE__) . DIRECTORY_SEPARATOR . "67647.mov";

$fi = new finfo(FILEINFO_MIME_TYPE);
var_dump($fi->file($f));
?>
+++DONE+++
