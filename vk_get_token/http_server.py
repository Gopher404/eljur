import config.config as config
import vk_get_token.vk as vk
from flask import Flask

app = Flask("vk")


@app.route('/get_token', methods=['GET'])
def hello_world():
    try:
        token = vk.get_token()
        return token
    except Exception as ex:
        print(ex)
        return ex.args[0], 500


def run():
    print("start server. bind: ", config.bind)
    app.run(host=config.bind[0], port=config.bind[1])
