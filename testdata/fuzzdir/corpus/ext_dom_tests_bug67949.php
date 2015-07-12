<?php

$html = <<<HTML
<div>data</div>
<a href="test">hello world</a>
HTML;
$doc = new DOMDocument;
$doc->loadHTML($html);

$nodes = $doc->getElementsByTagName('div');

echo "testing has_dimension\n";
var_dump(isset($nodes[0]));
var_dump(isset($nodes[1]));
var_dump(isset($nodes[-1]));

echo "testing property access\n";
var_dump($nodes[0]->textContent);
var_dump($nodes[1]->textContent);

echo "testing offset not a long\n";
$offset = ['test'];
var_dump($offset);
var_dump(isset($nodes[$offset]), $nodes[$offset]->textContent);
var_dump($offset);

$something = 'test';
$offset = &$something;

var_dump($offset);
var_dump(isset($nodes[$offset]), $nodes[$offset]->textContent);
var_dump($offset);

$offset = 'test';
var_dump($offset);
var_dump(isset($nodes[$offset]), $nodes[$offset]->textContent);
var_dump($offset);

echo "testing read_dimension with null offset\n";
var_dump($nodes[][] = 1);

echo "testing attribute access\n";
$anchor = $doc->getElementsByTagName('a')[0];
var_dump($anchor->attributes[0]->name);

echo "==DONE==\n";

