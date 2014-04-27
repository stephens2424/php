<?

$myvar = 1;
$yourvar = 2;

$func = function ($somevar, $othervar) use ($myvar, $yourvar) {
  return $myvar + $yourvar;
};
