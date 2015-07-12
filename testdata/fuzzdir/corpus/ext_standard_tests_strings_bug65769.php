<?php

$locales = array('sve', 'french', 'us', 'ru', 'czech', 'serbian');

foreach ($locales as $locale) {
	$locale = setlocale(LC_ALL, $locale);
	$lconv = localeconv();
	var_dump(
		$locale,
		$lconv['decimal_point'],
		$lconv['thousands_sep'],
		$lconv['int_curr_symbol'],
		$lconv['currency_symbol'],
		$lconv['mon_decimal_point'],
		$lconv['mon_thousands_sep']
	);
	echo '++++++++++++++++++++++', "\n";
}

?>
+++DONE+++
