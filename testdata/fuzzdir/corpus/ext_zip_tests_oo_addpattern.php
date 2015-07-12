<?php
$dirname = dirname(__FILE__) . '/';
include $dirname . 'utils.inc';
$file = $dirname . '__tmp_oo_addpattern.zip';

copy($dirname . 'test.zip', $file);
touch($dirname . 'foo.txt');
touch($dirname . 'bar.txt');

$zip = new ZipArchive();
if (!$zip->open($file)) {
        exit('failed');
}
$dir = realpath($dirname);
$options = array('add_path' => 'baz/', 'remove_path' => $dir);
if (!$zip->addPattern('/\.txt$/', $dir, $options)) {
        echo "failed\n";
}
if ($zip->status == ZIPARCHIVE::ER_OK) {
        dump_entries_name($zip);
        $zip->close();
} else {
        echo "failed\n";
}
?>
