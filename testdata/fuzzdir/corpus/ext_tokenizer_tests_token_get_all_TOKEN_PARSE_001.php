<?php
$tokens = token_get_all('<?php
X::continue;
X::$continue;
$x->$continue;
X::continue();
$x->continue();
X::class;

class X {
    const CONTINUE = 1;
    public $x = self::CONTINUE + 1;
}
', TOKEN_PARSE);

array_walk($tokens, function($tk) {
  if(is_array($tk)) {
    if(($t = token_name($tk[0])) == 'T_WHITESPACE') return;
    echo "L{$tk[2]}: ".$t." {$tk[1]}", PHP_EOL;
  }
  else echo $tk, PHP_EOL;
});

echo "Done";

?>
