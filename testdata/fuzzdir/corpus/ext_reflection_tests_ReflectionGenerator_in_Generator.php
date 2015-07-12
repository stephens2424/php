<?php

function call(ReflectionGenerator $ref, $method, $rec = true) {
	if ($rec) {
		call($ref, $method, false);
		return;
	}
	var_dump($ref->$method());
}

function doCalls(ReflectionGenerator $ref) {
	call($ref, "getTrace");
	call($ref, "getExecutingLine");
	call($ref, "getExecutingFile");
	call($ref, "getExecutingGenerator");
	call($ref, "getFunction");
	call($ref, "getThis");
}

($gen = (function() use (&$gen) {
	$ref = new ReflectionGenerator($gen);

	doCalls($ref);

	yield from (function() use ($ref) {
		doCalls($ref);
		yield; // Generator !
	})();
})())->next();

?>
