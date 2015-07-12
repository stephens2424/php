<?php
$a = '<?xml version="1.0" encoding="UTF-8"?>
<root a="b">
	<row b="y">
		<item s="t" />
	</row>
	<row p="c">
		<item y="n" />
	</row>
</root>';
$b = str_replace(array("\n", "\r", "\t"), "", $a);
$simple_xml = simplexml_load_string($b);
print_r($simple_xml);
?>
