<?php

// more info at http://sf.net/tracker/?func=detail&atid=390963&aid=1598422&group_id=27659

$contents = '
<wps:block>
<wps:var>
<wps:value/>
</wps:var>
</wps:block>';

$config = array(
'new-blocklevel-tags' => 'wps:block,wps:var,wps:value',
'newline' => 'LF'
);

$tidy = tidy_parse_string($contents, $config, 'utf8');
$tidy->cleanRepair();

var_dump($tidy->value);

?>
