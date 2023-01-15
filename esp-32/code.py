# SPDX-FileCopyrightText: 2021 ladyada for Adafruit Industries
# # SPDX-License-Identifier: MIT

import time
import board
import ssl
import socketpool 
import adafruit_sht31d
import wifi
import adafruit_minimqtt.adafruit_minimqtt as MQTT


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
print("conecting to the  broker")
mqtt_client.connect()
while True : 
    print("\nTemperature: %0.1f C" % sensor.temperature)
    print("Humidity: %0.1f %%" % sensor.relative_humidity)
    mqtt_client.loop()
    mqtt_client.publish('test/temperature', int(sensor.temperature))
    mqtt_client.publish('test/humidity', int(sensor.relative_humidity))
    time.sleep(10)
