from flask import Flask, request, jsonify
import os

app = Flask(__name__)

APIKEY=''
if 'APIKEY' in os.environ:
    APIKEY = os.environ['APIKEY']

@app.before_request
def before_request():
    data = request.json
    if 'APIKEY' in data and data['APIKEY'] != APIKEY:
        response = jsonify({"message": "Unauthorized. Please provide APIKEY."})
        return response, 401

@app.route('/initsession', methods=['POST'])
def init_session():
    data = request.json
    emulator = data['emulator'] # BizHawk or Citra
    rom = data['rom'] # File name of the rom in ~/roms
    settings = data['settings'] # Setting strings to set if needed, for the emulator, tracker/extensions, and randomizer
    patches = data['patches'] # List of patches to apply to the ROM

    return jsonify({"message": "Session initialized."})

@app.route('/', methods=['POST'])
def hello_world():
	return 'Hello World'

if __name__ == '__main__':
	app.run(host='0.0.0.0', port='5000')
