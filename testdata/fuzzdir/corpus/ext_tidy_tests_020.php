<?php

$tidy = new tidy();
$str  = <<<EOF
<p>Isto é um texto em Português<br>
para testes.</p>
EOF;

$tidy->parseString($str, array('output-xhtml'=>1), 'latin1');
$tidy->cleanRepair();
$tidy->diagnose();
var_dump(tidy_warning_count($tidy) > 0);
var_dump(strlen($tidy->errorBuffer) > 50);

echo $tidy;
?>
