import board 
import digitalio 
import time
import wifi
import socketpool
import adafruit_minimqtt.adafruit_minimqtt as MQTT
import json

switches = {
    'hum': 'A0',
    'heat': 'A1',
    'fan' : 'A2',
    'light': 'A3'
};
pins = {}
for x in switches:
    pin =  digitalio.DigitalInOut(getattr(board, switches[x]))
    pin.direction = digitalio.Direction.OUTPUT
    pin.value = False
    pins[switches[x]] = pin


try: 
    import secrets
except ImportError:
    print("Error importing secrets")
    raise    

mqtt_topic = 'mush/switch-group'

def connected(client, userdata, flags, rc): 
        client.subscribe(mqtt_topic + '/set')
    

def message(client, topic, message):
    print("New message on topic {0}: {1}".format(topic, message))
    try :
        msg = json.loads(message)
    except ValueError:
        print("Error parsing message")
    if('switch' in msg and msg['switch'] in switches):
        if(msg['value'] == 'on'):
            pins[switches[msg['switch']]].value = False;
        if(msg['value'] == 'off'):
            print(pins)
            pins[switches[msg['switch']]].value = True;
            



wifi.radio.connect(secrets.network_id, secrets.password)
pool = socketpool.SocketPool(wifi.radio)

mqtt_client = MQTT.MQTT(
    broker = "192.168.1.106",
    port= 1883,  
    username = "pi",
    password = "boopyou",
    socket_pool = pool,
)
mqtt_client.on_connect = connected
mqtt_client.on_message = message
# mqtt_client.on_disconnected = disconnected
mqtt_client.connect()


while True:
    mqtt_client.loop()
    time.sleep(1)