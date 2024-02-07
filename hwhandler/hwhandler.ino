#define AMT_BUTTONS 3
#define STARTING_PIN 2
//   9 -> CE  (nRF24)
//  10 -> CSN (nRF24)
//  11 -> MO  (nRF24)
//  12 -> MI  (nRF24)
//  13 -> SCK (nRF24)

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
    pinMode(i+STARTING_PIN, INPUT_PULLUP);
  }
}

void checkButton(int btn) {
  bool is_pressed = digitalRead(btn + STARTING_PIN) == LOW;
  if (is_pressed && currentState.button_pressed[btn] == false) {
    Serial.write("BTN_");
    Serial.println(btn);
    currentState.button_pressed[btn] = true;
  } else if (!is_pressed) {
    currentState.button_pressed[btn] = false;
  }
}
//#endregion

void setup()
{
    Serial.begin(57600);
    while(!Serial){}

    Serial.println("STARTING_UP");
    setInitialState();
}

void loop()
{
    unsigned long currentTime = millis();
    currentState.last_loop_time = currentTime;

    // Every 5 seconds, we send a ping to the computer
    // If the computer does not get it (with a delta),
    // it will restart the service to connect again
    if (currentTime - currentState.last_ping_time > 5000) {
        Serial.println("PING");
        currentState.last_ping_time = currentTime;
    }

    for (int i = 0; i < AMT_BUTTONS; i++){
      checkButton(i);
    }
}
