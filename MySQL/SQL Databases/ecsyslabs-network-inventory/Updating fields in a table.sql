# ***Example of updating a field in a table (Using Id Primary Key - Recommended):
# UPDATE device_inventory
# SET `Colocation` = 'FRECASTE LDC'
# WHERE `Id` = 5;

# ***Example of updating a field in a table (Using 2 other related fields - Could break if these change in the future):
# UPDATE device_inventory
# SET `Colocation` = 'FRECASTE LDC'
# WHERE `Class` = 'DELL R-710'
# AND `Hostname` = 'AURCOSTE001A-ERDC001A';

UPDATE device_inventory