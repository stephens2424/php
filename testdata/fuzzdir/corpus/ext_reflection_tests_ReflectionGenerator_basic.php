<?php

function foo() {
	yield;
}

$gens = [
	(new class() {
		function a() {
			yield from foo();
		}
	})->a(),
	(function() {
		yield;
	})(),
	foo(),
];

foreach ($gens as $gen) {
	var_dump($gen);

	$gen->valid(); // start Generator
	$ref = new ReflectionGenerator($gen);

	var_dump($ref->getTrace());
	var_dump($ref->getExecutingLine());
	var_dump($ref->getExecutingFile());
	var_dump($ref->getExecutingGenerator());
	var_dump($ref->getFunction());
	var_dump($ref->getThis());
}

?>
