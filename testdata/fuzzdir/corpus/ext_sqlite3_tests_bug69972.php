<?php
$db = new SQLite3(':memory:');
echo "SELECTING from invalid table\n";
$result = $db->query("SELECT * FROM non_existent_table");
echo "Closing database\n";
var_dump($db->close());
echo "Done\n";

// Trigger the use-after-free
echo "Error Code: " . $db->lastErrorCode() . "\n";
echo "Error Msg: " . $db->lastErrorMsg() . "\n";
?>
