#**Upload a file template (Windows Version):
# UPDATE `ecsyslabs-network-inventory`.device_inventory
# SET `Attachment` = LOAD_FILE('C:/path/to/your/file.pdf')
# WHERE Id = 5;

#**Upload a file template (Linux Version):
# -See where MySQL allows file access. You can type: SHOW VARIABLES LIKE 'secure_file_priv';
# -You’ll usually get something like /var/lib/mysql-files/. Files must live inside this folder for LOAD_FILE() to read them.
# -Put your file there and make it readable by mysqld; example of files:
# sudo cp /home/you/docs/device5.pdf /var/lib/mysql-files/
# sudo chmod 644 /var/lib/mysql-files/device5.pdf
# sudo chown mysql:mysql /var/lib/mysql-files/device5.pdf   # sometimes needed
# -If secure_file_priv is empty or NULL, set a folder in /etc/mysql/mysql.conf.d/mysqld.cnf (Debian/Ubuntu) or /etc/my.cnf (RHEL/CentOS): secure_file_priv=/var/lib/mysql-files
# -Then restart MySQL: sudo systemctl restart mysql (or mysqld/mariadb depending on distro).
# -Check size limits (big PDFs/images): SHOW VARIABLES LIKE 'max_allowed_packet';
# -Write file into the BLOB column:
# ALTER TABLE device_inventory MODIFY COLUMN `Attachment` VARCHAR(255);
# UPDATE device_inventory
# SET `Attachment` = '/mnt/files/device5.pdf'   -- or s3://bucket/key or https://...
# WHERE Id = 5;


UPDATE `ecsyslabs-network-inventory`.device_inventory
SET `Attachment` = LOAD_FILE('')
WHERE Id = ;
