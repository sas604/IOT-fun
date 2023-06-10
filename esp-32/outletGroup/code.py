import board
import digitalio
import time
import wifi
import socketpool
import adafruit_minimqtt.adafruit_minimqtt as MQTT
import json
import supervisor

switches = {
    'hum': 'A0',
    'heat': 'A1',
    'fan': 'A2',
    'light': 'A3'
}
pins = {}
for x in switches:
    pin = digitalio.DigitalInOut(getattr(board, switches[x]))
    pin.direction = digitalio.Direction.OUTPUT
    pin.value = True
    pins[switches[x]] = pin

btn = digitalio.DigitalInOut(board.BOOT0)
print(board.BOOT0)
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
    client.publish(mqtt_topic + '/controllerStatus', 'online', True, 0)
    client.subscribe(mqtt_topic + '/setall')
    client.subscribe(mqtt_topic+'/set/hum')
    client.subscribe(mqtt_topic+'/set/fun')
    client.subscribe(mqtt_topic+'/set/light')
    client.subscribe(mqtt_topic+'/set/heat')


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
    print("New message on topic {0}: {1}".format(topic, message))
    try:
        msg = json.loads(message)
    except ValueError:
        print("Error parsing message")
    if ('switch' in msg and msg['switch'] in switches):
        print('message: set ' + msg['switch'] + ' to ' + msg['value'])
        if (msg['value'] == 'on'):
            pins[switches[msg['switch']]].value = False
            client.publish(mqtt_topic + '/stateChangeComplete/' +
                           msg['switch'], json.dumps({'status': 'success', 'state': 'on'}), False, 0)
        if (msg['value'] == 'off'):
            client.publish(mqtt_topic + '/stateChangeComplete/' +
                           msg['switch'], json.dumps({'status': 'success', 'state': 'off'}), False, 0)
            pins[switches[msg['switch']]].value = True
        else:
            client.publish(mqtt_topic + '/stateChangeComplete/' + msg['switch'], json.dumps(
                {'status': 'error', 'state': 'error setting state'}), False, 0)


wifi.radio.connect(secrets.network_id, secrets.password)
pool = socketpool.SocketPool(wifi.radio)

mqtt_client = MQTT.MQTT(
    broker="192.168.1.106",
    port=1883,
    username="pi",
    password="boopyou",
    socket_pool=pool,
)
mqtt_client.on_connect = connected
mqtt_client.on_message = lambda: print('bla')
mqtt_client.add_topic_callback('mush/switch-group/setall', handleAll)
for name in switches:
    # using lambda here since passing the same function causing the client do disconect
    mqtt_client.add_topic_callback(
        mqtt_topic+'/set/' + name, lambda c, t, m: message(c, t, m))
mqtt_client.will_set(mqtt_topic + '/controllerStatus', 'ofline', 0, True)
mqtt_client.connect()


while True:
    try:
        mqtt_client.loop()
        cur_btn_state = btn.value
        if (cur_btn_state and not prev_btn_state):
            handlBtnPress()
        prev_btn_state = cur_btn_state
    except (ValueError, RuntimeError, OSError, ConnectionError) as e:
        print("Network error, reconnecting\n", str(e))
        time.sleep(10)
        supervisor.reload()
        continue
    time.sleep(0.1)
