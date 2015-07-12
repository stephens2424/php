<?php

function test($locale, $value) 
{
	$newlocale = setlocale(LC_ALL, $locale);
	$conv      = localeconv();
	$sep       = $conv['decimal_point'];

	printf("%s\n--------------------------\n", $newlocale);
	printf(" sep: %s\n", $sep);
	printf("  %%f: %f\n", $value);
	printf("  %%F: %F\n", $value);
	printf("date: %s\n", strftime('%x', mktime(0, 0, 0, 12, 5, 2014)));
	printf("\n");
}

test('german', 3.41);
test('english', 3.41);
test('french', 3.41);
test('german', 3.41);
