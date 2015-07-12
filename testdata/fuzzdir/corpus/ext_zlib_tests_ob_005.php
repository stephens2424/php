<?php
ob_start("ob_gzhandler");
ini_set("zlib.output_compression", 0);
echo "hi\n";
?>
