<?php

$tmpDomDocument = new DOMDocument();
        
$xml = '<?xml version="1.0" encoding="UTF-8"?><dummy xmlns:xfa="http://www.xfa.org/schema/xfa-data/1.0/"><xfa:data>
  <form1>
    <TextField1>Value A</TextField1>
    <TextField1>Value B</TextField1>
    <TextField1>Value C</TextField1>
  </form1>
</xfa:data></dummy>';

$tmpDomDocument->loadXML($xml);

$dataNodes = $tmpDomDocument->firstChild->childNodes->item(0)->childNodes;

var_dump($dataNodes->length);
$datasetDom = new DOMDocument();

foreach ($dataNodes AS $node) {
    $node = $datasetDom->importNode($node, true);
    var_dump($node);
}

?>
===DONE===
