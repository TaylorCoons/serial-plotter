int potentiometer;
int static_variable = 500;

void setup() {
  Serial.begin(9600);
}

uint8_t counter = 0;

void loop() {
  unsigned int randomNumber = random(25);
  counter++;
  if (counter >= 25) {
    counter = 0;
  }

  Serial.print("RandomWalk:");
  Serial.print(randomNumber);
  Serial.print(",");
  Serial.print("Ramp:");
  Serial.println(counter);
  delay(500);
}