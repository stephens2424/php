<?php
$objstore = new SplObjectStorage();
gc_collect_cycles();

var_export($objstore);
?>
