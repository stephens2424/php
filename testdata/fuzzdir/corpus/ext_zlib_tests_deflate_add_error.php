<?php
$badResource = fopen("php://memory", "r+");
var_dump(deflate_add($badResource, "test"));

$resource = deflate_init(ZLIB_ENCODING_DEFLATE);
$badFlushType = 6789;
var_dump(deflate_add($resource, "test", $badFlushType));
?>
