<?php
if (false !== ob_gzhandler("", PHP_OUTPUT_HANDLER_START)) {
	ini_set("zlib.output_compression", 0);
	ob_start("ob_gzhandler");
}
echo "hi\n";
?>
