<?php
$splFileObject = new SplFileObject(__FILE__);
$splFileObject->setMaxLineLen(3);
$line = $splFileObject->getCurrentLine();
var_dump($line === '<?p');
var_dump(strlen($line) === 3);
?>
