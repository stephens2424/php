<?php
$tarphar = new PharData(__DIR__.'/bug69958.tar');
$phar = $tarphar->convertToData(Phar::TAR); 
