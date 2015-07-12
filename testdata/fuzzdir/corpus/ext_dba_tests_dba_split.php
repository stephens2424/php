<?php
var_dump(dba_key_split("key1", "name"));
var_dump(dba_key_split(1));
var_dump(dba_key_split(null));
var_dump(dba_key_split(""));
var_dump(dba_key_split("name1"));
var_dump(dba_key_split("[key1"));
var_dump(dba_key_split("[key1]"));
var_dump(dba_key_split("key1]"));
var_dump(dba_key_split("[key1]name1"));
var_dump(dba_key_split("[key1]name1[key2]name2"));
var_dump(dba_key_split("[key1]name1"));

?>
===DONE===
