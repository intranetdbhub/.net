$language = "python"
$interface = "1.0"
#Cell Site Router - Power Edge Router Adjacency HEALTHCHECK SCRIPT

def Main():

           crt.Screen.Synchronous = True

           crt.Screen.Send("term len 0" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("sh ver | inc upt" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("sh log | i Cold" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("sh isis adj" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("sh arp" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("sh log start today | i %ROUTING-ISIS-5-ADJCHANGE" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("sh log start today | i %ROUTING-BGP-5" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("sh gnss-receiver | i \\\\"Status:\\\\\\\\|Locked at:\\\\\\\\|Lock Status\\\\\\\\|Alarms:\\\\\\\\|Satellite Count:\\\\" " + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("sh ptp int bri" + chr(13))

           crt.Screen.WaitForString("#")

           crt.Screen.Send("" + chr(13))

           crt.Screen.WaitForString("#")

Main()