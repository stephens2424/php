<?php
$doc = new DOMDocument();
$result = $doc->loadHTMLFile("");
assert('$result === false');
$doc = new DOMDocument();
$result = $doc->loadHTMLFile("text.html\0something");
assert('$result === false');
?>
