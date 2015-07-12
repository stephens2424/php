<?php
setlocale(LC_CTYPE, "UTF8", "en_US.UTF-8");
var_dump(escapeshellcmd('f{o}<€>'));
var_dump(escapeshellarg('f~|;*Þ?'));
var_dump(escapeshellcmd('?€®đæ?'));
var_dump(escapeshellarg('aŊł€'));

?>
