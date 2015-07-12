<?php
foreach (array(null, true, false, 1, "", new stdClass) as $invalid) {
    var_dump(filter_input_array(INPUT_POST, $invalid));
    var_dump(filter_var_array(array(), $invalid));
}
