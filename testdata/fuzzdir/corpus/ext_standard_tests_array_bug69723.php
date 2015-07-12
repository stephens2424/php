<?php
function byReference( & $array){
	foreach($array as &$item){
		$item['nanana'] = 'batman';
		$item['superhero'] = 'robin';
	}
}

$array = [
	[
	'superhero'=> 'superman',
	'nanana' => 'no nana'
	],
	[
	'superhero'=> 'acuaman',
	'nanana' => 'no nana'
	],

	];

var_dump(array_column($array, 'superhero'));
byReference($array);
var_dump(array_column($array, 'superhero'));
?>
