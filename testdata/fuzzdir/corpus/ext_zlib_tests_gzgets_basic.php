<?php

// note that gzgets is an alias to fgets. parameter checking tests will be
// the same as fgets

$f = dirname(__FILE__)."/004.txt.gz";
$h = gzopen($f, 'r');
$lengths = array(10, 14, 7, 99);
foreach ($lengths as $length) {
   var_dump(gzgets( $h, $length ) );
}

while (gzeof($h) === false) {
   var_dump(gzgets($h));
}
gzclose($h);


?>
===DONE===
