<?php
		$a = new tidy(dirname(__FILE__)."/007.html");
		echo "Current Value of 'tidy-mark': ";
		var_dump($a->getopt("tidy-mark"));
		echo "Current Value of 'error-file': ";
		var_dump($a->getopt("error-file"));
		echo "Current Value of 'tab-size': ";
		var_dump($a->getopt("tab-size"));

		var_dump($a->getopt('bogus-opt'));
		var_dump(tidy_getopt($a, 'non-ASCII string рсч'));
?>
