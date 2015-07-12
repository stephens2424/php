<?php
$buffer = '<html></html>';
$config = array(
  'indent' => true, // AutoBool
  'indent-attributes' => true, // Boolean
  'indent-spaces' => 3, // Integer
  'language' => 'de'); // String
$tidy = new tidy();
$tidy->parseString($buffer, $config);
$c = $tidy->getConfig();
var_dump($c['indent']);
var_dump($c['indent-attributes']);
var_dump($c['indent-spaces']);
var_dump($c['language']);
?>
