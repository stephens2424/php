<?php
$dirname = dirname(__FILE__) . '/';
include $dirname . 'utils.inc';
$file = $dirname . '__tmp_oo_addglob.zip';

copy($dirname . 'test.zip', $file);
touch($dirname . 'foo.txt');
touch($dirname . 'bar.baz');

$zip = new ZipArchive();
if (!$zip->open($file)) {
        exit('failed');
}
$options = array('add_path' => 'baz/', 'remove_all_path' => TRUE);
if (!$zip->addGlob($dirname . '*.{txt,baz}', GLOB_BRACE, $options)) {
        echo "failed1\n";
}
if ($zip->status == ZIPARCHIVE::ER_OK) {
        dump_entries_name($zip);
        $zip->close();
} else {
        echo "failed2\n";
}
?>
