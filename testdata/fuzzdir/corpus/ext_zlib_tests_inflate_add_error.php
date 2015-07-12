<?php
$badResource = fopen("php://memory", "r+");
var_dump(inflate_add($badResource, "test"));
$resource = inflate_init(ZLIB_ENCODING_DEFLATE);
$badFlushType = 6789;
var_dump(inflate_add($resource, "test", $badFlushType));
?>
