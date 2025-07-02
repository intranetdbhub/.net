# Database Table
data = [
  {"Device-ID": 3, "Class": "DELL R-710", "Colocation": "AURCOSTE LDC", "Description": "ESXi Server", "Hostname": "AURCOSTE001-ERDC001A", "Management-IP": "9.9.0.5", "Serial#": None, "Type": "Rack Server"},
  {"Device-ID": 2, "Class": "Cisco Catalyst Switch 3560 Series", "Colocation": "AURCOSTE LDC", "Description": "LDC Core Switch", "Hostname": "AURCOSTE001-CSS001A", "Management-IP": "10.68.1.2", "Serial#": None, "Type": "Enterprise Switch"},
  {"Device-ID": 1, "Class": "Cisco 2800 Series", "Colocation": "AURCOSTE LDC", "Description": "LDC Core Router", "Hostname": "AURCOSTE001-CSR001A", "Management-IP": "10.68.0.5", "Serial#": None, "Type": "Enterprise Router"}
]

# Table Print
for device in data:
    print(f"Device ID: {device['Device-ID']}, Class: {device['Class']}, Colocation: {device['Colocation']}, Description: {device['Description']}, Hostname: {device['Hostname']}, Management IP: {device['Management-IP']}, Serial#: {device['Serial#']}, Type: {device['Type']}")


