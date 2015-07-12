<?php
var_dump(
    "\x00\x02\x0b\xb7" === mb_convert_encoding("\x95\x34\xb2\x35", 'UTF-32', 'GB18030')
);
?>
