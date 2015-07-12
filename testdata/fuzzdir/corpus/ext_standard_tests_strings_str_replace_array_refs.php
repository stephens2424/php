<?php
$data = ['a' => 'b', 'numeric' => 1];
$ref = &$data;
$b = &$ref['a'];
$numeric = &$ref['numeric'];
var_dump(str_replace(array_keys($data), $data, "a numeric"));
var_dump($numeric);
var_dump($data['numeric']);
