<?php
print_r(unpack("Q", pack("Q", 0xfffffffffffe)));
print_r(unpack("Q", pack("Q", 0)));
print_r(unpack("Q", pack("Q", 0x8000000000000002)));
print_r(unpack("Q", pack("Q", -1)));
print_r(unpack("Q", pack("Q", 0x8000000000000000)));

print_r(unpack("J", pack("J", 0xfffffffffffe)));
print_r(unpack("J", pack("J", 0)));
print_r(unpack("J", pack("J", 0x8000000000000002)));
print_r(unpack("J", pack("J", -1)));
print_r(unpack("J", pack("J", 0x8000000000000000)));

print_r(unpack("P", pack("P", 0xfffffffffffe)));
print_r(unpack("P", pack("P", 0)));
print_r(unpack("P", pack("P", 0x8000000000000002)));
print_r(unpack("P", pack("P", -1)));
print_r(unpack("P", pack("P", 0x8000000000000000)));

print_r(unpack("q", pack("q", 0xfffffffffffe)));
print_r(unpack("q", pack("q", 0)));
print_r(unpack("q", pack("q", 0x8000000000000002)));
print_r(unpack("q", pack("q", -1)));
print_r(unpack("q", pack("q", 0x8000000000000000)));
?>
