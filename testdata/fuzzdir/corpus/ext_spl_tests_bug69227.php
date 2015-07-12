<?php

$s = new SplObjectStorage();
$s->attach($s);
gc_collect_cycles();
echo "ok";
?>
