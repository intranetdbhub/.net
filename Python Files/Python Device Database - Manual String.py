from flask import Flask, jsonify

application = Flask(__name__)

@application.route('/')
def home():
    return "I2016X Network Devices Inventory"

@application.route('/devices')
def devices():
    # Define the device data as a list of dictionaries
    devices = [
        {
            "Device_#": "1",
            "Device_ID": "AURCOSTE001A-CSR001A",
            "IOS": "Cisco",
            "Type": "router - physical",
            "Colocation": "AURCOSTE Local Datacenter | Rack 1 | 19 U - 20 U ",
            "Mgmt IP": "10.68.0.5",
        },
        {
            "Device_#": "2",
            "Device_ID": "AURCOSTE001A-CSW001A",
            "IOS": "Cisco",
            "Type": "switch - physical",
            "Colocation": "AURCOSTE Local Datacenter | Rack 1 | 21 U - 22 U ",
            "Mgmt IP": "10.68.0.5",
        },
        {
            "Device_#": "3",
            "Device_ID": "AURCOSTE001A-RDCS001A",
            "IOS": "ESXi",
            "Type": "server - physical",
            "Colocation": "AURCOSTE Local Datacenter | Rack 1 | 16 U - 17 U ",
            "Mgmt IP": "9.9.30.5",
        },
        {
            "Device_#": "4",
            "Device_ID": "OHSTE001A-WS002A",
            "IOS": "Windows Server",
            "Type": "server - Public AWS EC2 Instance",
            "Colocation": "OH Local Datacenter | OHSTE001A-VPC001A | OHSTE001A-SNET001A ",
            "Mgmt IP": "Assigned by AWS",
        },
        {
            "Device_#": "5",
            "Device_ID": "OHSTE001A-ALS001A",
            "IOS": "Amazon Linux",
            "Type": "server - Public AWS EC2 Instance",
            "Colocation": "OH Local Datacenter | OHSTE001A-VPC001A | OHSTE001A-SNET001A ",
            "Mgmt IP": "Assigned by AWS",
        },
        {
            "Device_#": "6",
            "Device_ID": "OHSTE002A-ALS001A",
            "IOS": "Amazon Linux",
            "Type": "server - Private AWS EC2 Instance",
            "Colocation": "OH Local Datacenter | OHSTE001A-VPC001A | OHSTE002A-SNET001A ",
            "Mgmt IP": "Assigned by AWS",
        },
        {
            "Device_#": "7",
            "Device_ID": "i2016xfileserver (WYSTE001A-WS001A)",
            "IOS": "None",
            "Type": "server - Microsoft Azure Blob Service",
            "Colocation": "WYS Local Datacenter | A-LDC Resource Group ",
            "Mgmt IP": "None",
        },
        {
            "Device_#": "8",
            "Device_ID": "WYSTE001A-WS001A",
            "IOS": "None - Blob Service",
            "Type": "server - Microsoft Azure Virtual Machine",
            "Colocation": "WYS Local Datacenter | A-LDC Resource Group | WYSTE001A-VNET001A ",
            "Mgmt IP": "Assigned by Azure",
        }
    ]
    return jsonify(devices)

if __name__ == "__main__":
    application.run(host='0.0.0.0', port=8080)
