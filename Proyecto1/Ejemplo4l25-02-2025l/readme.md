## Archivo para probar este ejemplo

mkdisk -size=10 -unit=k -path=/home/x/Desktop/DISCO.dk

fdisk -type=P -unit=k -name=Part1 -size=1 -path=/home/x/Desktop/disco.dk

fdisk -type=P -unit=k -name=Part2 -size=1 -path=/home/x/Desktop/disco.dk

fdisk -type=P -unit=k -name=Part3 -size=1 -path=/home/x/Desktop/disco.dk

fdisk -type=E -unit=k -name=Part4 -size=6 -path=/home/x/Desktop/disco.dk

fdisk -type=L -unit=k -name=Part5 -size=1 -path=/home/x/Desktop/disco.dk

fdisk -type=L -unit=k -name=Part6 -size=1 -path=/home/x/Desktop/disco.dk

mount -path=/home/x/Desktop/disco.dk -name=Part1

mount -path=/home/x/Desktop/disco.dk -name=Part4
mkfs -type=full -id=341a

