import board
import digitalio
import time
import wifi
import socketpool
import adafruit_minimqtt.adafruit_minimqtt as MQTT
import json
import supervisor

switches = {
    '1': 'A0',
    '2': 'A1',
    '3': 'A2',
    '4': 'A3'
}
pins = {}
for x in switches:
    pin = digitalio.DigitalInOut(getattr(board, switches[x]))
    pin.direction = digitalio.Direction.OUTPUT
    pin.value = True
    pins[switches[x]] = pin
pins["TX"] = digitalio.DigitalInOut(getattr(board, "TX"))
pins['TX'].direction = digitalio.Direction.OUTPUT
pins['TX'].value = True
btn = digitalio.DigitalInOut(board.BOOT0)
btn.direction = digitalio.Direction.INPUT
btn.pull = digitalio.Pull.UP

prev_btn_state = btn.value
btn_i = 0


try:
    import secrets
except ImportError:
    print("Error importing secrets")
    raise

mqtt_topic = 'mush/switch-group'


def connected(client, userdata, flags, rc):
    print("connected")
    client.subscribe(mqtt_topic + '/setall')


def handleAll(cient, topic, message):
    print(message)
    for pin in pins:
        if (message == 'off'):
            print(pins[pin].value)
            pins[pin].value = True
        if (message == 'on'):
            print(pins[pin].value)
            pins[pin].value = False


def handlBtnPress():
    for pin in pins:
        pins[pin].value = not pins[pin].value


def message(client, topic, message):
    id = topic.split('/').pop()
    if (id in switches):
        if (message == 'on'):
            pins[switches[id]].value = False

        if (message == 'off'):
            pins[switches[id]].value = True


wifi.radio.connect(secrets.network_id, secrets.password)
pool = socketpool.SocketPool(wifi.radio)

mqtt_client = MQTT.MQTT(
    broker="192.168.1.106",
    port=1883,
    username="pi",
    password="boopyou",
    socket_pool=pool,
    client_id='qtpi'
)

mqtt_client.on_connect = connected
mqtt_client.on_message = message
mqtt_client.add_topic_callback('mush/switch-group/setall', handleAll)
mqtt_client.will_set(mqtt_topic + '/controllerStatus', 'offline', 0, True)
mqtt_client.add_topic_callback(mqtt_topic+'/set/+',  message)

mqtt_client.connect()

mqtt_client.publish(mqtt_topic + '/controllerStatus', 'online', True, 0)
mqtt_client.subscribe(mqtt_topic+'/set/+')


while True:
    try:
        mqtt_client.loop()
        cur_btn_state = btn.value
        if (cur_btn_state and not prev_btn_state):
            handlBtnPress()
        prev_btn_state = cur_btn_state
        time.sleep(1)
    except (ValueError, RuntimeError, OSError, ConnectionError) as e:
        print("Network error, reconnecting\n", str(e))
        time.sleep(10)
        supervisor.reload()
        continue
