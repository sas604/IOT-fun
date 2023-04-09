# SPDX-FileCopyrightText: 2021 ladyada for Adafruit Industries
# # SPDX-License-Identifier: MIT

import time
import board
import busio
import ssl
import socketpool 
import adafruit_sht31d
import wifi
import adafruit_minimqtt.adafruit_minimqtt as MQTT
import terminalio
import supervisor
from adafruit_display_text import bitmap_label



try :
    import secrets
except ImportError: 
    print("WIFI secrets are kept in secrets.py")
    raise
wifi.radio.connect(secrets.network_id, secrets.password)


i2c = board.STEMMA_I2C()  # uses board.SCL and board.SDA
sensor = adafruit_sht31d.SHT31D(i2c)

def connected(client, userdata, flags, rc): 
    print('Conected to MQTT broker')
    client.subscribe('feed/onoff')
def disconnected(client, userdata, rc): 
    print('Disconected from broker')

def message(client, topic, message): 
    print("New message on topic  {0}: {1}".format(topic, message))    


# create a socket pool 
pool = socketpool.SocketPool(wifi.radio)
print(wifi.radio)
mqtt_client = MQTT.MQTT(
    broker = "192.168.1.106",
    port= 1883,  
    username = "pi",
    password = "boopyou",
    socket_pool = pool,
)

mqtt_client.on_connect = connected
mqtt_client.on_disconnected = disconnected
print("conecting to the  broker")
mqtt_client.connect()
while True :
    try: 
        text ="\nTemperature: %0.1fC \nHumidity: %0.1f%% \nMQTT:% s"% (sensor.temperature, sensor.relative_humidity, "Connected" if mqtt_client.is_connected() else "Disconnected" )
        text_area = bitmap_label.Label(terminalio.FONT, text=text, scale= 2)
        text_area.x = 0
        text_area.y = -10
        board.DISPLAY.show(text_area)
        mqtt_client.loop()
        mqtt_client.publish('test/temperature', str(sensor.temperature) + "," + str(sensor.relative_humidity)) 
    except (ValueError, RuntimeError, OSError, ConnectionError) as e:
        print("Network error, reconnecting\n", str(e))
        time.sleep(60)
        supervisor.reload()
        continue
    time.sleep(10)
