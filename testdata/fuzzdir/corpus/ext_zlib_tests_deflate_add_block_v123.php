<?php

$resource = deflate_init(ZLIB_ENCODING_GZIP);
var_dump(deflate_add($resource, "aaaaaaaaaaaaaaaaaaaaaa", ZLIB_BLOCK));

?>
===DONE===
