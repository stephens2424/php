<?php

// note that gzgets is an alias to fgets. parameter checking tests will be
// the same as gzgets

$f = dirname(__FILE__)."/004.txt.gz";
$h = gzopen($f, 'r');

$count = 0;
while (gzeof($h) === false) {
   $count++;
   echo fgetc( $h );
}

echo "\ncharacters counted=$count\n";
gzclose($h);

?>
===DONE===
