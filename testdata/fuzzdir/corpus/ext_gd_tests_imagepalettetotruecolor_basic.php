<?php
$im = imagecreate(100, 100);
var_dump(is_resource($im));
var_dump(imageistruecolor($im));
var_dump(imagepalettetotruecolor($im));
var_dump(imageistruecolor($im));
imagedestroy($im);
?>
