<?php 
/* $Id$ */

$xw = new XMLWriter();
$xw->openMemory();
$xw->setIndent(TRUE);
$xw->setIndentString('   ');
$xw->startDocument('1.0', "UTF-8");
$xw->startElement('root');
$xw->startElementNS('ns1', 'child1', 'urn:ns1');
$xw->writeAttributeNS('ns1', 'att1', 'urn:ns1', '<>"\'&'); 
$xw->writeElement('chars', "special characters: <>\"'&");
$xw->endElement();
$xw->endDocument();
// Force to write and empty the buffer
$output = $xw->flush(true);
print $output;
?>
