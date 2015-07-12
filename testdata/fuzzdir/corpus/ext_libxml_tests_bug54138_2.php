<?php
define('LIBXML_BIGLINES', 1<<22);
$foos = str_repeat('<foo/>' . PHP_EOL, 65535);
$xml = <<<XML
<?xml version="1.0" encoding="UTF-8"?>
<root>
$foos
<bar/>
</root>
XML;
$dom = new DOMDocument();
$dom->loadXML($xml, LIBXML_BIGLINES);
var_dump($dom->getElementsByTagName('bar')->item(0)->getLineNo());
?>
