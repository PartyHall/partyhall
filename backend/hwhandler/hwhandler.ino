// PARTYHALL HWHANDLER
// ESP32-C6 Version

// VERSION 2.0

#define LEDS_RELAY_PIN 15

// Those are the pin that the ESP will listen to
// Messages will be sent in the form of "BTN_[IDX OF THE PRESSED BUTTON]"
// @see https://docs.espressif.com/projects/esp-dev-kits/en/latest/esp32c6/_images/esp32-c6-devkitc-1-pin-layout.png
#define AMT_BUTTONS 14
int BUTTON_PINS[AMT_BUTTONS] = {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13};

// Non-locking serial read from http://www.gammon.com.au/serial
#define MAX_INPUT 50

struct State {
  bool button_pressed[AMT_BUTTONS];

  unsigned long last_loop_time = 0;
  unsigned long last_ping_time = 0;
};

State currentState;

//#region Button related stuff
void setInitialState() {
  currentState.last_loop_time = 0;
  currentState.last_ping_time = 0;

  for (int i = 0; i < AMT_BUTTONS; i++) {
    currentState.button_pressed[i] = false;
    pinMode(BUTTON_PINS[i], INPUT_PULLUP);
  }
}

void checkButton(int btn) {
  bool is_pressed = digitalRead(BUTTON_PINS[btn]) == LOW;
  if (is_pressed && currentState.button_pressed[btn] == false) {
    Serial.write("BTN_");
    Serial.println(btn);
    currentState.button_pressed[btn] = true;
  } else if (!is_pressed) {
    currentState.button_pressed[btn] = false;
  }
}
//#endregion

void process_data(const char* data) {
  char buffer[MAX_INPUT];
  strncpy(buffer, data, sizeof(buffer) - 1);
  buffer[sizeof(buffer) - 1] = '\0';

  char* command = strtok(buffer, " ");
  char* value_str = strtok(NULL, " ");

  if (command && value_str) {
    if (strcmp(command, "FLASH") == 0) {
      int value = atoi(value_str);

      if (value >= 0 && value <= 255) {
        analogWrite(LEDS_RELAY_PIN, value);
      }
    }
  }
}

void processIncomingByte(const byte inByte) {
  static char input_line[MAX_INPUT];
  static unsigned int input_pos = 0;

  switch (inByte) {
    case '\n':
      input_line[input_pos] = 0;
      process_data(input_line);
      input_pos = 0;
      break;

    case '\r':
      break;

    default:
      if (input_pos < (MAX_INPUT - 1)) {
        input_line[input_pos++] = inByte;
      }

      break;
  }
}

void setup() {
  pinMode(LEDS_RELAY_PIN, OUTPUT);

  Serial.begin(115200);
  while (!Serial) {}

  Serial.println("STARTING_UP");
  setInitialState();
}

void loop() {
  unsigned long currentTime = millis();
  currentState.last_loop_time = currentTime;

  // Every 5 seconds, we send a ping to the computer
  // If the computer does not get it (with a delta),
  // it will restart the service to connect again
  if (currentTime - currentState.last_ping_time > 5000) {
    Serial.println("PING");
    currentState.last_ping_time = currentTime;
  }

  for (int i = 0; i < AMT_BUTTONS; i++) {
    checkButton(i);
  }

  while (Serial.available() > 0) {
    processIncomingByte(Serial.read());
  }
}