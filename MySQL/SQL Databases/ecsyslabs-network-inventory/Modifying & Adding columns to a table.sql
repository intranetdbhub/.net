# **Examples of Commands to Modify a Table (Remember to start with: ALTER TABLE device_inventory) before adding the below:
# ADD COLUMN `Deployment Status` VARCHAR(200); < This is to add a regular column
# ADD COLUMN attachment LONGBLOB;  < This is to add a new blob column to attach files into a table
# DROP COLUMN `Deployment Status`; < To delete a column
# CHANGE COLUMN attachment Attachment LONGBLOB;  < This is to modify a column's name in a table
# MODIFY COLUMN `Hostname` VARCHAR(255); <This is to modify the size of the value in a cell
# MODIFY COLUMN `Description` VARCHAR(255) AFTER `Class`, MODIFY COLUMN `Colocation` VARCHAR(100) AFTER `Description`; < To modify the position or order of a column in a table

ALTER TABLE device_inventory



