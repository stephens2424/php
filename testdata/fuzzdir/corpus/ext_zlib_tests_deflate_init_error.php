<?php
var_dump(deflate_init(42));
var_dump(deflate_init(ZLIB_ENCODING_DEFLATE, ['level' => 42]));
var_dump(deflate_init(ZLIB_ENCODING_DEFLATE, ['level' => -2]));
var_dump(deflate_init(ZLIB_ENCODING_DEFLATE, ['memory' => 0]));
var_dump(deflate_init(ZLIB_ENCODING_DEFLATE, ['memory' => 10]));
?>
